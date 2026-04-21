package handler

import (
	"fmt"
	"net/http"
	usecasetest "my-go-server/internal/usecase/test"
)

type Handler struct {
	service usecasetest.MessageService
}

func NewHandler(service usecasetest.MessageService) *Handler{
	return &Handler{service: service}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	
	fmt.Fprint(w, h.service.GetMessage())
}