package handler

import (
	"fmt"
	usecasetest "my-go-server/internal/usecase/test"
	"net/http"
	"time"
)

type Handler struct {
	service usecasetest.MessageService
}

func NewHandler(service usecasetest.MessageService) *Handler{
	return &Handler{service: service}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	someWork()
	fmt.Fprint(w, h.service.GetMessage())
}

func someWork() {
	fmt.Println("Background work started...")
	time.Sleep(5*time.Second)
	fmt.Println("Background work finished...")
}