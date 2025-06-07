package account

import "time"

type UserDoc struct {
	UserId    int64     `bson:"_id"`
	Username  string    `bson:"username"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
}
