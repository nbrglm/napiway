package golang

import (
	"bytes"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
	"golang.org/x/tools/imports"
)

func TypesDataFromSpec(specification *spec.Specification) []TypeData {
	var types []TypeData
	for _, t := range specification.Schemas {
		types = append(types, TypeData{
			Name:        exportedName(t.Name),
			Description: t.Description,
			Fields:      getFieldsDataFromSpecFields(t.Properties, specification.Schemas),
			Enum:        t.Enum,
		})
	}
	sortTypesByName(&types)
	return types
}

func AuthMethodsFromSpec(specification *spec.Specification) []AuthMethodData {
	authMethods := make([]AuthMethodData, len(specification.Auth))
	for i, auth := range specification.Auth {
		authMethods[i] = AuthMethodData{
			ID:            auth.ID,
			Name:          exportedName(auth.Name),
			TransportName: auth.TransportName,
			Type:          AuthMethodType(auth.Type),
			Description:   auth.Description,
			Format:        auth.Format,
		}
	}
	sortAuthMethodsByID(&authMethods)
	return authMethods
}

func RequestResponsesDataFromEndpointDef(endpointIdx int, specification *spec.Specification) (RequestData, error) {
	endpoint := specification.Endpoints[endpointIdx]

	var authMethodAll, authMethodAny []AuthMethodData

	if endpoint.Auth != nil {
		var err error
		authMethodAll, err = getAuthMethods(specification, endpoint.Auth.All)
		if err != nil {
			return RequestData{}, err
		}
		authMethodAny, err = getAuthMethods(specification, endpoint.Auth.Any)
		if err != nil {
			return RequestData{}, err
		}
	}

	var requestBodyName *string
	if endpoint.BodyName != nil {
		requestBodyName = new(string)
		*requestBodyName = exportedName(*endpoint.BodyName)
	}

	responses := make([]ResponseData, len(endpoint.Responses))
	has413 := false
	for i, resp := range endpoint.Responses {
		if resp.Status == 413 {
			has413 = true
		}
		var responseBodyName *string
		if resp.BodyName != nil {
			responseBodyName = new(string)
			*responseBodyName = exportedName(*resp.BodyName)
		}
		responses[i] = ResponseData{
			StatusCode:       resp.Status,
			Name:             exportedName(endpoint.Name + strconv.Itoa(resp.Status)),
			Description:      resp.Description,
			RawBody:          resp.RawBody,
			ContentType:      *resp.ContentType,
			Headers:          mapSpecParamToParamData(resp.Headers),
			ResponseBodyName: responseBodyName,
		}
	}

	// If the responses don't include a 413 (Payload Too Large), add a default one for the case when the request body exceeds MaxBodyBytes
	if !has413 && !endpoint.RawBody && requestBodyName != nil {
		// If RawBody is true, it means the generated code will not be reading/unmarshalling the request body, so we don't need to worry about adding a 413 response.
		// If requestBodyName is nil, it means there is no request body, so we also don't need to worry about adding a 413 response.
		desc := "Payload Too Large - the request body exceeds the maximum allowed size"
		responses = append(responses, ResponseData{
			StatusCode:       413,
			Name:             exportedName(endpoint.Name + "413Response"),
			Description:      &desc,
			RawBody:          true, // Since the body is too large to be read, we set RawBody to true to indicate that the generated code should not try to read/unmarshal the response body.
			Headers:          []ParamData{},
			ResponseBodyName: nil,
		})
	}

	sortResponsesByStatusCode(&responses)

	return RequestData{
		Name:            exportedName(endpoint.Name + "Req"),
		Description:     endpoint.Description,
		Method:          string(endpoint.Method),
		Path:            endpoint.Path,
		MaxBodyBytes:    endpoint.MaxBodyBytes,
		ContentType:     *endpoint.ContentType,
		RawBody:         endpoint.RawBody,
		RequestBodyName: requestBodyName,
		PathParams:      mapSpecParamToParamData(endpoint.PathParams),
		QueryParams:     mapSpecParamToParamData(endpoint.QueryParams),
		HeaderParams:    mapSpecParamToParamData(endpoint.Headers),
		AuthAll:         authMethodAll,
		AuthAny:         authMethodAny,
		Responses:       responses,
	}, nil
}

func mapSpecParamToParamData(params []spec.Param) []ParamData {
	resParams := make([]ParamData, len(params))
	for i, pathParam := range params {
		resParams[i] = ParamData{
			Name:          exportedName(pathParam.Name),
			TransportName: pathParam.TransportName,
			Type:          getPathParamTypeFromSpecPathParamType(pathParam.Type),
			Required:      pathParam.Required,
			Description:   pathParam.Description,
		}
	}
	sortParamsByName(&resParams)
	return resParams
}

func getPathParamTypeFromSpecPathParamType(paramType spec.ParamType) string {
	switch paramType {
	case spec.ParamTypeString:
		return TypeStrString
	case spec.ParamTypeInteger:
		return TypeStrInteger
	case spec.ParamTypeDouble:
		return TypeStrDouble
	case spec.ParamTypeBoolean:
		return TypeStrBoolean
	default:
		return string(paramType)
	}
}

func getAuthMethods(specification *spec.Specification, ids []string) ([]AuthMethodData, error) {
	ams := make([]AuthMethodData, 0, len(ids))
	missingAuths := []string{}

	for _, id := range ids {
		idx := slices.IndexFunc(specification.Auth, func(auth spec.AuthMethod) bool {
			return auth.ID == id
		})
		if idx == -1 {
			missingAuths = append(missingAuths, id)
		} else {
			auth := specification.Auth[idx]
			ams = append(ams, AuthMethodData{
				ID:            auth.ID,
				Name:          exportedName(auth.Name),
				TransportName: auth.TransportName,
				Type:          AuthMethodType(auth.Type),
				Description:   auth.Description,
				Format:        auth.Format,
			})
		}
	}
	if len(missingAuths) > 0 {
		return nil, fmt.Errorf("auth methods with ids %v not found in specification", missingAuths)
	}

	sortAuthMethodsByID(&ams)
	return ams, nil
}

func getFieldsDataFromSpecFields(fields []*spec.SchemaField, schemas []*spec.Schema) []TypeFieldData {
	if len(fields) == 0 {
		return nil
	}
	res := make([]TypeFieldData, len(fields))
	for i, field := range fields {
		typ, isPrimitive := getTypeDataFieldTypeFromSpecFieldType(field.Type)
		isEnum := IsTypeEnum(isPrimitive, typ, schemas)
		ptrType := false
		if field.IsArray {
			ptrType = false
		} else if !field.Required {
			ptrType = true
		} else if !isPrimitive && !isEnum {
			ptrType = true
		}
		var tagBuilder strings.Builder
		tagBuilder.WriteString("json:")
		tagBuilder.WriteRune('"')
		tagBuilder.WriteString(exportedName(field.Name))
		tagBuilder.WriteString(getOmitEmpty(field.Required))
		tagBuilder.WriteRune('"')
		res[i] = TypeFieldData{
			Name:               exportedName(field.Name),
			Description:        field.Description,
			Type:               typ,
			PtrType:            ptrType,
			IsNonPrimitiveType: !isPrimitive,
			Tag:                tagBuilder.String(),
			IsArray:            field.IsArray,
			IsEnum:             isEnum,
			Required:           field.Required,
			NonEmpty:           field.NonEmpty,
		}
	}
	sortTypeFieldsByName(&res)
	return res
}

// Returns the Go type string for a given SchemaFieldType, and a boolean indicating whether the type is a primitive type (i.e. one of the types in TypeStr) or not. If the type is not a primitive type, the returned string is just the exported name of the SchemaFieldType, and it is assumed that there will be a struct generated for this type in Types.
func getTypeDataFieldTypeFromSpecFieldType(fieldType spec.SchemaFieldType) (string, bool) {
	switch fieldType {
	case spec.SchemaFieldTypeString:
		return TypeStrString, true
	case spec.SchemaFieldTypeInteger:
		return TypeStrInteger, true
	case spec.SchemaFieldTypeDouble:
		return TypeStrDouble, true
	case spec.SchemaFieldTypeBoolean:
		return TypeStrBoolean, true
	case spec.SchemaFieldTypeFreeFormObject:
		return TypeStrFreeFormObject, true
	default:
		return exportedName(string(fieldType)), false
	}
}

func exportedName(name string) string {
	if len(name) == 0 {
		return name
	}
	return strings.ToUpper(string(name[0])) + name[1:]
}

func getOmitEmpty(required bool) string {
	if required {
		return ""
	}
	return ",omitempty"
}

func sortAuthMethodsByID(authMethods *[]AuthMethodData) {
	slices.SortFunc(*authMethods, func(a, b AuthMethodData) int {
		return strings.Compare(a.ID, b.ID)
	})
}

func sortTypeFieldsByName(fields *[]TypeFieldData) {
	slices.SortFunc(*fields, func(a, b TypeFieldData) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func sortResponsesByStatusCode(responses *[]ResponseData) {
	slices.SortFunc(*responses, func(a, b ResponseData) int {
		return a.StatusCode - b.StatusCode
	})
}

func sortTypesByName(types *[]TypeData) {
	slices.SortFunc(*types, func(a, b TypeData) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func sortParamsByName(params *[]ParamData) {
	slices.SortFunc(*params, func(a, b ParamData) int {
		return strings.Compare(a.Name, b.Name)
	})
}

func ExecuteTemplate(name string, tmplData any) ([]byte, error) {
	tmpl, err := template.ParseFS(goTemplates, "templates/*.tmpl", "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, name, tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}

func formatAndWriteFile(filePath string, content []byte) error {
	formattedContent, fmtErr := formatWithImports(content)
	if fmtErr == nil {
		content = formattedContent
	}
	err := utils.WriteFile(filePath, content)
	if err != nil {
		return fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	// handle formatting error after writing the file to ensure the file is created even if formatting fails
	// Helps in debugging issues without losing the generated code
	if fmtErr != nil {
		return fmt.Errorf("failed to format code (%s) with imports: %w", filePath, fmtErr)
	}
	return nil
}

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

// generateAndWriteHelperFuncsFile generates the content of the helperFuncs.go file, which contains helper functions that are used across multiple generated files.
func generateAndWriteHelperFuncsFile(packageName, clientName, version, filePath string) error {
	var buf bytes.Buffer
	tmpl, err := template.ParseFS(goTemplates, "templates/helperFuncsFile.tmpl")
	if err != nil {
		return fmt.Errorf("failed to parse helper funcs template: %w", err)
	}

	err = tmpl.ExecuteTemplate(&buf, "goHelperFuncsFile", map[string]string{
		"PackageName": packageName,
		"ClientName":  clientName,
		"Version":     version,
	})
	if err != nil {
		return fmt.Errorf("failed to execute helper funcs template: %w", err)
	}
	return formatAndWriteFile(filePath, buf.Bytes())
}

func IsTypeEnum(isPrimitive bool, typ string, schemas []*spec.Schema) bool {
	if !isPrimitive {
		// Check if the type is an enum by looking for a schema with the same name and checking if it has a non-empty Enum field
		for _, schema := range schemas {
			if schema.Name == typ && len(schema.Enum) > 0 {
				return true
			}
		}
	}
	return false
}
