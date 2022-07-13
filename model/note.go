package model

import (
	"errors"
	"fmt"
	"time"

	"github.com/kevinl75/macmahome-backend/utils"
	"gorm.io/gorm"
)

type Note struct {
	NoteId      uint      `gorm:"primaryKey" json:"note_id"`
	NoteName    string    `json:"note_name"`
	NoteContent string    `json:"note_content"`
	NoteDate    time.Time `json:"note_date"`
	ProjectId   uint      `json:"project_id"`
}

func (n *Note) CreateNote() error {

	dbConn := utils.NewDBConnection()
	tx := dbConn.Begin()

	var result *gorm.DB
	fmt.Print(n)
	if n.ProjectId == 0 {
		result = tx.Omit("ProjectId").Create(&n)
	} else {
		result = tx.Create(&n)
	}

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.First(&n)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (n *Note) UpdateNote() error {

	dbConn := utils.NewDBConnection()
	tx := dbConn.Begin()

	var result *gorm.DB
	if n.ProjectId == 0 {
		result = tx.Omit("ProjectId").Updates(&n)
	} else {
		result = tx.Updates(&n)
	}

	if result.Error != nil {
		tx.Rollback()
		fmt.Println(result.Error)
		return result.Error
	}

	result = tx.First(&n)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (n Note) DeleteNote() error {
	dbConn := utils.NewDBConnection()

	result := dbConn.Delete(&n)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func ReturnNote(id uint) (Note, error) {

	var note Note
	dbConn := utils.NewDBConnection()
	result := dbConn.First(&note, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Note{}, nil
		}
		return Note{}, result.Error
	}

	return note, nil
}

func ReturnNotes() ([]Note, error) {

	var notes []Note
	dbConn := utils.NewDBConnection()
	result := dbConn.Find(&notes)

	if result.Error != nil {
		return []Note{}, result.Error
	}

	return notes, nil
}

func ReturnNotesByProjectId(projectId uint) ([]Note, error) {
	var notes []Note
	dbConn := utils.NewDBConnection()
	result := dbConn.Where("project_id = ?", projectId).Find(&notes)

	if result.Error != nil {
		return []Note{}, result.Error
	}

	return notes, nil
}
