package dto

import "time"

type BaseResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	BlankToken        string    `json:"-"`
	BlankTokenExpired time.Time `json:"-"`
	AccessToken       string    `json:"accessToken"`
	RefreshToken      string    `json:"-"`
	SessionId         string    `json:"-"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
