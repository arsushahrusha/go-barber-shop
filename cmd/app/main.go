package main

//http://localhost:8080/test

import (
	"fmt"
	"net/http"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	repository := &TestRepository{}
	service := &TestService{repository: *repository}
	
	fmt.Fprint(w, service.GetMessage())
}

type TestRepository struct{}
	func (s *TestRepository) GetMessage() string {
		return "znayu ya chto porvetsa pizdyukha!"
	}

type TestService struct{
	repository TestRepository
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