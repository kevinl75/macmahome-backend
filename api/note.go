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

// createNote permit to handle a POST request to create a new Note.
// @summary      POST /note
// @description  Permit to create a new note entity from a JSON object. Note could be linked
//               to a project or not. Return the created entity
// @tags         note
// @accept       json
// @produce      json
// @param        note  body  model.Note  true  "New note to create"
// @Success      201  {object}  model.Note
// @Failure      400  {error}  error  "Error with the format of the provided data."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /note [post]
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

// returnNote permit to return a note based on its id
// @summary      GET /note/:id
// @description  Permit to return a Note entity based on its id.
// @tags         note
// @accept       mpfd
// @produce      json
// @param        id  path  uint  true  "Id of the note to return"
// @Success      200  {object}  model.Note
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      404  {error}  error  "Entity to return not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /note/:id [get]
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
		errorMsg := fmt.Errorf("no note entity with id %d", id).Error()
		c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
		return
	}

	c.JSON(http.StatusOK, note)
}

// returnNotes permit to return all the Note entities
// @summary      GET /note
// @description  Permit to return all the Note entities.
// @tags         note
// @accept       mpfd
// @produce      json
// @Success      200  {object}  []model.Note
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /note [get]
func returnNotes(c *gin.Context) {

	notes, err := model.ReturnNotes()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, notes)
}

// returnNotesByProjectId permit to return all the Note entities of a project with the project id
// @summary      GET /project/:id/note
// @description  Permit to return all the Note entities of a project with the project id.
// @tags         note
// @accept       mpfd
// @produce      json
// @Success      200  {object}  []model.Note
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /project/:id/note [get]
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

// updateNote permit to handle a PATCH request to update a Note.
// @summary      PATCH /note
// @description  Permit to update a note entity with a JSON object. Return the updated entity
// @tags         note
// @accept       json
// @produce      json
// @param        note  body  model.Note  true  "Updated note entity"
// @Success      200  {object}  model.Note
// @Failure      400  {error}  error  "Error with the format of the provided data."
// @Failure      404  {error}  error  "Entity to delete not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /note [patch]
func updateNote(c *gin.Context) {
	var updatedNote model.Note

	if err := c.ShouldBindJSON(&updatedNote); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := updatedNote.UpdateNote()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorMsg := fmt.Errorf("no note entity with id %d", updatedNote.NoteId).Error()
			c.JSON(http.StatusNotFound, gin.H{"error": errorMsg})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedNote)
}

// deleteNote permit to delete a note based on its id
// @summary      DELETE /note/:id
// @description  Permit to delete a note based on its id. Return the deleted entity.
// @tags         note
// @accept       mpfd
// @produce      json
// @param        id  path  uint  true  "Id of the note to delete"
// @Success      200  {object}  model.Note
// @Failure      400  {error}  error  "Error with the format of the provided id."
// @Failure      404  {error}  error  "Entity to delete not found error."
// @Failure      500  {error}  error  "Internal error during the the request handling."
// @Router       /note/:id [delete]
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
