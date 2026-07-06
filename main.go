package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var addr string

func main() {
	flag.StringVar(&addr, "addr", ":80", "http service address")
	flag.Parse()

	log.Fatal(startServer())
}

type HiArgs struct {
	Name string `json:"name" jsonschema:"the name to say hi to"`
}

func SayHi(_ context.Context, _ *mcp.CallToolRequest, args HiArgs) (*mcp.CallToolResult, any, error) {
	name, _ := os.Hostname()
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Hello, %s from %s %s!", args.Name, name, addr)},
		},
	}, nil, nil
}

func PromptHi(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Code review prompt",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: fmt.Sprintf("Say hi to %s from %s", req.Params.Arguments["name"], addr)}},
		},
	}, nil
}

// newHandler builds the MCP streamable HTTP handler serving the greeter server.
func newHandler() http.Handler {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter_s1"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	server.AddPrompt(&mcp.Prompt{Name: "greet"}, PromptHi)

	return mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
}

func startServer() error {
	handler := newHandler()
	log.Printf("MCP Server handler listening at %s", addr)

	return http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(rw, req)
	}))
}
