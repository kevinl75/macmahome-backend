package test

import (
	"bytes"
	"encoding/json"

	"github.com/kevinl75/macmahome-backend/model"
	"github.com/stretchr/testify/assert"
)

// Basic imports

func (suite *RouteAPITests) TestTaskPostRoute() {

	// We assert a normal insertion with no error will return code 201 and the newly created object.
	postRequestData, _ := json.Marshal(suite.TaskItems[2])
	w := ServeTestHTTPRequest("POST", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)
	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), task, suite.TaskItems[2])

	// We assert the insertion of an entity already inserted will return an error code 500.
	w = ServeTestHTTPRequest("POST", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)
	assert.Equal(suite.T(), 500, w.Code)

	suite.TestDbConn.Delete(&task)

	// We assert the insertion of an entity poorly formatted will return an error 400.
	w = ServeTestHTTPRequest("POST", "/task", bytes.NewBufferString(suite.FakeJSONTaskItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert the insertion of a task linked to an unknown project entity will fail and return a code 500

	// This test has been desactivated because the SQLite database used to run the test does not
	// support the foreign key constraint. This could be configured.

	// postRequestData, _ = json.Marshal(suite.TaskItems[3])
	// w = ServeTestHTTPRequest("POST", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)
	// assert.Equal(suite.T(), 500, w.Code)
}

func (suite *RouteAPITests) TestTaskGetRoute() {

	// We assert fetching an existing entity will return it correctly with the code 200
	w := ServeTestHTTPRequest("GET", "/task/1", nil, suite.TestRouter)
	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), task.TaskId, suite.TaskItems[0].TaskId)
	assert.Equal(suite.T(), task, suite.TaskItems[0])

	// We assert fetching an entity that is not in our database will return a code 400
	w = ServeTestHTTPRequest("GET", "/task/4", nil, suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert fetching an existing with a none positive integer id will return a code 404
	w = ServeTestHTTPRequest("GET", "/task/test", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
	w = ServeTestHTTPRequest("GET", "/task/-2", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert fetching the route without ids will return all the entities inserted during the setup.
	var tasks []model.Task
	w = ServeTestHTTPRequest("GET", "/task", nil, suite.TestRouter)
	json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 2, len(tasks))

	// We assert that if we fetch the tasks linked to the project with id 1, it will return the correct tasks and the code 200
	w = ServeTestHTTPRequest("GET", "/project/1/task", nil, suite.TestRouter)
	json.Unmarshal(w.Body.Bytes(), &tasks)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), 1, len(tasks))
	assert.Equal(suite.T(), suite.TaskItems[0], tasks[0])
}

func (suite *RouteAPITests) TestTaskDeleteRoute() {

	// We assert that trying to delete an existing object will effectively delete it and return the entity
	w := ServeTestHTTPRequest("DELETE", "/task/2", nil, suite.TestRouter)
	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), task, suite.TaskItems[1])
	suite.TestDbConn.Create(&task)

	// We assert that trying to delete a not existing entity will return a code 404
	w = ServeTestHTTPRequest("DELETE", "/task/4", nil, suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert that trying to delete an entity with an id that is not a positive integer will return a code 400
	w = ServeTestHTTPRequest("DELETE", "/task/test", nil, suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *RouteAPITests) TestTaskPatchRoute() {

	// We assert that patching an entity with a well formatted object will effectively update it and return a code 200
	updatedTask := suite.TaskItems[1]
	updatedTask.TaskName = "Updated name."
	patchRequestData, _ := json.Marshal(updatedTask)
	w := ServeTestHTTPRequest("PATCH", "/task", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)
	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), task, updatedTask)
	assert.NotEqual(suite.T(), task.TaskName, suite.TaskItems[1].TaskName)
	suite.TestDbConn.Save(&suite.TaskItems[1])

	// We assert that patching a not existing entity will return a code 404
	patchRequestData, _ = json.Marshal(suite.TaskItems[2])
	w = ServeTestHTTPRequest("PATCH", "/task", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	assert.Equal(suite.T(), 404, w.Code)

	// We assert that patching a entity with a poorly formatted object will return a code 400
	w = ServeTestHTTPRequest("PATCH", "/task", bytes.NewBufferString(suite.FakeJSONTaskItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)

	// We assert that patching a entity to link it to a unknown project will return a code 500

	// This test has been desactivated because the SQLite database used to run the test does not
	// support the foreign key constraint. This could be configured.

	// taskToUpdate := suite.TaskItems[0]
	// taskToUpdate.ProjectId = 50
	// patchRequestData, _ = json.Marshal(taskToUpdate)
	// w = ServeTestHTTPRequest("PATCH", "/task", bytes.NewBuffer(patchRequestData), suite.TestRouter)
	// assert.Equal(suite.T(), 500, w.Code)
}
