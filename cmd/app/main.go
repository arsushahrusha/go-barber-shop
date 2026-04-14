package main

//http://localhost:8080/test

import (
	"fmt"
	"net/http"
)

func testHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "znayu ya chto porvetsa pizdyukha!")
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