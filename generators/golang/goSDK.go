package golang

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/nbrglm/napiway/spec"
	"github.com/nbrglm/napiway/utils"
)

func GenerateGoSDK(cfg *spec.GoSDKGeneration, spc *spec.Specification) error {
	if absOutputDir, err := filepath.Abs(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to get absolute path of output directory: %w", err)
	} else {
		cfg.OutputDir = absOutputDir
	}

	// clear the content of the output directory
	if err := utils.ClearOutputDir(cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to clear output directory: %w", err)
	}

	if err := runGoModInit(cfg.ModuleName, cfg.OutputDir); err != nil {
		return fmt.Errorf("failed to run go mod init: %w", err)
	}

	modulePaths := strings.Split(cfg.ModuleName, "/")
	packageName := modulePaths[len(modulePaths)-1]

	// Types file
	if err := generateAndWriteSdkTypesFile(cfg, spc, packageName); err != nil {
		return fmt.Errorf("failed to generate and write types file: %w", err)
	}

	// ReqRes Files
	if err := generateAndWriteSdkReqResFiles(cfg, spc, packageName); err != nil {
		return fmt.Errorf("failed to generate and write request and response files: %w", err)
	}

	// helpers file
	helpersFilePath := filepath.Join(cfg.OutputDir, "helperFuncs.go")
	if err := generateAndWriteHelperFuncsFile(packageName, spc.ApiName, spc.Version, helpersFilePath); err != nil {
		return fmt.Errorf("failed to generate and write helper functions file: %w", err)
	}

	// client file
	clientFileEndpoints := make([]EndpointData, len(spc.Endpoints))
	for idx := range spc.Endpoints {
		reqData, err := RequestResponsesDataFromEndpointDef(idx, spc)
		if err != nil {
			return fmt.Errorf("failed to get request and response data from endpoint (%s) definition: %w", spc.Endpoints[idx].Name, err)
		}
		clientFileEndpoints[idx] = EndpointData{
			Name:    spc.Endpoints[idx].Name,
			Request: reqData,
		}
	}
	clientFilePath := filepath.Join(cfg.OutputDir, "client.go")
	clientFileData := GoSdkClientFileData{
		PackageName:   packageName,
		ClientName:    exportedName(strings.ReplaceAll(spc.ApiName, " ", "")),
		ClientVersion: spc.Version,
		Endpoints:     clientFileEndpoints,
	}
	clientFileContent, err := ExecuteTemplate("sdkClientFile", clientFileData)
	if err != nil {
		return fmt.Errorf("failed to execute client file template: %w", err)
	}
	clientFileContentFormatted, formatErr := formatWithImports(clientFileContent)
	if formatErr == nil {
		clientFileContent = clientFileContentFormatted
	}
	err = utils.WriteFile(clientFilePath, clientFileContent)
	if err != nil {
		return fmt.Errorf("failed to write client file %s: %w", clientFilePath, err)
	}
	if formatErr != nil {
		return fmt.Errorf("failed to format client file %s: %w", clientFilePath, formatErr)
	}

	// Write the License File, if any is provided
	if cfg.LicenseFile != nil {
		licenseFilePath := filepath.Join(cfg.OutputDir, "LICENSE")
		bytes, err := os.ReadFile(*cfg.LicenseFile)
		if err != nil {
			return fmt.Errorf("Error reading provided go sdk license file %s: %w", *cfg.LicenseFile, err)
		}
		err = utils.WriteFile(licenseFilePath, bytes)
		if err != nil {
			return fmt.Errorf("Error writing go sdk license file %s: %w", licenseFilePath, err)
		}
	}

	if err := runGoModTidy(cfg.OutputDir); err != nil {
		return err
	}
	return nil
}

func generateAndWriteSdkTypesFile(cfg *spec.GoSDKGeneration, spc *spec.Specification, packageName string) error {
	types := TypesDataFromSpec(spc)
	authMethodsAllDefined := AuthMethodsFromSpec(spc)
	fileData := GoTypesFileData{
		PackageName: packageName,
		Types:       types,
		AuthMethods: authMethodsAllDefined,
	}
	filePath := filepath.Join(cfg.OutputDir, "types.go")
	content, err := ExecuteTemplate("sdkTypesFile", fileData)
	if err != nil {
		return fmt.Errorf("failed to execute template: %w", err)
	}
	err = formatAndWriteFile(filePath, content)
	if err != nil {
		return err
	}
	return nil
}

func generateAndWriteSdkReqResFiles(cfg *spec.GoSDKGeneration, spc *spec.Specification, packageName string) error {
	for idx := range spc.Endpoints {
		filePath := filepath.Join(cfg.OutputDir, fmt.Sprintf("%s.go", exportedName(spc.Endpoints[idx].Name)))

		reqData, err := RequestResponsesDataFromEndpointDef(idx, spc)
		if err != nil {
			return fmt.Errorf("failed to get request and response data from endpoint (%s) definition: %w", spc.Endpoints[idx].Name, err)
		}

		fileData := GoReqResFileData{
			PackageName: packageName,
			RequestData: reqData,
		}
		content, err := ExecuteTemplate("sdkReqResFile", fileData)
		if err != nil {
			return fmt.Errorf("failed to execute template for endpoint (%s): %w", spc.Endpoints[idx].Name, err)
		}
		err = formatAndWriteFile(filePath, content)
		if err != nil {
			return err
		}
	}
	return nil
}
