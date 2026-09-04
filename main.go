package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var addr string
var legacy bool

func main() {
	flag.StringVar(&addr, "addr", ":80", "http service address")
	flag.BoolVar(&legacy, "legacy", false, "use legacy mode")
	flag.Parse()

	ctx := context.Background()
	log.Fatal(startServer(ctx))
}

type HiArgs struct {
	Name string `json:"name" jsonschema:"the name to say hi to"`
}

func SayHi(ctx context.Context, req *mcp.CallToolRequest, args HiArgs) (
	*mcp.CallToolResult,
	any,
	error) {
	name, _ := os.Hostname()
	return &mcp.CallToolResult{

		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Hello, %s from %s %s!", args.Name, name, addr)},
		},
	}, nil, nil
}

func SayHiComplexArray(ctx context.Context, req *mcp.CallToolRequest, args HiArgs) (
	*mcp.CallToolResult,
	any,
	error) {
	name, _ := os.Hostname()
	return &mcp.CallToolResult{}, []struct{ Name string }{{args.Name}, {addr}, {name}}, nil
}

func SayHiComplexString(ctx context.Context, req *mcp.CallToolRequest, args HiArgs) (
	*mcp.CallToolResult,
	any,
	error) {
	name, _ := os.Hostname()
	return &mcp.CallToolResult{}, fmt.Sprintf("Hello, %s from %s %s!", args.Name, name, addr), nil
}

func SayHiInputRequest(ctx context.Context, req *mcp.CallToolRequest, args HiArgs) (
	*mcp.CallToolResult,
	any,
	error) {
	if len(req.Params.InputResponses) == 0 {
		return &mcp.CallToolResult{
			InputRequests: mcp.InputRequestMap{
				"name": &mcp.ElicitParams{
					Message: "What's your test name?",
					RequestedSchema: &jsonschema.Schema{
						Type: "object",
						Properties: map[string]*jsonschema.Schema{
							"name": {Type: "string"},
						},
					},
				}}}, nil, nil
	} else {
		return &mcp.CallToolResult{}, fmt.Sprintf("Hello, %s from %s %s!", args.Name, req.Params.InputResponses["name"].(*mcp.ElicitResult).Content["name"], addr), nil
	}
}

func PromptHi(ctx context.Context, params *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	return &mcp.GetPromptResult{
		Description: "Code review prompt",
		Messages: []*mcp.PromptMessage{
			{Role: "user", Content: &mcp.TextContent{Text: fmt.Sprintf("Say hi to %s from %s", params.Params.Arguments["name"], addr)}},
		},
	}, nil
}

func startServer(ctx context.Context) error {
	server := mcp.NewServer(&mcp.Implementation{Name: "greeter_s1"}, &mcp.ServerOptions{})
	mcp.AddTool(server, &mcp.Tool{Name: "greet", Description: "say hi"}, SayHi)
	mcp.AddTool(server, &mcp.Tool{Name: "greet_complex_array", Description: "say hi with complex array"}, SayHiComplexArray)
	mcp.AddTool(server, &mcp.Tool{Name: "greet_complex_string", Description: "say hi with complex string"}, SayHiComplexString)
	mcp.AddTool(server, &mcp.Tool{Name: "greet_inputrequest", Description: "say hi with inputRequest"}, SayHiInputRequest)
	server.AddPrompt(&mcp.Prompt{Name: "greet"}, PromptHi)
	server.AddResource(&mcp.Resource{
		Name:     "info",
		MIMEType: "text/plain",
		URI:      "embedded:info",
	}, func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		return &mcp.ReadResourceResult{Contents: []*mcp.ResourceContents{
			{
				Text: fmt.Sprintf("Say hi to %s from %s", req.Params.URI, addr),
			},
		}}, nil
	})

	handler := mcp.NewStreamableHTTPHandler(func(*http.Request) *mcp.Server {
		return server
	}, &mcp.StreamableHTTPOptions{
		Stateless: !legacy,
	})
	log.Printf("MCP Server handler listening at %s", addr)

	return http.ListenAndServe(addr, http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler.ServeHTTP(rw, req)
	}))
}
