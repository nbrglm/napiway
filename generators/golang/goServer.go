package golang

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"
	"text/template"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
)

func GenerateServerHelpers(cfg spec.GoServerGeneration, api spec.Specification) error {
	if absOutputDir, err := filepath.Abs(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to get absolute path of output directory: %w", err)
	} else {
		cfg.OutputDir = absOutputDir
	}

	// clear the content of the output directory
	if err := utils.ClearOutputDir(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to clear output directory: %w", err)
	}

	// generate a file for each endpoint
	for _, endpointName := range utils.SortedMapKeys(api.Endpoints) {
		filePath := path.Join(cfg.OutputDir, exportedName(endpointName)+"Handler.go")
		content, err := GenerateGoServerHandlerFile(cfg, &api, endpointName)
		if err != nil {
			// handle error
			return fmt.Errorf("failed to generate handler for %s: %w", endpointName, err)
		}
		err = utils.WriteFile(filePath, content)
		if err != nil {
			// handle error
			return fmt.Errorf("failed to write file %s: %w", filePath, err)
		}
	}

	helpersFileContent, err := generateGoHelperFuncsFile(cfg.PackageName)
	if err != nil {
		return fmt.Errorf("failed to generate helper funcs file: %w", err)
	}
	helpersFilePath := path.Join(cfg.OutputDir, "helperFuncs.go")
	err = utils.WriteFile(helpersFilePath, helpersFileContent)
	if err != nil {
		return fmt.Errorf("failed to write helper funcs file: %w", err)
	}

	// run fmt and goimports
	if err := runGoImports(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to run goimports: %w", err)
	}
	return nil
}

func GenerateGoServerHandlerFile(cfg spec.GoServerGeneration, api *spec.Specification, endpointName string) ([]byte, error) {
	endpoint, exists := api.Endpoints[endpointName]
	if !exists {
		return nil, fmt.Errorf("endpoint %s not found in specification", endpointName)
	}

	// Collect required auth methods for the endpoint
	reqAuthMethodsAll, reqAuthMethodsAny := collectAnyAndAllAuthMethods(endpoint, api)

	request := collectRequest(endpoint, endpointName, reqAuthMethodsAll, reqAuthMethodsAny)
	responses := collectResponses(endpoint, endpointName)

	tmplData := GoServerFileTemplateData{
		PackageName: cfg.PackageName,
		Request:     request,
		Responses:   responses,
	}

	tmpl, err := template.ParseFS(goTemplates, "templates/*.tmpl", "templates/**/*.tmpl")
	if err != nil {
		return nil, fmt.Errorf("failed to parse templates: %w", err)
	}

	var buf bytes.Buffer
	err = tmpl.ExecuteTemplate(&buf, "goServerHandlerFile", tmplData)
	if err != nil {
		return nil, fmt.Errorf("failed to execute template: %w", err)
	}
	// return the generated content
	return buf.Bytes(), nil
}
