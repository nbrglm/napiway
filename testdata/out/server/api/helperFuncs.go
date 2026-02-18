package api

import (
	"fmt"
	"strconv"
	"strings"
)

func parsefloat64Param(param string, paramName string, required bool) (*float64, error) {
	param = strings.TrimSpace(param)
	if param == "" {
		if required {
			return nil, fmt.Errorf("missing required parameter '%s'", paramName)
		}
		return nil, nil
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(param), 64)
	if err != nil {
		return nil, fmt.Errorf("invalid number parameter '%s': %v", paramName, err)
	}

	return &value, nil
}

func parseboolParam(param string, paramName string, required bool) (*bool, error) {
	param = strings.TrimSpace(param)
	if param == "" {
		if required {
			return nil, fmt.Errorf("missing required parameter '%s'", paramName)
		}
		return nil, nil
	}

	value, err := strconv.ParseBool(param)
	if err != nil {
		return nil, fmt.Errorf("invalid boolean parameter '%s': %v", paramName, err)
	}

	return &value, nil
}

func parsestringParam(param string, paramName string, required bool) (*string, error) {
	param = strings.TrimSpace(param)
	if param == "" {
		if required {
			return nil, fmt.Errorf("missing required parameter '%s'", paramName)
		}
		return nil, nil
	}
	return &param, nil
}

func paramToString(param interface{}, paramName string, goType string, required bool) (string, error) {
	if param == nil {
		if required {
			return "", fmt.Errorf("missing required parameter '%s'", paramName)
		}
		return "", nil
	}
	strValue := ""
	switch goType {
	case "string":
		strValue, _ = param.(string)
	case "*string":
		ptrValue, _ := param.(*string)
		if ptrValue != nil {
			strValue = *ptrValue
		}
	case "float64":
		floatValue, ok := param.(float64)
		if !ok {
			return "", fmt.Errorf("invalid float64 parameter '%s'", paramName)
		}
		strValue = fmt.Sprintf("%f", floatValue)
	case "*float64":
		ptrValue, ok := param.(*float64)
		if !ok {
			return "", fmt.Errorf("invalid *float64 parameter '%s'", paramName)
		}
		if ptrValue != nil {
			strValue = fmt.Sprintf("%f", *ptrValue)
		}
	case "bool":
		boolValue, ok := param.(bool)
		if !ok {
			return "", fmt.Errorf("invalid bool parameter '%s'", paramName)
		}
		strValue = fmt.Sprintf("%t", boolValue)
	case "*bool":
		ptrValue, ok := param.(*bool)
		if !ok {
			return "", fmt.Errorf("invalid *bool parameter '%s'", paramName)
		}
		if ptrValue != nil {
			strValue = fmt.Sprintf("%t", *ptrValue)
		}
	default:
		return "", fmt.Errorf("unsupported goType '%s' for parameter '%s'", goType, paramName)
	}
	if strValue == "" && required {
		return "", fmt.Errorf("missing required parameter '%s'", paramName)
	}
	return strValue, nil
}
