package typescript

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
)

func GenerateTSSDK(genCfg spec.TsSDKGeneration, config *spec.Config) (err error) {
	if err := utils.ClearOutputDir(genCfg.OutputDir); err != nil {
		return fmt.Errorf("failed to clear output directory: %w", err)
	}

	// write package.json
	err = writePackageJsonFile(genCfg, config)
	if err != nil {
		return err
	}

	// write tsconfig.json
	err = writeTsConfigFile(genCfg)
	if err != nil {
		return err
	}

	// write gitignore
	err = writeGitIgnoreFile(genCfg)
	if err != nil {
		return err
	}

	// TODO: write README.md

	// gather data
	endpoints := make([]EndpointData, len(config.Spec.Endpoints))
	for idx := range config.Spec.Endpoints {
		reqData, err := RequestResponsesDataFromEndpointDef(idx, &config.Spec)
		if err != nil {
			return fmt.Errorf("failed to get request and response data from endpoint (%s) definition: %w", config.Spec.Endpoints[idx].Name, err)
		}
		endpoints[idx] = EndpointData{
			Name:    config.Spec.Endpoints[idx].Name,
			Request: reqData,
		}
	}

	clientName := exportedName(strings.ReplaceAll(config.Spec.ApiName, " ", ""))

	apiFileData := TsSdkApiFileData{
		ClientName:    clientName,
		ClientVersion: config.Spec.Version,
		Endpoints:     endpoints,
	}

	// write api.ts
	err = writeApiFile(apiFileData, config, genCfg)
	if err != nil {
		return err
	}

	requests := make([]RequestData, len(endpoints))
	for idx, endpoint := range endpoints {
		requests[idx] = endpoint.Request
	}
	types := TypesDataFromSpec(config)
	modelsFileData := TsSdkModelsFileData{
		ClientName: clientName,

		Types:    types,
		Requests: requests,
	}

	// write models.ts
	err = writeModelsFile(modelsFileData, genCfg)
	if err != nil {
		return err
	}

	// run npm i && npm run build
	err = utils.ExecCommand("npm install", genCfg.OutputDir)
	if err != nil {
		return err
	}

	err = utils.ExecCommand("npm run build", genCfg.OutputDir)
	if err != nil {
		return err
	}

	// copy license file if specified
	if genCfg.LicenseFile != nil {
		licenseFilePath := filepath.Join(genCfg.OutputDir, "LICENSE")
		bytes, err := os.ReadFile(*genCfg.LicenseFile)
		if err != nil {
			return fmt.Errorf("Error reading provided ts sdk license file %s: %w", *genCfg.LicenseFile, err)
		}
		err = utils.WriteFile(licenseFilePath, bytes)
		if err != nil {
			return fmt.Errorf("Error writing ts sdk license file %s: %w", licenseFilePath, err)
		}
	}
	return nil
}

func writeModelsFile(modelsFileData TsSdkModelsFileData, genCfg spec.TsSDKGeneration) error {
	modelsFileContent, err := createTsSDKModelsFile(&modelsFileData)
	if err != nil {
		return err
	}
	modelsFilePath := path.Join(genCfg.OutputDir, "models.ts")
	err = utils.WriteFile(modelsFilePath, modelsFileContent)
	if err != nil {
		return err
	}
	return nil
}

func writeApiFile(apiFileData TsSdkApiFileData, config *spec.Config, genCfg spec.TsSDKGeneration) error {
	apiFileContent, err := createTsSDKApiFile(&apiFileData, &config.Spec)
	if err != nil {
		return err
	}
	apiFilePath := path.Join(genCfg.OutputDir, "api.ts")
	err = utils.WriteFile(apiFilePath, apiFileContent)
	if err != nil {
		return err
	}
	return nil
}

func writeGitIgnoreFile(genCfg spec.TsSDKGeneration) error {
	gitignoreFilePath := path.Join(genCfg.OutputDir, ".gitignore")
	gitignoreFileContent, err := createTsSDKGitignoreFile()
	if err != nil {
		return err
	}
	err = utils.WriteFile(gitignoreFilePath, gitignoreFileContent)
	if err != nil {
		return err
	}
	return nil
}

func writeTsConfigFile(genCfg spec.TsSDKGeneration) error {
	tsconfigFilePath := path.Join(genCfg.OutputDir, "tsconfig.json")
	tsconfigFileContent, err := createTsSDKTsconfigFile()
	if err != nil {
		return err
	}
	err = utils.WriteFile(tsconfigFilePath, tsconfigFileContent)
	if err != nil {
		return err
	}
	return nil
}

func writePackageJsonFile(genCfg spec.TsSDKGeneration, config *spec.Config) error {
	packageJsonFileData := TsSdkPackageJsonTemplateData{
		PackageName: genCfg.PackageName,
		Version:     config.Spec.Version,
		Description: genCfg.Description,
		Website:     genCfg.Website,
		Repository:  genCfg.Repository,
		License:     genCfg.License,
		Author:      genCfg.Author,
		Keywords:    genCfg.Keywords,
	}
	packageJsonFilePath := path.Join(genCfg.OutputDir, "package.json")
	packageJsonFileContent, err := createTsSDKPackageJsonFile(&packageJsonFileData)
	if err != nil {
		return err
	}
	err = utils.WriteFile(packageJsonFilePath, packageJsonFileContent)
	if err != nil {
		return err
	}
	return nil
}

func createTsSDKModelsFile(templateData *TsSdkModelsFileData) ([]byte, error) {
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

func createTsSDKApiFile(data *TsSdkApiFileData, api *spec.Specification) ([]byte, error) {
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
	return strings.ToUpper(name[:1]) + name[1:]
}

func TypesDataFromSpec(specification *spec.Config) []TypeData {
	types := make([]TypeData, len(specification.Schemas))
	for idx, schema := range specification.Schemas {
		types[idx] = TypeData{
			Name:        exportedName(schema.Name),
			Description: schema.Description,
			Fields:      getFieldsDataFromSpecFields(schema.Properties),
		}
	}
	sortTypesByName(&types)
	return types
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

	var reqBodyName *string
	if endpoint.BodyName != nil {
		reqBodyNameVal := exportedName(*endpoint.BodyName)
		reqBodyName = &reqBodyNameVal
	}

	responses := make([]ResponseData, len(endpoint.Responses))
	has413 := false
	for i, resp := range endpoint.Responses {
		if resp.StatusCode == 413 {
			has413 = true
		}
		var respBodyName *string
		if resp.BodyName != nil {
			respBodyNameVal := exportedName(*resp.BodyName)
			respBodyName = &respBodyNameVal
		}
		responses[i] = ResponseData{
			StatusCode:       resp.StatusCode,
			Name:             exportedName(endpoint.Name + strconv.Itoa(resp.StatusCode)),
			Description:      resp.Description,
			RawBody:          resp.RawBody,
			ContentType:      *resp.ContentType,
			Headers:          mapSpecParamToParamData(resp.Headers),
			ResponseBodyName: respBodyName,
		}
	}

	// If the responses don't include a 413 (Payload Too Large), add a default one for the case when the request body exceeds MaxBodyBytes
	if !has413 && !endpoint.RawBody && reqBodyName != nil {
		// If RawBody is true, it means the generated code will not be reading/unmarshalling the request body, so we don't need to worry about adding a 413 response.
		// If requestBodyName is nil, it means there is no request body, so we also don't need to worry about adding a 413 response.
		desc := "Payload Too Large - the request body exceeds the maximum allowed size"
		responses = append(responses, ResponseData{
			StatusCode:       413,
			Name:             exportedName(endpoint.Name + "413"),
			Description:      &desc,
			RawBody:          false, // no raw body too
			Headers:          nil,
			ResponseBodyName: nil, // no body
		})
	}

	sortResponsesByStatusCode(&responses)

	return RequestData{
		Name:            exportedName(endpoint.Name + "Req"),
		Description:     endpoint.Description,
		Method:          string(endpoint.Method),
		Path:            endpoint.Path,
		ContentType:     *endpoint.ContentType,
		RawBody:         endpoint.RawBody,
		RequestBodyName: reqBodyName,
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

func getFieldsDataFromSpecFields(fields []*spec.SchemaField) []TypeFieldData {
	fieldsData := make([]TypeFieldData, len(fields))
	for idx, field := range fields {
		typ, _ := getTypeDataFieldTypeFromSpecFieldType(field.Type)
		fieldsData[idx] = TypeFieldData{
			Name:        exportedName(field.Name),
			Description: field.Description,
			Type:        typ,
			IsArray:     field.IsArray,
			Required:    field.Required,
			NonEmpty:    field.NonEmpty,
		}
	}

	sortTypeFieldsByName(&fieldsData)
	return fieldsData
}

// Returns the TS type string for a given SchemaFieldType, and a boolean indicating whether the type is a primitive type or not
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
