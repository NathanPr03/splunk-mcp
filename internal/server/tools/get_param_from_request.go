package tools

import (
	"errors"
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
