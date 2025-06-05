package account

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server/internal/account/users/jwts"
	"server/pkg/database"
	"time"
)

type UserBan struct {
	BanID     primitive.ObjectID `bson:"_id,omitempty"`
	UserID    int                `bson:"user_id"`
	BannedBy  int                `bson:"banned_by,omitempty"`
	BanReason string             `bson:"ban_reason,omitempty"`
	BanStart  time.Time          `bson:"ban_start_time,omitempty"`
	BanEnd    *time.Time         `bson:"ban_end_time,omitempty"`
	IsActive  bool               `bson:"is_active"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// IsUserBanned checks if the user is currently banned (is_active=1 and ban_end_time is null or in the future)
func IsUserBanned(userID int) (bool, *UserBan, error) {
	conn, err := database.GetMongoConnector()
	if err != nil {
		return false, nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
		"$or": []bson.M{
			{"ban_end_time": bson.M{"$eq": nil}},
			{"ban_end_time": bson.M{"$gt": time.Now()}},
		},
	}
	findOpts := options.FindOne().SetSort(bson.D{{"ban_start_time", -1}})
	var ban UserBan
	err = conn.DB.Collection("users_bans").FindOne(ctx, filter, findOpts).Decode(&ban)
	if err != nil {
		if err.Error() == "mongo: no documents in result" {
			return false, nil, nil
		}
		return false, nil, err
	}
	return true, &ban, nil
}

// BanUser inserts a new ban record for the user
func BanUser(userID int, bannedBy *int, reason string, banEnd time.Time) error {
	conn, err := database.GetMongoConnector()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	now := time.Now()
	banDoc := bson.M{
		"user_id":        userID,
		"banned_by":      bannedBy,
		"ban_reason":     reason,
		"ban_end_time":   banEnd,
		"is_active":      true,
		"created_at":     now,
		"updated_at":     now,
		"ban_start_time": now,
	}
	_, err = conn.DB.Collection("users_bans").InsertOne(ctx, banDoc)
	if err == nil {
		jwts.RemoveUserJWTsFromWhitelist(userID)
	}
	return err
}

// UnbanUser sets is_active=false for all active bans of the user
func UnbanUser(userID int) error {
	conn, err := database.GetMongoConnector()
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	filter := bson.M{
		"user_id":   userID,
		"is_active": true,
	}
	update := bson.M{
		"$set": bson.M{
			"is_active":  false,
			"updated_at": time.Now(),
		},
	}
	_, err = conn.DB.Collection("users_bans").UpdateMany(ctx, filter, update)
	return err
}
