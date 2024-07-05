package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"password"`
	Email     string             `bson:"email" json:"email"`
	Role      string             `bson:"role" json:"role"`
	Active    bool               `bson:"active" json:"active"`
	Created   time.Time          `bson:"created" json:"created"`
	Updated   time.Time          `bson:"updated" json:"updated"`
	Avatar    string             `bson:"avatar" json:"avatar"`
	Favorites []string           `bson:"favorites" json:"favorites"`
	Friends   []string           `bson:"friends" json:"friends"`
	Rooms     []string           `bson:"rooms" json:"rooms"`
	Watched   []string           `bson:"watched" json:"watched"`
}
