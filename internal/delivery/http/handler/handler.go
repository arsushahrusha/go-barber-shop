package handler

import (
	"fmt"
	"io"
	usecasetest "my-go-server/internal/usecase/test"
	dbservice "my-go-server/internal/usecase/db"
	"net/http"
	"time"
)

type Handler struct {
	service usecasetest.MessageService
	dbService dbservice.DBServiceMessage
}

func NewHandler(service usecasetest.MessageService, dbService dbservice.DBServiceMessage) *Handler{
	return &Handler{
		service: service,
		dbService: dbService,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	someWork()
	fmt.Fprint(w, h.service.GetMessage())
}

func (h *Handler) HandleDBTest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method is not allowed", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		http.Error(w, "Failed to read request Body", http.StatusBadRequest)
		return
	}

	value := string(body)

	if value == "" {
		http.Error(w, "Empty body", http.StatusBadRequest)
		return
	}

	_, err = h.dbService.Save(value)
	if err != nil {
		http.Error(w, "Failed to save value", http.StatusInternalServerError)
		return
	}
	
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("saved"))
}

func someWork() {
	fmt.Println("Background work started...")
	time.Sleep(5*time.Second)
	fmt.Println("Background work finished...")
}