package account

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"goauthx/internal/db"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// 包级变量，避免每次请求都编译正则
var usernamePattern = regexp.MustCompile(`^[a-z0-9_]+$`)

// 辅助函数，批量去除空格
func trimRegisterRequest(req *RegisterRequest) {
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.Email = strings.TrimSpace(req.Email)
	req.Captcha = strings.TrimSpace(req.Captcha)
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	Captcha  string `json:"captcha"`
}

type RegisterResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 注册核心逻辑，供 HTTP handler 和命令复用
// 不再处理验证码校验
func RegisterUser(req *RegisterRequest) (RegisterResponse, int) {
	trimRegisterRequest(req)
	req.Username = strings.ToLower(req.Username)

	if req.Username == "" || req.Password == "" || req.Email == "" {
		return RegisterResponse{Code: 1, Message: "Missing fields"}, http.StatusBadRequest
	}
	if !usernamePattern.MatchString(req.Username) {
		return RegisterResponse{Code: 1, Message: "Username must be lowercase letters, numbers, or underscores"}, http.StatusBadRequest
	}

	conn, err := db.GetMongoConnector()
	if err != nil {
		return RegisterResponse{Code: 2, Message: "Database connection error"}, http.StatusInternalServerError
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{"username": req.Username},
			{"email": req.Email},
		},
	}
	count, err := conn.DB.Collection("users").CountDocuments(ctx, filter)
	if err != nil {
		return RegisterResponse{Code: 2, Message: "Database error"}, http.StatusInternalServerError
	}
	if count > 0 {
		return RegisterResponse{Code: 1, Message: "Username or email already exists"}, http.StatusConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return RegisterResponse{Code: 2, Message: "Password encryption failed"}, http.StatusInternalServerError
	}

	userId, err := db.GetNextSequenceValue("user_id")
	if err != nil {
		return RegisterResponse{Code: 2, Message: "Failed to generate userId"}, http.StatusInternalServerError
	}

	userDoc := UserDoc{
		UserId:    userId,
		Username:  req.Username,
		Password:  string(hashedPassword),
		Email:     req.Email,
		CreatedAt: time.Now(),
	}

	_, err = conn.DB.Collection("users").InsertOne(ctx, userDoc)
	if err != nil {
		return RegisterResponse{Code: 2, Message: "Register failed"}, http.StatusInternalServerError
	}

	return RegisterResponse{Code: 0, Message: "Register success"}, http.StatusOK
}
