package internal

import (
	"context"
	"net/http"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/rs/cors"
	"github.com/traefik/whoamimcp/internal/tools"
)

func StartServer(ctx context.Context, name string, addr string) error {

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		info := tools.GetWhoamiInfo(r, name)
		w.WriteHeader(http.StatusOK)
		info.WriteTo(w)
	})

	mux.Handle("/sse", NewMcpHandler(name, &mcp.StreamableHTTPOptions{}))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"*",
		},
		ExposedHeaders: []string{
			"Content-Length",
		},
		AllowCredentials: true,
		MaxAge:           -1,
	})

	srv := &http.Server{
		Addr:         addr,
		ReadTimeout:  0,
		WriteTimeout: 0,
		IdleTimeout:  5 * time.Minute,
		Handler:      c.Handler(mux),
	}
	return srv.ListenAndServe()
}
