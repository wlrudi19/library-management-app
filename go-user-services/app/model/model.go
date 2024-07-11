package model

import "time"

type UserResponse struct {
	Id       int       `json:"id"`
	Name     string    `json:"name"`
	Created  time.Time `json:"created"`
	Password string    `json:"password"`
}

type UserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Id          int    `json:"id"`
	AccessToken string `json:"access-token"`
}
