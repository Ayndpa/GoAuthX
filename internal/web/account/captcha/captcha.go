package captcha

import (
	"encoding/json"
	"fmt"
	"github.com/patrickmn/go-cache"
	"goauthx/internal/config"
	"goauthx/internal/smtp"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	// 使用 go-cache 作为内存验证码存储，带TTL
	captchaCache = cache.New(5*time.Minute, 10*time.Minute)
	rnd          = rand.New(rand.NewSource(time.Now().UnixNano()))
	// 邮箱和IP限速缓存，TTL为1分钟
	rateLimitEmailCache = cache.New(1*time.Minute, 2*time.Minute)
	rateLimitIPCache    = cache.New(1*time.Minute, 2*time.Minute)
)

type CaptchaRequest struct {
	Email string `json:"email"`
}

type CaptchaResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// 生成6位数字验证码
func generateCaptchaCode() string {
	return fmt.Sprintf("%06d", rnd.Intn(1000000))
}

func HandleCaptcha(w http.ResponseWriter, r *http.Request) {
	var req CaptchaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp := CaptchaResponse{Code: 1, Message: "Invalid request"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		resp := CaptchaResponse{Code: 1, Message: "Missing email"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// 限速逻辑
	clientIP := r.Header.Get("X-Real-IP")
	if clientIP == "" {
		clientIP = r.Header.Get("X-Forwarded-For")
	}
	if clientIP == "" {
		clientIP = strings.Split(r.RemoteAddr, ":")[0]
	}
	if _, found := rateLimitEmailCache.Get(req.Email); found {
		w.WriteHeader(http.StatusTooManyRequests)
		resp := CaptchaResponse{Code: 3, Message: "Too many requests for this email, please try again later"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}
	if _, found := rateLimitIPCache.Get(clientIP); found {
		w.WriteHeader(http.StatusTooManyRequests)
		resp := CaptchaResponse{Code: 3, Message: "Too many requests from this IP, please try again later"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	code := generateCaptchaCode()
	captchaCache.Set(req.Email, code, 10*time.Minute)

	templatePath := "./resources/template/email/captcha.html"
	htmlBytes, err := os.ReadFile(templatePath)
	if err != nil {
		captchaCache.Delete(req.Email)
		w.WriteHeader(http.StatusInternalServerError)
		resp := CaptchaResponse{Code: 2, Message: "Failed to load email template"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	cfg := config.GetConfig()
	htmlBody := strings.ReplaceAll(string(htmlBytes), "{{CODE}}", code)
	htmlBody = strings.ReplaceAll(htmlBody, "{{NAME}}", cfg.Name)
	subject := fmt.Sprintf("您的 %s 验证码", cfg.Name)
	if err := email.SendEmail([]string{req.Email}, subject, htmlBody); err != nil {
		captchaCache.Delete(req.Email)
		w.WriteHeader(http.StatusInternalServerError)
		resp := CaptchaResponse{Code: 2, Message: "Failed to send email"}
		_ = json.NewEncoder(w).Encode(resp)
		return
	}

	// 新增：设置邮箱和IP的限速缓存
	rateLimitEmailCache.Set(req.Email, true, cache.DefaultExpiration)
	rateLimitIPCache.Set(clientIP, true, cache.DefaultExpiration)

	resp := CaptchaResponse{Code: 0, Message: "Captcha sent"}
	_ = json.NewEncoder(w).Encode(resp)
}

// 验证验证码是否正确，并在成功后删除
func VerifyCaptcha(email, code string) bool {
	val, found := captchaCache.Get(email)
	if !found {
		return false
	}
	stored, ok := val.(string)
	if ok && stored == code {
		captchaCache.Delete(email) // 验证成功后删除
		return true
	}
	return false
}
