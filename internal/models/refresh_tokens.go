package models

import "time"

type RefreshToken struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
}
