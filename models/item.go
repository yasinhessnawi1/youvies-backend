package models

import "time"

type User struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	Active    bool      `json:"active"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Avatar    string    `json:"avatar"`
	Favorites []string  `json:"favorites"`
	Friends   []string  `json:"friends"`
	Rooms     []string  `json:"rooms"`
	Watched   []string  `json:"watched"`
}
type UserResponse struct {
	ID        string    `json:"id,omitempty"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Active    bool      `json:"active"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	Avatar    string    `json:"avatar"`
	Favorites []string  `json:"favorites"`
	Friends   []string  `json:"friends"`
	Rooms     []string  `json:"rooms"`
	Watched   []string  `json:"watched"`
}
