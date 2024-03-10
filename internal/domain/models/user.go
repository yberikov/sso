package models

import "time"

type User struct {
	ID        int64
	Email     string
	PassHash  []byte
	Name      string
	Telephone string
	DateBirth time.Time
}
