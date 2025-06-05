package web

import (
	"fmt"
	"log"
	"net/http"
	"goauthx/pkg/config"
)

func StartServer(handler http.Handler) error {
	cfg := config.GetConfig()
	addr := fmt.Sprintf(":%d", cfg.HTTPServer.Port)
	if cfg.HTTPServer.EnableSSL {
		log.Printf("Starting HTTPS server on %s\n", addr)
		return http.ListenAndServeTLS(
			addr,
			cfg.HTTPServer.SSLCertFile,
			cfg.HTTPServer.SSLKeyFile,
			handler,
		)
	}
	log.Printf("Starting HTTP server on %s\n", addr)
	return http.ListenAndServe(addr, handler)
}