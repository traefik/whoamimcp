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

	ctx := context.Background()
	log.Fatal(startServer(ctx))
}

type HiArgs struct {
	Name string `json:"name" jsonschema:"the name to say hi to"`
}

func SayHi(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[HiArgs]) (*mcp.CallToolResultFor[struct{}], error) {
	name, _ := os.Hostname()
	return &mcp.CallToolResultFor[struct{}]{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Hello, %s from %s %s!", params.Arguments.Name, name, addr)},
		},
	}, nil
}

func PromptHi(ctx context.Context, ss *mcp.ServerSession, params *mcp.GetPromptParams) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Code review prompt",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: fmt.Sprintf("Say hi to %s from %s", params.Arguments["name"], addr)}},
		},
	}, nil
}

func startServer(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter_s1"}, nil)
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	server.AddPrompt(&mcp.Prompt{Name: "greet"}, PromptHi)

	// server.AddResource(&mcp.Resource{
	// 	Name:     "info",
	// 	MIMEType: "text/plain",
	// 	URI:      "embedded:info",
	// }, handleEmbeddedResource)

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, nil)
	log.Printf("MCP Server handler listening at %s", addr)

	return http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(rw, req)
	}))
}
