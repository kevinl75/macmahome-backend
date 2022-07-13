package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/model"
	"gorm.io/gorm"
)

func createNote(c *gin.Context) {
	var newNote model.Note

	if err := c.ShouldBindJSON(&newNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := newNote.CreateNote()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newNote)
}

func returnNote(c *gin.Context) {

	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note, err := model.ReturnNote(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if (note == model.Note{}) {
		errorMsg := fmt.Errorf("no task entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, note)
}

func returnNotes(c *gin.Context) {

	notes, err := model.ReturnNotes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func returnNotesByProjectId(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	notes, err := model.ReturnNotesByProjectId(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

func updateNote(c *gin.Context) {
	var updatedNote model.Note

	if err := c.ShouldBindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := updatedNote.UpdateNote()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := fmt.Errorf("no project entity with id %d", updatedNote.NoteId).Error()
			c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedNote)
}

func deleteNote(c *gin.Context) {
	rawId := c.Param("id")
	id, err := strconv.ParseUint(rawId, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	noteToDelete, err := model.ReturnNote(uint(id))

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if (noteToDelete == model.Note{}) {
		errorMsg := fmt.Errorf("no task entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	noteToDelete.DeleteNote()
	c.JSON(http.StatusOK, noteToDelete)
}
