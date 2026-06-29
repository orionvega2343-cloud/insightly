package models

import "time"

type Queries struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	FileId    int       `json:"file_id"`
	CreatedAt time.Time `json:"created_at"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
}
