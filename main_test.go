package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// TestGreetToolCall exercises the full MCP client flow with the standalone SSE
// stream enabled (the default). On an older go-sdk server that never flushed the
// GET SSE response, Connect would hang here; it must succeed and greet must work.
func TestGreetToolCall(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "v1.0.0"}, nil)
	// DisableStandaloneSSE is intentionally left false to cover the SSE listen path.
	transport := &mcp.StreamableClientTransport{Endpoint: srv.URL}

	session, err := client.Connect(ctx, transport, nil)
	if err != nil {
		t.Fatalf("connect: %v", err)
	}
	defer session.Close() // nolint

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("list tools: %v", err)
	}
	if len(tools.Tools) != 1 || tools.Tools[0].Name != "greet" {
		t.Fatalf("expected a single %q tool, got %+v", "greet", tools.Tools)
	}

	res, err := session.CallTool(ctx, &mcp.CallToolParams{
		Name:      "greet",
		Arguments: map[string]any{"name": "you"},
	})
	if err != nil {
		t.Fatalf("call tool: %v", err)
	}
	if res.IsError {
		t.Fatalf("tool returned an error result: %+v", res)
	}
	if len(res.Content) == 0 {
		t.Fatal("empty tool result content")
	}
	text, ok := res.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("expected text content, got %T", res.Content[0])
	}
	if !strings.Contains(text.Text, "Hello, you") {
		t.Fatalf("unexpected greeting: %q", text.Text)
	}
}

// TestStandaloneSSEFlushesImmediately asserts that a GET to the endpoint opens the
// standalone SSE stream and returns its headers right away, rather than blocking
// until the first server-sent byte. A server that buffers headers makes this GET
// hang, and the short client timeout below turns that into a failure.
func TestStandaloneSSEFlushesImmediately(t *testing.T) {
	srv := httptest.NewServer(newHandler())
	defer srv.Close()

	sessionID := initSession(t, srv.URL)

	getReq, err := http.NewRequest(http.MethodGet, srv.URL, nil)
	if err != nil {
		t.Fatalf("new GET request: %v", err)
	}
	getReq.Header.Set("Accept", "text/event-stream")
	getReq.Header.Set("Mcp-Session-Id", sessionID)
	getReq.Header.Set("Mcp-Protocol-Version", "2025-06-18")

	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Do(getReq)
	if err != nil {
		t.Fatalf("standalone SSE GET did not return headers in time (regression): %v", err)
	}
	defer resp.Body.Close() // nolint

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %q", ct)
	}
}

// initSession performs the MCP initialize handshake over raw HTTP and returns the
// assigned session ID.
func initSession(t *testing.T, baseURL string) string {
	t.Helper()

	const initBody = `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","capabilities":{},"clientInfo":{"name":"t","version":"1"}}}`
	req, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(initBody))
	if err != nil {
		t.Fatalf("new initialize request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/event-stream")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("initialize: %v", err)
	}
	defer resp.Body.Close() // nolint

	sessionID := resp.Header.Get("Mcp-Session-Id")
	if sessionID == "" {
		t.Fatal("no Mcp-Session-Id returned by initialize")
	}

	notifyBody := `{"jsonrpc":"2.0","method":"notifications/initialized"}`
	notifyReq, err := http.NewRequest(http.MethodPost, baseURL, strings.NewReader(notifyBody))
	if err != nil {
		t.Fatalf("new initialized notification: %v", err)
	}
	notifyReq.Header.Set("Content-Type", "application/json")
	notifyReq.Header.Set("Accept", "application/json, text/event-stream")
	notifyReq.Header.Set("Mcp-Session-Id", sessionID)
	notifyReq.Header.Set("Mcp-Protocol-Version", "2025-06-18")

	notifyResp, err := client.Do(notifyReq)
	if err != nil {
		t.Fatalf("initialized notification: %v", err)
	}
	notifyResp.Body.Close() // nolint

	return sessionID
}
