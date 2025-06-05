package users

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	users_bans "server/internal/account/users/bans"
	"server/internal/account/users/jwts"
	"server/pkg/database"
	httpServer "server/pkg/web/http"
	"strings"
	"time"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Token   string `json:"token,omitempty"`
}

type LoginHandler struct{}

func (h *LoginHandler) Path() string {
	return "/user/login"
}

func (h *LoginHandler) Method() string {
	return "POST"
}

func (h *LoginHandler) Handle(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var req LoginRequest
	encoder := json.NewEncoder(w)
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(LoginResponse{Code: 1, Message: "Invalid request"})
		return
	}
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	if req.Username == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		_ = encoder.Encode(LoginResponse{Code: 1, Message: "Missing fields"})
		return
	}

	conn, err := database.GetMongoConnector()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(LoginResponse{Code: 1, Message: "Database connection error"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{}
	switch {
	case isEmail(req.Username):
		filter["email"] = req.Username
	case isNumeric(req.Username):
		filter["userId"] = toInt64(req.Username)
	default:
		filter["username"] = req.Username
	}

	var user UserDoc // 使用统一的 UserDoc
	err = conn.DB.Collection("users").FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.WriteHeader(http.StatusUnauthorized)
			_ = encoder.Encode(LoginResponse{Code: 1, Message: "User not found"})
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			_ = encoder.Encode(LoginResponse{Code: 1, Message: "Database error"})
		}
		return
	}

	// 只允许 userId
	userID := int(user.UserId)

	// 检查用户是否被封禁
	banned, banInfo, err := users_bans.IsUserBanned(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(LoginResponse{Code: 4, Message: "Ban check failed"})
		return
	}
	if banned {
		msg := "User is banned"
		if banInfo != nil {
			if banInfo.BanReason != "" {
				msg += ": " + banInfo.BanReason
			}
			if banInfo.BanEnd != nil {
				msg += " (Until: " + banInfo.BanEnd.Format("2006-01-02 15:04:05") + ")"
			}
		}
		_ = encoder.Encode(LoginResponse{Code: 5, Message: msg})
		return
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)) != nil {
		w.WriteHeader(http.StatusUnauthorized)
		_ = encoder.Encode(LoginResponse{Code: 2, Message: "Incorrect password"})
		return
	}

	token, err := jwts.GenerateJWT(userID, 72*time.Hour)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_ = encoder.Encode(LoginResponse{Code: 3, Message: "Token generation failed"})
		return
	}

	_ = encoder.Encode(LoginResponse{Code: 0, Message: "Login success", Token: token})
}

// 判断是否为邮箱（与 register.go 保持一致，使用正则）
func isEmail(s string) bool {
	emailRegexp := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegexp.MatchString(s)
}

// 判断是否为纯数字
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// 字符串转 int64
func toInt64(s string) int64 {
	var n int64
	for i := 0; i < len(s); i++ {
		n = n*10 + int64(s[i]-'0')
	}
	return n
}

func init() {
	handler := &LoginHandler{}
	httpServer.HttpManagerInstance.RegisterHandler(handler)
}
