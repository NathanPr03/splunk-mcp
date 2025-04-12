package client

import (
	"context"
	"fmt"
	mcpClient "github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
)

func NewClient() *mcpClient.SSEMCPClient {
	baseUrl := "http://localhost:1690"

	sseMcpClient, err := mcpClient.NewSSEMCPClient(baseUrl + "/sse")
	if err != nil {
		fmt.Printf("Client error: %v\n", err)
		return nil
	}

	sseMcpClient.Start(context.Background())

	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name:    "nate-agent",
		Version: "1.0.0",
	}

	_, err = sseMcpClient.Initialize(context.Background(), initRequest)
	if err != nil {
		fmt.Printf("Initialize error: %v\n", err)
		return nil
	}

	return sseMcpClient
}
