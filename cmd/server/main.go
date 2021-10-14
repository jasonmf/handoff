package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	listen := ":8000"
	if v := os.Getenv("LISTEN"); v != "" {
		listen = v
	}

	assets := "/assets"
	if v := os.Getenv("ASSETS"); v != "" {
		assets = v
	}

	log.Printf("starting %s on %s", assets, listen)
	err := http.ListenAndServe(listen, http.FileServer(http.Dir(assets)))
	if err != nil {
		fmt.Println("Failed to start server", err)
		return
	}
}
