package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/traefik/whoamimcp/internal/keys"
	"github.com/traefik/whoamimcp/internal/prompts"
	"github.com/traefik/whoamimcp/internal/tools"
)

func NewMcpHandler(name string, options *mcp.StreamableHTTPOptions) *mcp.StreamableHTTPHandler {
	return mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		requestCtx := context.WithValue(req.Context(), keys.HttpRequest, req)
		requestCtx = context.WithValue(requestCtx, keys.Name, name)

		*req = *req.WithContext(requestCtx)

		options := mcp.ServerOptions{
			KeepAlive: time.Minute * 5,
		}
		server := mcp.NewServer(&mcp.Implementation{
			Name:    name,
			Version: "0.1.0",
		}, &options)

		mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "Say hi with hostname and server info"}, tools.SayHi)
		mcp.AddTool(server, &mcp.Tool{Name: "whoami", Description: "Get detailed system and request information"}, tools.Whoami)
		mcp.AddTool(server, &mcp.Tool{Name: "raise_error", Description: "Simulates an error on the MCP server"}, tools.RaiseError)
		server.AddPrompt(&mcp.Prompt{Name: "greet", Description: "Greeting prompt"}, prompts.PromptHi)

		return server
	}, options)
}
