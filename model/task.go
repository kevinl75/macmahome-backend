package model

import (
	"errors"
	"fmt"
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
	tx := dbConn.Begin()

	var result *gorm.DB
	if t.ProjectId == 0 {
		result = tx.Omit("ProjectId").Create(&t)
	} else {
		result = tx.Create(&t)
	}

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.First(&t)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (t *Task) UpdateTask() error {

	dbConn := utils.NewDBConnection()
	tx := dbConn.Begin()

	var result *gorm.DB
	if t.ProjectId == 0 {
		result = tx.Omit("ProjectId").Updates(&t)
	} else {
		result = tx.Updates(&t)
	}

	if result.Error != nil {
		tx.Rollback()
		fmt.Println(result.Error)
		return result.Error
	}

	result = tx.First(&t)

	if result.Error != nil {
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

func ReturnTask(id uint) (Task, error) {

	var task Task
	dbConn := utils.NewDBConnection()
	result := dbConn.First(&task, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Task{}, nil
		}
		return Task{}, result.Error
	}

	return task, nil
}

func ReturnTasks() ([]Task, error) {

	var tasks []Task
	dbConn := utils.NewDBConnection()
	result := dbConn.Find(&tasks)

	if result.Error != nil {
		return []Task{}, result.Error
	}

	return tasks, nil
}
