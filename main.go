package main

// @title Macmahome application
// @description Backend API for the Macmahome application. It permit to manage all the
// Macmahome application entities.
// @version 0.1
// @contact.name Kevin L
// @contact.email <kevin.letupe@gmail.com>
// @license.name Apache License
// @license.url http://www.apache.org/licenses/

// @accept json
// @produce json

// @host localhost:8080
// @BasePath /
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
