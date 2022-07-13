package test

// Basic imports
import (
	"bytes"
	"encoding/json"

	"github.com/kevinl75/macmahome-backend/model"
	"github.com/stretchr/testify/assert"
)

func (suite *RouteAPITests) TestNotePostRoute() {

	// We assert a normal insertion with no error will return code 201 and the newly created object.
	postRequestData, _ := json.Marshal(suite.NoteItems[2])
	w := ServeTestHTTPRequest("POST", "/note", bytes.NewBuffer(postRequestData), suite.TestRouter)
	var note model.Note
	json.Unmarshal(w.Body.Bytes(), &note)
	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), note, suite.NoteItems[2])

	// We assert the insertion of an entity already inserted will return an error code 500.
	w = ServeTestHTTPRequest("POST", "/note", bytes.NewBuffer(postRequestData), suite.TestRouter)
	assert.Equal(suite.T(), 500, w.Code)

	suite.TestDbConn.Delete(&note)

	// We assert the insertion of an entity poorly formatted will return an error 400.
	w = ServeTestHTTPRequest("POST", "/note", bytes.NewBufferString(suite.FakeJSONNoteItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert the insertion of a note linked to an unknown project entity will fail and return a code 500

	// This test has been desactivated because the SQLite database used to run the test does not
	// support the foreign key constraint. This could be configured.

	// postRequestData, _ = json.Marshal(suite.NoteItems[3])
	// w = ServeTestHTTPRequest("POST", "/note", bytes.NewBuffer(postRequestData), suite.TestRouter)
	// assert.Equal(suite.T(), 500, w.Code)
}

func (suite *RouteAPITests) TestNoteGetRoute() {

	// We assert fetching an existing entity will return it correctly with the code 200
	w := ServeTestHTTPRequest("GET", "/note/1", nil, suite.TestRouter)
	var note model.Note
	json.Unmarshal(w.Body.Bytes(), &note)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), note.NoteId, suite.NoteItems[0].NoteId)
	assert.Equal(suite.T(), note, suite.NoteItems[0])

	// We assert fetching an entity that is not in our database will return a code 400
	w = ServeTestHTTPRequest("GET", "/note/4", nil, suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert fetching an existing with a none positive integer id will return a code 404
	w = ServeTestHTTPRequest("GET", "/note/test", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
	w = ServeTestHTTPRequest("GET", "/note/-2", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert fetching the route without ids will return all the entities inserted during the setup.
	var notes []model.Note
	w = ServeTestHTTPRequest("GET", "/note", nil, suite.TestRouter)
	json.Unmarshal(w.Body.Bytes(), &notes)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 2, len(notes))

	// We assert that if we fetch the notes linked to the project with 1, it will return the correct note and the code 200
	w = ServeTestHTTPRequest("GET", "/project/1/note", nil, suite.TestRouter)
	json.Unmarshal(w.Body.Bytes(), &notes)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 1, len(notes))
	assert.Equal(suite.T(), suite.NoteItems[0], notes[0])
}

func (suite *RouteAPITests) TestNoteDeleteRoute() {

	// We assert that trying to delete an existing object will effectively delete it and return the entity
	w := ServeTestHTTPRequest("DELETE", "/note/2", nil, suite.TestRouter)
	var note model.Note
	json.Unmarshal(w.Body.Bytes(), &note)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), note, suite.NoteItems[1])
	suite.TestDbConn.Create(&note)

	// We assert that trying to delete a not existing entity will return a code 404
	w = ServeTestHTTPRequest("DELETE", "/note/4", nil, suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert that trying to delete an entity with an id that is not a positive integer will return a code 400
	w = ServeTestHTTPRequest("DELETE", "/note/test", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *RouteAPITests) TestNotePatchRoute() {

	// We assert that patching an entity with a well formatted object will effectively update it and return a code 200
	updatedNote := suite.NoteItems[1]
	updatedNote.NoteName = "Updated name."
	patchRequestData, _ := json.Marshal(updatedNote)
	w := ServeTestHTTPRequest("PATCH", "/note", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	var note model.Note
	json.Unmarshal(w.Body.Bytes(), &note)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), note, updatedNote)
	assert.NotEqual(suite.T(), note.NoteName, suite.NoteItems[1].NoteName)
	suite.TestDbConn.Save(&suite.NoteItems[1])

	// We assert that patching a not existing entity will return a code 404
	patchRequestData, _ = json.Marshal(suite.NoteItems[2])
	w = ServeTestHTTPRequest("PATCH", "/note", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert that patching a entity with a poorly formatted object will return a code 400
	w = ServeTestHTTPRequest("PATCH", "/note", bytes.NewBufferString(suite.FakeJSONNoteItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert that patching a entity to link it to a unknown project will return a code 500

	// This test has been desactivated because the SQLite database used to run the test does not
	// support the foreign key constraint. This could be configured.

	// noteToUpdate := suite.NoteItems[0]
	// noteToUpdate.ProjectId = 50
	// patchRequestData, _ = json.Marshal(noteToUpdate)
	// w = ServeTestHTTPRequest("PATCH", "/note", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	// assert.Equal(suite.T(), 500, w.Code)
}
