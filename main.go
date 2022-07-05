package main

import (
	"github.com/kevinl75/macmahome-backend/api"
)

func main() {

	router := api.NewRouter()

	router.Run()
}
