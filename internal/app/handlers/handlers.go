package handlers

import "gorm.io/gorm"

type Handler struct {
	DB *gorm.DB
}

type JsonResponse struct {
	Data       interface{} `json:"data"`
	Pagination Pagination  `json:"pagination"`
}

type Pagination struct {
	NextCursor string `json:"nextCursor"`
	PrevCursor string `json:"prevCursor"`
}

func NewHandler(db *gorm.DB) *Handler {
	return &Handler{
		DB: db,
	}
}
