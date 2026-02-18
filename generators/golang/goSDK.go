package golang

import (
	"bytes"
	"fmt"
	"path"
	"strings"
	"text/template"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
)

func GenerateGoSDK(cfg spec.GoSDKGeneration, api spec.Specification) error {
	utils.ClearOutputDir(cfg.OutputDir)

	runGoModInit(cfg.ModuleName, cfg.OutputDir)

	modulePaths := strings.Split(cfg.ModuleName, "/")
	packageName := modulePaths[len(modulePaths)-1]

	clientFilePath := path.Join(cfg.OutputDir, "client.go")
	fileData := &GoSdkClientAndTypesFileTemplateData{
		PackageName: packageName,
		ClientName:  exportedName(strings.ReplaceAll(api.ApiName, " ", "")),
		Endpoints:   []GoSdkEndpointDef{},
	}

	clientFileContent, err := createGoSDKClientFile(fileData, &api)
	if err != nil {
		return err
	}
	err = utils.WriteFile(clientFilePath, clientFileContent)
	if err != nil {
		return err
	}

	// types file
	typesFilePath := path.Join(cfg.OutputDir, "types.go")
	typesFileContent, err := createGoSDKTypesFile(fileData)
	if err != nil {
		// handle error
		return err
	}
	err = utils.WriteFile(typesFilePath, typesFileContent)
	if err != nil {
		// handle error
		return err
	}

	// helpers file
	helpersFileContent, err := generateGoHelperFuncsFile(packageName)
	if err != nil {
		return fmt.Errorf("failed to generate helper funcs file: %w", err)
	}
	helpersFilePath := path.Join(cfg.OutputDir, "helperFuncs.go")
	err = utils.WriteFile(helpersFilePath, helpersFileContent)
	if err != nil {
		return fmt.Errorf("failed to write helper funcs file: %w", err)
	}

	// run go mod tidy and goimports

	if err := runGoModTidy(cfg.OutputDir); err != nil {
		return err
	}

	return runGoImports(cfg.OutputDir)
}

func createGoSDKClientFile(data *GoSdkClientAndTypesFileTemplateData, api *spec.Specification) ([]byte, error) {
	for _, endpointName := range utils.SortedMapKeys(api.Endpoints) {
		endpoint := api.Endpoints[endpointName]

		sdkEndpoint := GoSdkEndpointDef{
			Name:      exportedName(endpointName),
			Responses: make(map[int]GoResponseStructDef, len(endpoint.Responses)),
		}

		// Collect required auth methods for the endpoint
		reqAuthMethodsAll, reqAuthMethodsAny := collectAnyAndAllAuthMethods(endpoint, api)

		sdkEndpoint.Request = collectRequest(endpoint, endpointName, reqAuthMethodsAll, reqAuthMethodsAny)

		responses := collectResponses(api.Endpoints[endpointName], endpointName)
		for _, resp := range responses {
			sdkEndpoint.Responses[resp.StatusCode] = resp
		}

		data.Endpoints = append(data.Endpoints, sdkEndpoint)
	}

	tmpl, err := template.ParseFS(goTemplates, "templates/*.tmpl", "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "goSdkClientFile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}

func createGoSDKTypesFile(data *GoSdkClientAndTypesFileTemplateData) ([]byte, error) {
	tmpl, err := template.ParseFS(goTemplates, "templates/*.tmpl", "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "goSdkTypesFile", data)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	return buf.Bytes(), nil
}
