package main

import (
	"net/http"
	"fmt"
)

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "20")
}