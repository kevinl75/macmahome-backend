package model

type Note struct {
	NoteId      uint `gorm:"primaryKey"`
	NoteName    string
	NoteContent string
}
