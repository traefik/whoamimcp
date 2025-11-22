package tools

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/traefik/whoamimcp/internal/keys"
)

type WhoamiInfo struct {
	Hostname   string      `json:"hostname,omitempty"`
	Name       string      `json:"name,omitempty"`
	Headers    http.Header `json:"headers,omitempty"`
	URL        string      `json:"url,omitempty"`
	Scheme     string      `json:"scheme,omitempty"`
	Host       string      `json:"host,omitempty"`
	Port       uint32      `json:"port,omitempty"`
	Method     string      `json:"method,omitempty"`
	RemoteAddr string      `json:"remoteAddr,omitempty"`
	UserAgent  string      `json:"userAgent,omitempty"`
	Proto      string      `json:"proto,omitempty"`
	TLS        bool        `json:"tls,omitempty"`
}

func (i *WhoamiInfo) WriteTo(w io.Writer) (int64, error) {
	var total int
	var err error
	write := func(format string, args ...interface{}) {
		if err != nil {
			return
		}
		var n int
		n, err = fmt.Fprintf(w, format, args...)
		total += n
	}
	write("Name: %s\n", i.Name)
	write("Hostname: %s\n", i.Hostname)
	write("RemoteAddr: %s\n", i.RemoteAddr)
	write("Method: %s\n", i.Method)
	write("URL: %s\n", i.URL)
	write("Proto: %s\n", i.Proto)
	write("Scheme: %s\n", i.Scheme)
	write("Host: %s\n", i.Host)
	write("Port: %d\n", i.Port)
	write("User-Agent: %s\n", i.UserAgent)
	write("TLS: %v\n", i.TLS)
	write("Headers:\n")
	for name, values := range i.Headers {
		for _, value := range values {
			write("%s: %s\n", name, value)
		}
	}
	return int64(total), err
}

func (i *WhoamiInfo) String() string {
	var output strings.Builder
	i.WriteTo(&output)
	return output.String()
}

func GetWhoamiInfo(req *http.Request, name string) WhoamiInfo {
	hostname, _ := os.Hostname()
	port, _ := strconv.ParseUint(req.URL.Port(), 10, 32)
	return WhoamiInfo{
		Hostname:   hostname,
		Name:       name,
		URL:        req.URL.RequestURI(),
		Scheme:     req.URL.Scheme,
		Host:       req.Host,
		Port:       uint32(port),
		Method:     req.Method,
		RemoteAddr: req.RemoteAddr,
		UserAgent:  req.UserAgent(),
		Proto:      req.Proto,
		TLS:        req.TLS != nil,
		Headers:    req.Header,
	}
}

type WhoamiArgs struct {
}

func Whoami(ctx context.Context, ss *mcp.ServerSession, params *mcp.CallToolParamsFor[WhoamiArgs]) (*mcp.CallToolResultFor[WhoamiInfo], error) {

	name, ok := ctx.Value(keys.Name).(string)
	if !ok {
		return &mcp.CallToolResultFor[WhoamiInfo]{}, fmt.Errorf("Internal error: failed to get name")
	}
	req, ok := ctx.Value(keys.HttpRequest).(*http.Request)
	if !ok {
		return &mcp.CallToolResultFor[WhoamiInfo]{}, fmt.Errorf("Internal error: failed to get request")
	}

	info := GetWhoamiInfo(req, name)
	return &mcp.CallToolResultFor[WhoamiInfo]{
		Content:           []mcp.Content{},
		StructuredContent: info,
	}, nil
}
