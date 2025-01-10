package dto

import "time"

type BaseResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	BlankToken        string    `json:"-"`
	BlankTokenExpired time.Time `json:"-"`
	SessionId         string    `json:"-"`
	ID                string    `json:"id"`
	FullName          string    `json:"fullName"`
	Email             string    `json:"email"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
