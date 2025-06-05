package jwts

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/pkg/config"
	"server/pkg/database"
	"time"
)

var jwtSecret = []byte(config.GetConfig().JWTSecret)

type Claims struct {
	UserID int    `json:"user_id"`
	JTI    string `json:"jti"`
	jwt.RegisteredClaims
}

// JWTRecord 用于MongoDB存储
type JWTRecord struct {
	UserID    int       `bson:"user_id"`
	JTI       string    `bson:"jti"`
	ExpiresAt time.Time `bson:"expires_at"`
}

// getJWTCollection 获取 users_jwts 集合
func getJWTCollection() (*mongo.Collection, error) {
	conn, err := database.GetMongoConnector()
	if err != nil {
		return nil, err
	}
	coll := conn.DB.Collection("users_jwts")
	// 创建TTL索引（如果不存在）
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "expires_at", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(0),
	}
	_, _ = coll.Indexes().CreateOne(context.Background(), indexModel)
	return coll, nil
}

// GenerateJWT 签发JWT并存入MongoDB
func GenerateJWT(userID int, duration time.Duration) (string, error) {
	jti := uuid.NewString()
	expireAt := time.Now().Add(duration)
	claims := Claims{
		UserID: userID,
		JTI:    jti,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expireAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        jti,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	coll, err := getJWTCollection()
	if err != nil {
		return "", err
	}
	_, err = coll.InsertOne(context.Background(), JWTRecord{
		UserID:    userID,
		JTI:       jti,
		ExpiresAt: expireAt,
	})
	if err != nil {
		return "", err
	}
	return signed, nil
}

// ParseJWT 验证JWT，校验MongoDB白名单并滑动续期
func ParseJWT(tokenString string) (bool, *Claims) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return false, nil
		}
		return jwtSecret, nil
	})
	if err != nil {
		return false, nil
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return false, nil
	}
	coll, err := getJWTCollection()
	if err != nil {
		return false, nil
	}
	var record JWTRecord
	err = coll.FindOne(context.Background(), bson.M{"jti": claims.JTI, "user_id": claims.UserID}).Decode(&record)
	if err != nil {
		return false, nil
	}
	// 检查是否过期
	if time.Now().After(record.ExpiresAt) {
		// 已过期，自动由TTL清理
		return false, nil
	}
	// 滑动续期：如果距离过期小于一半，则延长
	ttl := record.ExpiresAt.Sub(time.Now())
	origTTL := claims.ExpiresAt.Time.Sub(claims.IssuedAt.Time)
	if ttl < origTTL/2 {
		newExpire := time.Now().Add(origTTL)
		_, _ = coll.UpdateOne(context.Background(),
			bson.M{"jti": claims.JTI, "user_id": claims.UserID},
			bson.M{"$set": bson.M{"expires_at": newExpire}},
		)
	}
	return true, claims
}

// RemoveJWTFromWhitelist 移除指定 jti（强制下线单个会话）
func RemoveJWTFromWhitelist(jti string) {
	coll, err := getJWTCollection()
	if err != nil {
		return
	}
	_, _ = coll.DeleteOne(context.Background(), bson.M{"jti": jti})
}

// RemoveUserJWTsFromWhitelist 移除指定用户的所有jti（强制下线该用户所有会话）
func RemoveUserJWTsFromWhitelist(userID int) {
	coll, err := getJWTCollection()
	if err != nil {
		return
	}
	_, _ = coll.DeleteMany(context.Background(), bson.M{"user_id": userID})
}
