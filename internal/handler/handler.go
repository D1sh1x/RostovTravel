package handler

import "hackaton/internal/service"

type Handler struct {
	service service.ServiceInterface
}

func NewHandler(s service.ServiceInterface) *Handler {
	return &Handler{
		service: s,
	}
}
