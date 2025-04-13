package tools

import (
	"errors"
	"fmt"
	"github.com/mark3labs/mcp-go/mcp"
)

func GetParamFromRequest(request mcp.CallToolRequest, paramToSearchFor string) (string, error) {
	param, ok := request.Params.Arguments[paramToSearchFor]
	if !ok {
		return "", errors.New("parameter not found")
	}

	paramStr, ok := param.(string)
	if !ok {
		return "", errors.New("parameter can not be converted to string")
	}

	return paramStr, nil
}

func GetBoolParamFromRequest(request mcp.CallToolRequest, paramName string) (bool, error) {
	param, ok := request.Params.Arguments[paramName]
	if !ok {
		return false, fmt.Errorf("parameter %s not provided", paramName)
	}

	boolVal, ok := param.(bool)
	if !ok {
		return false, fmt.Errorf("parameter %s is not a boolean", paramName)
	}

	return boolVal, nil
}

func GetIntParamFromRequest(request mcp.CallToolRequest, paramName string) (int, error) {
	param, ok := request.Params.Arguments[paramName]
	if !ok {
		return 0, fmt.Errorf("parameter %s not provided", paramName)
	}
	intVal, ok := param.(int)
	if !ok {
		return 0, fmt.Errorf("parameter %s is not an integer", paramName)
	}

	return intVal, nil
}
