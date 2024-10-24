package models

import (
    "github.com/google/uuid"
)

type User struct {
    ID uuid.UUID `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Credits int `json:"credits"`
    Role int `json:"role"`
    Status int `json:"status"`
    Country string `json:"country"`
    Rating int `json:"rating"`
    Avatar string `json:"avatar"`
    CreationDate string `json:"creationdate"`
}

type LoginRequest struct {
    ID uuid.UUID `json:"id"`
    Username string `json:"username"`
    Password string `json:"password"`
    Status int `json:"status"`
    Role int `json:"role"`
}

type BanRequest struct {
    Username string `json:"username"`
    Reason string `json:"reason"`
}
