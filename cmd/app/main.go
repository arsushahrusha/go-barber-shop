package main

//http://localhost:8080/test

import (
	"fmt"
	"net/http"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	var repository MessageRepository
	repository = &TestRepository{}
	var service MessageService
	service = &TestService{repository: repository}
	
	fmt.Fprint(w, service.GetMessage())
}

type MessageRepository interface {
	GetMessage() string
}

type TestRepository struct{}
	func (r *TestRepository) GetMessage() string {
		return "znayu ya chto porvetsa pizdyukha!"
	}

type MessageService interface {
	GetMessage() string
}

type TestService struct{
	repository MessageRepository
}
	func (s *TestService) GetMessage() string {
		return s.repository.GetMessage()
	}

func setupRoutes() {
	http.HandleFunc("/test", testHandler)
}

func main() {
	setupRoutes()

	err := http.ListenAndServe(":8080", nil)
	if err!=nil {
		panic(err)
	}
}