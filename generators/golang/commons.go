package golang

import (
	"bytes"
	"fmt"
	"slices"
	"strings"
	"text/template"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
	"golang.org/x/tools/imports"
)

func generateGoHelperFuncsFile(packageName string) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFS(goTemplates, "templates/helperFuncsFile.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse helper funcs template: %w", err)
	}

	err = tmpl.ExecuteTemplate(&buf, "goHelperFuncsFile", map[string]string{
		"PackageName": packageName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to execute helper funcs template: %w", err)
	}
	return buf.Bytes(), nil
}

func collectAnyAndAllAuthMethods(endpoint *spec.Endpoint, api *spec.Specification) (reqAuthMethodsAll, reqAuthMethodsAny []GoAuthMethod) {
	if endpoint.Auth != nil {
		for _, auth := range api.Auth {
			if slices.Contains(endpoint.Auth.All, auth.ID) {
				reqAuthMethodsAll = append(reqAuthMethodsAll, GoAuthMethod{
					ID:            auth.ID,
					Name:          auth.Name,
					Type:          getGoAuthMethodType(auth.Type),
					TransportName: auth.TransportName,
					Description:   auth.Description,
					Format:        auth.Format,
				})
				continue
			}

			if slices.Contains(endpoint.Auth.Any, auth.ID) {
				reqAuthMethodsAny = append(reqAuthMethodsAny, GoAuthMethod{
					ID:            auth.ID,
					Name:          auth.Name,
					Type:          getGoAuthMethodType(auth.Type),
					TransportName: auth.TransportName,
					Description:   auth.Description,
					Format:        auth.Format,
				})
			}
		}
	}
	return
}

func collectRequest(endpoint *spec.Endpoint, endpointName string, reqAuthMethodsAll, reqAuthMethodsAny []GoAuthMethod) GoRequestStructDef {
	headers, query, path := []GoParamDef{}, []GoParamDef{}, []GoParamDef{}

	for _, header := range endpoint.Headers {
		headers = append(headers, getGoParamDef(header))
	}
	for _, queryParam := range endpoint.QueryParams {
		query = append(query, getGoParamDef(queryParam))
	}
	for _, pathParam := range endpoint.PathParams {
		path = append(path, getGoParamDef(pathParam))
	}

	requestStructs := []GoStructDef{}
	var reqBodyName *string
	// Collect structs from RequestBody
	if endpoint.RequestBody != nil && endpoint.ContentType != nil && *endpoint.ContentType == "application/json" {
		name := fmt.Sprintf("%sRequestBody", exportedName(endpointName))
		reqBodyName = &name
		requestStructs = append(requestStructs, collectStructsFromRequestBody(exportedName(endpointName), endpoint.RequestBody)...)
	}

	slices.SortFunc(requestStructs, func(a, b GoStructDef) int {
		return strings.Compare(a.Name, b.Name)
	})

	return GoRequestStructDef{
		Name:              exportedName(endpointName) + "Request",
		Description:       endpoint.Description,
		ContentType:       *endpoint.ContentType,
		Method:            string(endpoint.Method),
		Path:              endpoint.Path,
		MaxBodyBytes:      endpoint.MaxBodyBytes,
		HeaderParams:      headers,
		QueryParams:       query,
		PathParams:        path,
		SupportingStructs: requestStructs,
		RequestBodyName:   reqBodyName,
		AuthAll:           reqAuthMethodsAll,
		AuthAny:           reqAuthMethodsAny,
	}
}

func collectResponses(endpoint *spec.Endpoint, endpointName string) []GoResponseStructDef {
	responseStructs := []GoResponseStructDef{}
	for _, status := range utils.SortedMapKeys(endpoint.Responses) {
		response := endpoint.Responses[status]
		respName := fmt.Sprintf("%s%dResponse", exportedName(endpointName), status)
		structs := []GoStructDef{}
		var respBodyName *string
		if response.Body != nil && response.ContentType != nil && *response.ContentType == "application/json" {
			name := respName + "Body"
			respBodyName = &name
			structs = append(structs, collectStructsFromResponseBody(respName, response.Body)...)
		}
		headers := []GoParamDef{}
		for _, header := range response.Headers {
			headers = append(headers, getGoParamDef(header))
		}
		responseStructs = append(responseStructs, GoResponseStructDef{
			StatusCode:        status,
			Name:              respName,
			Description:       response.Description,
			ContentType:       *response.ContentType,
			Headers:           headers,
			SupportingStructs: structs,
			ResponseBodyName:  respBodyName,
		})
	}
	slices.SortFunc(responseStructs, func(a, b GoResponseStructDef) int {
		return strings.Compare(a.Name, b.Name)
	})
	return responseStructs
}

func collectStructsFromResponseBody(parentName string, body *spec.HTTPBody) []GoStructDef {
	structs := []GoStructDef{}

	field := spec.Field{
		Type:        spec.FieldTypeObject,
		Description: body.Description,
		Properties:  body.Properties,
		Required:    true,
	}

	collectStructs(parentName, "Body", &field, &structs)
	return structs
}

func collectStructsFromRequestBody(parentName string, body *spec.HTTPBody) []GoStructDef {
	structs := []GoStructDef{}

	field := spec.Field{
		Type:        spec.FieldTypeObject,
		Description: body.Description,
		Properties:  body.Properties,
		Required:    true,
	}

	collectStructs(parentName, "RequestBody", &field, &structs)
	return structs
}

func collectStructs(parentName string, fieldName string, field *spec.Field, out *[]GoStructDef) {
	if field.Type != spec.FieldTypeObject {
		return
	}

	structName := parentName + exportedName(fieldName)

	structDef := GoStructDef{
		Name:   structName,
		Fields: []GoFieldDef{},
	}

	for _, propName := range utils.SortedMapKeys(field.Properties) {
		prop := field.Properties[propName]
		goFieldType, shouldRecurseValidate := goTypeFromSpecField(prop, structName+exportedName(propName))

		goTypeElem := ""
		if prop.Type == spec.FieldTypeArray {
			goTypeElem = goFieldType[2:] // trim the "[]"
		}

		fieldDef := GoFieldDef{
			Name:            exportedName(propName),
			Description:     prop.Description,
			GoType:          goFieldType,
			Tag:             `json:"` + exportedName(propName) + omitEmptyTag(prop.Required) + `"`,
			RecurseValidate: shouldRecurseValidate,
			IsArray:         prop.Type == spec.FieldTypeArray,
			ElemGoType:      goTypeElem,
			Required:        prop.Required,
			NonEmpty:        prop.NonEmpty,
		}

		structDef.Fields = append(structDef.Fields, fieldDef)

		// Recursively collect nested structs
		switch prop.Type {
		case spec.FieldTypeObject:
			collectStructs(structName, propName, &prop, out)
		case spec.FieldTypeArray:
			// If the property is an array, check if the items are objects
			if prop.Items != nil && prop.Items.Type == spec.FieldTypeObject {
				collectStructs(structName, propName+"Item", prop.Items, out)
			}
		}
	}
	*out = append(*out, structDef)
}

func exportedName(name string) string {
	// return cases.Title(language.English).String(name)
	return strings.ToUpper(name[:1]) + name[1:]
}

func goTypeFromSpecField(field spec.Field, parentName string) (typ string, shouldRecurseValidate bool) {
	switch field.Type {
	case spec.FieldTypeString:
		return "string", false
	case spec.FieldTypeNumber:
		return "float64", false
	case spec.FieldTypeBoolean:
		return "bool", false
	case spec.FieldTypeObject:
		return exportedName(parentName), true
	case spec.FieldTypeArray:
		elemType, elemRecurse := goTypeFromSpecField(*field.Items, parentName+"Item")
		return "[]" + elemType, elemRecurse
	default:
		panic("unsupported field type " + field.Type)
	}
}

func goTypeFromParam(param spec.Param) string {
	switch param.Type {
	case spec.ParamTypeString:
		return "string"
	case spec.ParamTypeNumber:
		return "float64"
	case spec.ParamTypeBoolean:
		return "bool"
	default:
		panic("unsupported param type " + param.Type)
	}
}

func omitEmptyTag(required bool) string {
	if required {
		return ""
	}
	return ",omitempty"
}

func getGoAuthMethodType(authType spec.AuthMethodType) GoAuthMethodType {
	switch authType {
	case spec.AuthMethodHeader:
		return GoAuthMethodTypeHeader
	default:
		panic("unsupported auth method type " + string(authType))
	}
}

func getGoParamDef(param spec.Param) GoParamDef {
	return GoParamDef{
		Name:          exportedName(param.Name),
		TransportName: param.TransportName,
		GoType:        goTypeFromParam(param),
		Description:   param.Description,
		Required:      param.Required,
	}
}

// utils
// run go mod tidy
func runGoModTidy(dir string) error {
	cmdGoModTidy := "go mod tidy"

	return utils.ExecCommand(cmdGoModTidy, dir)
}

// run go mod init
func runGoModInit(moduleName string, dir string) error {
	cmdGoModInit := fmt.Sprintf("go mod init %s", moduleName)

	return utils.ExecCommand(cmdGoModInit, dir)
}

func formatWithImports(src []byte) ([]byte, error) {
	opts := &imports.Options{
		TabWidth:  8,
		TabIndent: true,
		Comments:  true,
	}

	formattedSrc, err := imports.Process("", src, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to format code with goimports: %w", err)
	}
	return formattedSrc, nil
}
