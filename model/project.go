package model

import (
	"errors"

	"github.com/kevinl75/macmahome-backend/utils"
	"gorm.io/gorm"
)

type Project struct {
	ProjectId          uint   `gorm:"primaryKey" json:"project_id"`
	ProjectName        string `json:"project_name"`
	ProjectDescription string `json:"project_description"`
	Tasks              []Task `json:"project_tasks"`
}

func (p *Project) CreateProject() error {

	dbConn := utils.NewDBConnection()
	tx := dbConn.Begin()

	var result *gorm.DB
	result = tx.Create(&p)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.First(&p)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (p *Project) UpdateProject() error {

	dbConn := utils.NewDBConnection()
	tx := dbConn.Begin()

	var result *gorm.DB
	result = tx.Updates(&p)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	result = tx.First(&p)

	if result.Error != nil {
		tx.Rollback()
		return result.Error
	}

	tx.Commit()
	return nil
}

func (p Project) DeleteProject() error {
	dbConn := utils.NewDBConnection()

	result := dbConn.Unscoped().Delete(&p)

	if result.RowsAffected == 0 {
		return result.Error
	}

	return nil
}

func (p Project) IsEqual(p2 Project) bool {
	if p.ProjectId != p2.ProjectId {
		return false
	}
	if p.ProjectName != p2.ProjectName {
		return false
	}
	if p.ProjectDescription != p2.ProjectDescription {
		return false
	}
	for id, task := range p.Tasks {
		if task != p2.Tasks[id] {
			return false
		}
	}
	return true
}

func ReturnProject(id uint) (Project, error) {

	var project Project
	dbConn := utils.NewDBConnection()
	result := dbConn.Preload("Tasks").First(&project, id)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return Project{}, nil
		}
		return Project{}, result.Error
	}

	return project, nil
}

func ReturnProjects() ([]Project, error) {

	var projects []Project
	dbConn := utils.NewDBConnection()
	result := dbConn.Preload("Tasks").Find(&projects)

	if result.Error != nil {
		return []Project{}, result.Error
	}

	return projects, nil
}
