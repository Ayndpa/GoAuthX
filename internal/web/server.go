package web

import (
	"fmt"
	"goauthx/internal/config"
	"log"
	"net/http"
)

import (
	"goauthx/internal/web/account/captcha"
	"goauthx/internal/web/account/users"
)

func StartServer() error {
	cfg := config.GetConfig()
	addr := fmt.Sprintf(":%d", cfg.HTTPServer.Port)

	http.HandleFunc("/captcha", captcha.HandleCaptcha)
	http.HandleFunc("/login", users.HandleLogin)
	http.HandleFunc("/register", users.HandleRegister)

	if cfg.HTTPServer.EnableSSL {
		log.Printf("Starting HTTPS server on %s\n", addr)
		return http.ListenAndServeTLS(
			addr,
			cfg.HTTPServer.SSLCertFile,
			cfg.HTTPServer.SSLKeyFile,
			nil,
		)
	}
	log.Printf("Starting HTTP server on %s\n", addr)
	return http.ListenAndServe(addr, nil)
}
