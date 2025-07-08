package models

import "time"

type User struct {
	ID         string    `json:"id"`
	Username   string    `json:"username" gorm:"unique"`
	Password   string    `json:"password"`
	CreatedAt  time.Time `json:"created_at"`
	UpdateddAt time.Time `json:"updated_at"`
	DeletedAt  time.Time `json:"deleted_at"`
}

type Todo struct {
    ID        string `gorm:"primaryKey"`
    Title     string
    Completed bool
    UserID    string
    CreatedAt time.Time
}