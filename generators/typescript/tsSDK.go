package typescript

import (
	"bytes"
	"fmt"
	"path"
	"slices"
	"strings"
	"text/template"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
)

func GenerateTSSDK(cfg spec.TsSDKGeneration, api spec.Specification) (err error) {
	utils.ClearOutputDir(cfg.OutputDir)

	// write package.json
	packageJsonFileData := TsSdkPackageJsonTemplateData{
		PackageName: cfg.PackageName,
		Version:     api.Version,
		Description: cfg.Description,
		Website:     cfg.Website,
		Repository:  cfg.Repository,
		License:     cfg.License,
		Author:      cfg.Author,
		Keywords:    cfg.Keywords,
	}
	packageJsonFilePath := path.Join(cfg.OutputDir, "package.json")
	packageJsonFileContent, err := createTsSDKPackageJsonFile(&packageJsonFileData)
	if err != nil {
		return err
	}
	err = utils.WriteFile(packageJsonFilePath, packageJsonFileContent)
	if err != nil {
		return err
	}

	// write tsconfig.json
	tsconfigFilePath := path.Join(cfg.OutputDir, "tsconfig.json")
	tsconfigFileContent, err := createTsSDKTsconfigFile()
	if err != nil {
		return err
	}
	err = utils.WriteFile(tsconfigFilePath, tsconfigFileContent)
	if err != nil {
		return err
	}

	// write gitignore
	gitignoreFilePath := path.Join(cfg.OutputDir, ".gitignore")
	gitignoreFileContent, err := createTsSDKGitignoreFile()
	if err != nil {
		return err
	}
	err = utils.WriteFile(gitignoreFilePath, gitignoreFileContent)
	if err != nil {
		return err
	}

	// TODO: write README.md

	templateData := TsSdkApiAndModelsFilesTemplateData{
		ClientName:    exportedName(strings.ReplaceAll(api.ApiName, " ", "")),
		ClientVersion: api.Version,
		Endpoints:     []TsSdkEndpointDef{},
	}

	// write api.ts
	apiFileContent, err := createTsSDKApiFile(&templateData, &api)
	if err != nil {
		return err
	}
	apiFilePath := path.Join(cfg.OutputDir, "api.ts")
	err = utils.WriteFile(apiFilePath, apiFileContent)
	if err != nil {
		return err
	}

	// write models.ts
	modelsFileContent, err := createTsSDKModelsFile(&templateData)
	if err != nil {
		return err
	}
	modelsFilePath := path.Join(cfg.OutputDir, "models.ts")
	err = utils.WriteFile(modelsFilePath, modelsFileContent)
	if err != nil {
		return err
	}

	// run npm i && npm run build
	err = utils.ExecCommand("npm install", cfg.OutputDir)
	if err != nil {
		return err
	}

	err = utils.ExecCommand("npm run build", cfg.OutputDir)
	if err != nil {
		return err
	}
	return nil
}

func createTsSDKModelsFile(templateData *TsSdkApiAndModelsFilesTemplateData) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFS(tsTemplates, "templates/*")
	if err != nil {
		return nil, err
	}
	tmpl = tmpl.Option("missingkey=error")
	err = tmpl.ExecuteTemplate(&buf, "models.ts", templateData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func createTsSDKApiFile(data *TsSdkApiAndModelsFilesTemplateData, api *spec.Specification) ([]byte, error) {

	for _, endpointName := range utils.SortedMapKeys(api.Endpoints) {
		endpoint := api.Endpoints[endpointName]

		sdkEndpoint := TsSdkEndpointDef{
			Name:      exportedName(endpointName),
			Responses: make(map[int]TsResponseDef),
		}

		// Collect request authentication methods
		reqAuthMethodsAll, reqAuthMethodsAny := collectAnyAndAllAuthMethods(endpoint, api)

		sdkEndpoint.Request = collectRequest(endpoint, endpointName, reqAuthMethodsAll, reqAuthMethodsAny)

		responses := collectResponses(endpoint, endpointName)
		for _, response := range responses {
			sdkEndpoint.Responses[response.StatusCode] = response
		}

		data.Endpoints = append(data.Endpoints, sdkEndpoint)
	}

	var buf bytes.Buffer
	tmpl, err := template.ParseFS(tsTemplates, "templates/*")
	if err != nil {
		return nil, err
	}
	err = tmpl.ExecuteTemplate(&buf, "api.ts", data)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func collectResponses(endpoint *spec.Endpoint, endpointName string) []TsResponseDef {
	responses := []TsResponseDef{}
	for _, status := range utils.SortedMapKeys(endpoint.Responses) {
		response := endpoint.Responses[status]
		respName := fmt.Sprintf("%s%dResponse", exportedName(endpointName), status)
		types := []TsTypeDef{}
		var respBodyName *string
		if response.Body != nil && response.ContentType != nil && *response.ContentType == "application/json" {
			name := fmt.Sprintf("%s%dResponseBody", exportedName(endpointName), status)
			respBodyName = &name
			types = append(types, collectTypesFromResponseBody(respName, response.Body)...)
		}
		headers := []TsParamDef{}
		for _, header := range response.Headers {
			headers = append(headers, getTsParamDef(header))
		}
		responses = append(responses, TsResponseDef{
			StatusCode:       status,
			Name:             respName,
			Description:      response.Description,
			ContentType:      *response.ContentType,
			Headers:          headers,
			SupportingTypes:  types,
			ResponseBodyName: respBodyName,
		})
	}
	slices.SortFunc(responses, func(a, b TsResponseDef) int {
		return strings.Compare(a.Name, b.Name)
	})
	return responses
}

func collectTypesFromResponseBody(respName string, body *spec.HTTPBody) []TsTypeDef {
	types := []TsTypeDef{}

	field := spec.Field{
		Type:        spec.FieldTypeObject,
		Description: body.Description,
		Properties:  body.Properties,
		Required:    true, // response body is always required
	}

	collectTypes(respName, "Body", &field, &types)
	return types
}

func collectRequest(endpoint *spec.Endpoint, endpointName string, reqAuthMethodsAll, reqAuthMethodsAny []TsAuthMethodDef) TsRequestDef {
	headers, query, path := []TsParamDef{}, []TsParamDef{}, []TsParamDef{}
	for _, header := range endpoint.Headers {
		headers = append(headers, getTsParamDef(header))
	}
	for _, queryParam := range endpoint.QueryParams {
		query = append(query, getTsParamDef(queryParam))
	}
	for _, pathParam := range endpoint.PathParams {
		path = append(path, getTsParamDef(pathParam))
	}

	requestTypes := []TsTypeDef{}
	var requestBodyName *string

	// collect types from RequestBody
	if endpoint.RequestBody != nil && endpoint.ContentType != nil && *endpoint.ContentType == "application/json" {
		name := fmt.Sprintf("%sRequestBody", exportedName(endpointName))
		requestBodyName = &name
		requestTypes = append(requestTypes, collectTypesFromRequestBody(exportedName(endpointName), endpoint.RequestBody)...)
	}

	slices.SortFunc(requestTypes, func(a, b TsTypeDef) int {
		return strings.Compare(a.Name, b.Name)
	})

	return TsRequestDef{
		Name:            exportedName(endpointName) + "Request",
		Description:     endpoint.Description,
		ContentType:     *endpoint.ContentType,
		Method:          string(endpoint.Method),
		Path:            endpoint.Path,
		HeaderParams:    headers,
		QueryParams:     query,
		PathParams:      path,
		SupportingTypes: requestTypes,
		RequestBodyName: requestBodyName,
		AuthAll:         reqAuthMethodsAll,
		AuthAny:         reqAuthMethodsAny,
	}
}

func collectTypesFromRequestBody(parentName string, body *spec.HTTPBody) []TsTypeDef {
	types := []TsTypeDef{}

	field := spec.Field{
		Type:        spec.FieldTypeObject,
		Description: body.Description,
		Properties:  body.Properties,
		Required:    true, // request body is always required
	}

	collectTypes(parentName, "RequestBody", &field, &types)
	return types
}

func collectTypes(parentName, fieldName string, field *spec.Field, types *[]TsTypeDef) {
	if field.Type != spec.FieldTypeObject {
		return
	}

	typeName := parentName + exportedName(fieldName)

	typeDef := TsTypeDef{
		Name:   typeName,
		Fields: []TsFieldDef{},
	}

	for _, propName := range utils.SortedMapKeys(field.Properties) {
		prop := field.Properties[propName]
		tsFieldType, shouldRecurseValidate := tsTypeFromSpecField(prop, typeName+exportedName(propName))

		tsElemType := ""
		if prop.Type == spec.FieldTypeArray {
			tsElemType = strings.TrimPrefix(tsFieldType, "Array<")
			tsElemType = strings.TrimSuffix(tsElemType, ">")
		}

		fieldDef := TsFieldDef{
			Name:            exportedName(propName),
			Description:     prop.Description,
			Type:            tsFieldType,
			RecurseValidate: shouldRecurseValidate,
			IsArray:         prop.Type == spec.FieldTypeArray,
			ElemType:        tsElemType,
			Required:        prop.Required,
			NonEmpty:        prop.NonEmpty,
		}
		typeDef.Fields = append(typeDef.Fields, fieldDef)

		switch prop.Type {
		case spec.FieldTypeObject:
			collectTypes(typeName, exportedName(propName), &prop, types)
		case spec.FieldTypeArray:
			if prop.Items != nil && prop.Items.Type == spec.FieldTypeObject {
				collectTypes(typeName, exportedName(propName)+"Item", prop.Items, types)
			}
		}
	}
	*types = append(*types, typeDef)
}

func tsTypeFromSpecField(field spec.Field, parentName string) (typ string, shouldRecurseValidate bool) {
	switch field.Type {
	case spec.FieldTypeString:
		return "string", false
	case spec.FieldTypeNumber:
		return "number", false
	case spec.FieldTypeBoolean:
		return "boolean", false
	case spec.FieldTypeObject:
		return exportedName(parentName), true
	case spec.FieldTypeArray:
		elemType, elemRecurse := tsTypeFromSpecField(*field.Items, parentName+"Item")
		return fmt.Sprintf("Array<%s>", elemType), elemRecurse
	default:
		panic("unsupported field type: " + string(field.Type))
	}
}

func getTsParamDef(param spec.Param) TsParamDef {
	return TsParamDef{
		Name:          exportedName(param.Name),
		TransportName: param.TransportName,
		Type:          tsTypeFromParam(param),
		Description:   param.Description,
		Required:      param.Required,
	}
}

func tsTypeFromParam(param spec.Param) string {
	switch param.Type {
	case spec.ParamTypeString:
		return "string"
	case spec.ParamTypeNumber:
		return "number"
	case spec.ParamTypeBoolean:
		return "boolean"
	default:
		panic("unsupported param type: " + string(param.Type))
	}
}

func createTsSDKPackageJsonFile(packageJsonFileData *TsSdkPackageJsonTemplateData) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFS(tsTemplates, "templates/*")
	if err != nil {
		return nil, err
	}
	err = tmpl.ExecuteTemplate(&buf, "package.json", packageJsonFileData)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func createTsSDKTsconfigFile() ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFS(tsTemplates, "templates/*")
	if err != nil {
		return nil, err
	}
	err = tmpl.ExecuteTemplate(&buf, "tsconfig.json", nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func createTsSDKGitignoreFile() ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFS(tsTemplates, "templates/*")
	if err != nil {
		return nil, err
	}
	err = tmpl.ExecuteTemplate(&buf, ".gitignore", nil)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func exportedName(name string) string {
	// return cases.Title(language.English).String(name)
	return strings.ToUpper(name[:1]) + name[1:]
}

func collectAnyAndAllAuthMethods(endpoint *spec.Endpoint, api *spec.Specification) (reqAuthMethodsAll, reqAuthMethodsAny []TsAuthMethodDef) {
	if endpoint.Auth != nil {
		for _, auth := range api.Auth {
			if slices.Contains(endpoint.Auth.All, auth.ID) {
				reqAuthMethodsAll = append(reqAuthMethodsAll, TsAuthMethodDef{
					ID:            auth.ID,
					Name:          auth.Name,
					Type:          getAuthMethodType(auth.Type),
					TransportName: auth.TransportName,
					Description:   auth.Description,
					Format:        auth.Format,
				})
				continue
			}

			if slices.Contains(endpoint.Auth.Any, auth.ID) {
				reqAuthMethodsAny = append(reqAuthMethodsAny, TsAuthMethodDef{
					ID:            auth.ID,
					Name:          auth.Name,
					Type:          getAuthMethodType(auth.Type),
					TransportName: auth.TransportName,
					Description:   auth.Description,
					Format:        auth.Format,
				})
			}
		}
	}
	return
}

func getAuthMethodType(authType spec.AuthMethodType) TsAuthMethodType {
	switch authType {
	case spec.AuthMethodHeader:
		return TsAuthMethodTypeHeader
	default:
		panic("unsupported auth method type: " + string(authType))
	}
}
