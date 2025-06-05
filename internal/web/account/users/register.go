package users

import (
	"encoding/json"
	"net/http"
	httpServer "server/pkg/web/http"
)

// RegisterRequest represents the structure of the registration request
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Captcha  string `json:"captcha"` // 新增验证码字段
}

// RegisterResponse represents the structure of the registration response
type RegisterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// RegisterHandler handles the user registration requests
type RegisterHandler struct{}

// Path returns the HTTP path for the registration handler
func (h *RegisterHandler) Path() string {
	return "/user/register"
}

// Method returns the HTTP method for the registration handler
func (h *RegisterHandler) Method() string {
	return "POST"
}

// Handle processes the registration request
func (h *RegisterHandler) Handle(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(RegisterResponse{Code: 1, Message: "Invalid request"})
		return
	}

	resp, status := RegisterUser(&req, true)
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(resp)
}

func init() {
	handler := &RegisterHandler{}
	httpServer.HttpManagerInstance.RegisterHandler(handler)
}
