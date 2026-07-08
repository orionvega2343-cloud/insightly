package models

type Files struct {
	Id     int    `json:"id" db:"id"`
	UserId int    `json:"user_id" db:"user_id"`
	Name   string `json:"filename" db:"filename"`
	Path   string `json:"filepath" db:"filepath"`
}
