package main

//http://localhost:8080/test

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request)  {
		fmt.Fprint(w, "znayu ya chto porvetsa pizdyukha!")
	})
	err := http.ListenAndServe(":8080", nil)
	if (err!=nil) {
		panic(err)
	}
}