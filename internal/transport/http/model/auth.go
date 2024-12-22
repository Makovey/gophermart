package model

type AuthRequest struct {
	Login    string `json:"login" validate:"required,max=14"`
	Password string `json:"password" validate:"required,max=32"`
}
