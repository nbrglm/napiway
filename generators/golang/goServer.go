package golang

import (
	"fmt"
	"path/filepath"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
)

func GenerateServerHelpers(cfg *spec.GoServerGeneration, api *spec.Config) error {
	if absOutputDir, err := filepath.Abs(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to get absolute path of output directory: %w", err)
	} else {
		cfg.OutputDir = absOutputDir
	}

	// clear the content of the output directory
	if err := utils.ClearOutputDir(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to clear output directory: %w", err)
	}

	if err := generateAndWriteServerTypesFile(cfg, api); err != nil {
		return fmt.Errorf("failed to generate and write server types file: %w", err)
	}

	if err := generateAndWriteServerReqResFiles(cfg, api); err != nil {
		return fmt.Errorf("failed to generate and write server request and response files: %w", err)
	}

	helpersFilePath := filepath.Join(cfg.OutputDir, "helperFuncs.go")
	if err := generateAndWriteHelperFuncsFile(cfg.PackageName, api.Spec.ApiName, api.Spec.Version, helpersFilePath); err != nil {
		return fmt.Errorf("failed to generate and write helper functions file: %w", err)
	}

	return nil
}

func generateAndWriteServerTypesFile(cfg *spec.GoServerGeneration, api *spec.Config) error {
	types := TypesDataFromSpec(api)
	fileData := GoTypesFileData{
		PackageName: cfg.PackageName,
		Types:       types,
	}
	filePath := filepath.Join(cfg.OutputDir, "types.go")
	content, err := ExecuteTemplate("serverTypesFile", fileData)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	err = formatAndWriteFile(filePath, content)
	if err != nil {
		return err
	}
	return nil
}

func generateAndWriteServerReqResFiles(cfg *spec.GoServerGeneration, api *spec.Config) error {
	for idx := range api.Spec.Endpoints {
		filePath := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s.go", exportedName(api.Spec.Endpoints[idx].Name)))

		reqData, err := RequestResponsesDataFromEndpointDef(idx, &api.Spec)
		if err != nil {
			return fmt.Errorf("failed to get request and response data from endpoint (%s) definition: %w", api.Spec.Endpoints[idx].Name, err)
		}

		fileData := GoReqResFileData{
			PackageName: cfg.PackageName,
			RequestData: reqData,
		}
		content, err := ExecuteTemplate("serverReqResFile", fileData)
		if err != nil {
			return fmt.Errorf("failed to execute template for endpoint (%s): %w", api.Spec.Endpoints[idx].Name, err)
		}
		err = formatAndWriteFile(filePath, content)
		if err != nil {
			return err
		}
	}
	return nil
}
