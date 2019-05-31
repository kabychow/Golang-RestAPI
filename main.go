package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"os"
	"two-server/app"
	"two-server/controllers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/user/register", controllers.CreateAccount).Methods("POST")
	router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	router.Use(app.JwtAuthentication)

	err := http.ListenAndServe(":" + os.Getenv("app_port"), router)
	if err != nil {
		fmt.Println(err)
	}
}