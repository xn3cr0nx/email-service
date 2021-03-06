// +build docs

package main

import (
	"fmt"
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
	_ "github.com/xn3cr0nx/email-service/docs" // docs is generated by Swag CLI, you have to import it.
)

const address = "http://localhost:8085"

func main() {
	log.Printf("Serving API Swagger docs at: %s/", address)
	http.ListenAndServe(":8085", httpSwagger.Handler(
		httpSwagger.URL(fmt.Sprintf("%s/swagger/doc.json", address)), //The url pointing to API definition"
	))
}
