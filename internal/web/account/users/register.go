package users

import (
	"encoding/json"
	account "goauthx/internal/account"
	"goauthx/internal/web/account/captcha"
	"net/http"
)

// Handle processes the registration request
func HandleRegister(w http.ResponseWriter, r *http.Request) {
	var req account.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(account.RegisterResponse{Code: 1, Message: "Invalid request"})
		return
	}

	// 验证码校验逻辑移到这里
	if req.Captcha == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(account.RegisterResponse{Code: 1, Message: "Captcha required"})
		return
	}
	if !captcha.VerifyCaptcha(req.Email, req.Captcha) {
		w.WriteHeader(http.StatusUnauthorized)
		_ = json.NewEncoder(w).Encode(account.RegisterResponse{Code: 4, Message: "Invalid or expired captcha"})
		return
	}

	resp, status := account.RegisterUser(&req)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}
