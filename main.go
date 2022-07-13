package main

import (
	"github.com/kevinl75/macmahome-backend/api"
	"github.com/kevinl75/macmahome-backend/model"
	"github.com/kevinl75/macmahome-backend/utils"
)

func MigrateModels() {

	dbConn := utils.NewDBConnection()
	err := dbConn.AutoMigrate(&model.Project{}, &model.Task{}, &model.Note{})
	if err != nil {
		panic("an error occured during the migration.")
	}
}

func main() {

	MigrateModels()

	router := api.NewRouter()

	router.Run()
}
