package prompts

import (
	"context"
	"fmt"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func PromptHi(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	hostname, _ := os.Hostname()
	return &mcp.GetPromptResult{
		Description: "Code review prompt",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: fmt.Sprintf("Say hi to %s. Regards, %s", params.Arguments["name"], hostname)}},
		},
	}, nil
}
