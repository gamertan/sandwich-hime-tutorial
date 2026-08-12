// SPDX-License-Identifier: 0BSD

package main

import (
	"log"
	"net/http"

	"gitea.speelman.ca/gamertan/sandwich-hime-tutorial/internal/server"
)

func main() {
	const address = "127.0.0.1:8080"
	log.Printf("listening on http://%s", address)
	log.Fatal(http.ListenAndServe(address, server.New()))
}
