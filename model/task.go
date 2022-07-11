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
	TaskDuration   uint      `json:"task_duration"`
	TaskDate       time.Time `json:"task_date"`
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

	if result.Error != nil {
		return result.Error
	}

	result = dbConn.First(&t)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (t *Task) UpdateTask() error {

	dbConn := utils.NewDBConnection()
	tx := dbConn.Begin()

	var result *gorm.DB

	if t.ProjectId == 0 {
		result = tx.Debug().Omit("ProjectId").Updates(&t)
	} else {
		result = tx.Debug().Updates(&t)
	}

	if result.Error != nil {
		log.Print("update failed")
		tx.Rollback()
		return result.Error
	}

	result = tx.First(&t)

	if result.Error != nil {
		log.Print("update failed")
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
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

func ReturnTask(id string) Task {

	var task Task
	dbConn := utils.NewDBConnection()
	result := dbConn.First(&task, id)

	if result.RowsAffected == 0 {
		return Task{}
	}

	return task
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
