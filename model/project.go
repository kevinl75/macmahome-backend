package model

type Project struct {
	ProjectId          uint `gorm:"primaryKey"`
	ProjectName        string
	ProjectDescription string
	Tasks              []Task //`gorm:"foreignKey:ProjectId"`
}
