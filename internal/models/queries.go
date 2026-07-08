package models

import "time"

type Queries struct {
	Id        int       `json:"id" db:"id"`
	UserId    int       `json:"user_id" db:"user_id"`
	FileId    int       `json:"file_id" db:"file_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Question  string    `json:"question" db:"question"`
	Answer    string    `json:"answer" db:"answer"`
}
