package model

import (
	"errors"
	"log"
	"time"

	"github.com/kevinl75/macmahome-backend/utils"
	"gorm.io/gorm"
)

// Task entity definition
type Task struct {
	TaskId         uint      `gorm:"primaryKey" json:"task_id"`
	TaskName       string    `json:"task_name"`
	TaskIsComplete bool      `json:"task_is_complete"`
	TaskDuration   int       `json:"task_duration"`
	TaskDate       time.Time `json:"task_date" time_format:"2006-01-02"`
	ProjectId      uint      `json:"project_id"`
}

func (t *Task) CreateTask() error {

	dbConn := utils.NewDBConnection()

	var result *gorm.DB
	if t.ProjectId == 0 {
		result = dbConn.Omit("ProjectId").Create(&t)
	} else {
		result = dbConn.Create(&t)
	}

	if result.RowsAffected == 0 {
		log.Print("Insertion Failed")
		return errors.New("insertion failed")
	}

	result = dbConn.First(&t)

	if result.RowsAffected == 0 {
		log.Print("select failed")
		return errors.New("select failed")
	}

	return nil
}

func (t *Task) UpdateTask() error {

	dbConn := utils.NewDBConnection()

	var result *gorm.DB

	if t.ProjectId == 0 {
		result = dbConn.Omit("ProjectId").Save(&t)
	} else {
		result = dbConn.Save(&t)
	}

	if result.RowsAffected == 0 {
		log.Print("update failed")
		return errors.New("select failed")
	}

	result = dbConn.First(&t)

	if result.RowsAffected == 0 {
		log.Print("select failed")
		return errors.New("select failed")
	}

	return nil
}

func (t Task) DeleteTask() error {
	dbConn := utils.NewDBConnection()

	result := dbConn.Unscoped().Delete(&t)

	if result.RowsAffected == 0 {
		log.Print("delete failed")
		return result.Error
	}

	return nil
}

func ReturnTask(id string) (Task, error) {

	var task Task
	dbConn := utils.NewDBConnection()
	result := dbConn.First(&task, id)

	if result.RowsAffected == 0 {
		log.Print("select failed")
		return Task{}, errors.New("select failed")
	}

	return task, nil
}

func ReturnTasks() ([]Task, error) {

	var tasks []Task
	dbConn := utils.NewDBConnection()
	result := dbConn.Find(&tasks)

	if result.RowsAffected == 0 {
		log.Print("fetch all failed.")
		return []Task{}, errors.New("fetch all failed")
	}

	return tasks, nil
}
