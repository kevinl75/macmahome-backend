package test

// Basic imports
import (
	"bytes"
	"encoding/json"

	"github.com/kevinl75/macmahome-backend/model"
	"github.com/stretchr/testify/assert"
)

func (suite *RouteAPITests) TestProjectPostRoute() {

	// We assert a normal insertion with no error will return code 201 and the newly created object.
	postRequestData, _ := json.Marshal(suite.ProjectItems[2])
	w := ServeTestHTTPRequest("POST", "/project", bytes.NewBuffer(postRequestData), suite.TestRouter)
	var project model.Project
	json.Unmarshal(w.Body.Bytes(), &project)
	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), project.ProjectName, suite.ProjectItems[2].ProjectName)

	// We assert that the insertion of an entity already inserted will return an error code 500.
	w = ServeTestHTTPRequest("POST", "/project", bytes.NewBuffer(postRequestData), suite.TestRouter)
	assert.Equal(suite.T(), 500, w.Code)

	suite.TestDbConn.Delete(&project)

	// We assert that the insertion of an entity badly formated will return an error 400.
	w = ServeTestHTTPRequest("POST", "/project", bytes.NewBufferString(suite.FakeJSONProjectItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *RouteAPITests) TestProjectGetRoute() {

	// We assert fetching an existing entity will return it correctly with the code 200
	w := ServeTestHTTPRequest("GET", "/project/2", nil, suite.TestRouter)
	var project model.Project
	json.Unmarshal(w.Body.Bytes(), &project)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), project.ProjectId, suite.ProjectItems[1].ProjectId)
	assert.Equal(suite.T(), true, project.IsEqual(suite.ProjectItems[1]))

	// We assert fetching an existing project entity with link to note and task will return it correctly with the code 200
	w = ServeTestHTTPRequest("GET", "/project/1", nil, suite.TestRouter)
	project2 := suite.ProjectItems[0]
	project2.Tasks = []model.Task{suite.TaskItems[0]}
	project2.Notes = []model.Note{suite.NoteItems[0]}
	json.Unmarshal(w.Body.Bytes(), &project)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), project.ProjectId, project2.ProjectId)
	assert.Equal(suite.T(), true, project.IsEqual(project2))

	// We assert fetching an entity that is not in our database will return a code 400
	w = ServeTestHTTPRequest("GET", "/project/4", nil, suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert fetching an existing with a none positive integer id will return a code 404
	w = ServeTestHTTPRequest("GET", "/project/test", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
	w = ServeTestHTTPRequest("GET", "/project/-2", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert fetching the route without ids will return all the entities inserted during the setup.
	var projects []model.Project
	w = ServeTestHTTPRequest("GET", "/project", nil, suite.TestRouter)
	json.Unmarshal(w.Body.Bytes(), &projects)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 2, len(projects))
}

func (suite *RouteAPITests) TestProjectDeleteRoute() {

	// We assert that trying to delete an existing object will effectively delete it and return the entity
	w := ServeTestHTTPRequest("DELETE", "/project/2", nil, suite.TestRouter)
	var project model.Project
	json.Unmarshal(w.Body.Bytes(), &project)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), true, project.IsEqual(suite.ProjectItems[1]))
	// assert.Equal(suite.T(), suite.ProjectItems[1], project)
	suite.TestDbConn.Create(&project)

	// We assert that trying to delete a not existing entity will return a code 404
	w = ServeTestHTTPRequest("DELETE", "/project/4", nil, suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert that trying to delete an entity with an id that is not a positive integer will return a code 400
	w = ServeTestHTTPRequest("DELETE", "/project/test", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert that deleted a project link to existing tasks and notes will return a code 500
	w = ServeTestHTTPRequest("DELETE", "/project/1", nil, suite.TestRouter)
	assert.Equal(suite.T(), 200, w.Code)
	suite.TestDbConn.Create(&(suite.ProjectItems[0]))
}

func (suite *RouteAPITests) TestProjectPatchRoute() {

	// We assert that patching an entity with a well formatted object will effectively update it and return a code 200
	updatedProject := suite.ProjectItems[1]
	updatedProject.ProjectName = "Updated name."
	patchRequestData, _ := json.Marshal(updatedProject)
	w := ServeTestHTTPRequest("PATCH", "/project", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	var project model.Project
	json.Unmarshal(w.Body.Bytes(), &project)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), project, updatedProject)
	assert.NotEqual(suite.T(), project.ProjectName, suite.ProjectItems[1].ProjectName)
	suite.TestDbConn.Save(&suite.ProjectItems[1])

	// We assert that patching a not existing entity will return a code 404
	patchRequestData, _ = json.Marshal(suite.ProjectItems[2])
	w = ServeTestHTTPRequest("PATCH", "/project", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert that patching a entity with a poorly formatted object will return a code 400
	w = ServeTestHTTPRequest("PATCH", "/project", bytes.NewBufferString(suite.FakeJSONProjectItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
}
