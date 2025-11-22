package tools

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type ErrorArgs struct {
	ShouldFail bool   `json:"should_fail" jsonschema:"should the request fail"`
	Message    string `json:"error_message" jsonschema:"what message should say"`
}

func RaiseError(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[ErrorArgs]) (*mcp.CallToolResultFor[any], error) {
	if params.Arguments.ShouldFail {
		return nil, fmt.Errorf("%s", params.Arguments.Message)
	}
	return &mcp.CallToolResultFor[any]{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: params.Arguments.Message,
			},
		},
	}, nil
}
