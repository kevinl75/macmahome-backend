package test

// Basic imports
import (
	"bytes"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/api"
	"github.com/kevinl75/macmahome-backend/model"
	"github.com/kevinl75/macmahome-backend/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type TaskRouteAPITests struct {
	suite.Suite
	TestDbConn       *gorm.DB
	TaskItems        []model.Task
	TestingDate      time.Time
	TestRouter       *gin.Engine
	FakeJSONTaskItem string
}

// before all the suite
func (suite *TaskRouteAPITests) SetupSuite() {

	os.Setenv("MACMAHOME_TEST", "true")

	dbConn := utils.NewDBConnection()
	err := dbConn.AutoMigrate(&model.Project{}, &model.Task{}, &model.Note{})
	if err != nil {
		panic("an error occured during the migration.")
	}

	suite.TestDbConn = dbConn
	suite.TestRouter = api.NewRouter()
	suite.TestingDate, _ = time.Parse("2006-01-02", "2022-01-01")
	suite.TaskItems = []model.Task{
		{TaskId: 1, TaskName: "Test 1", TaskIsComplete: false, TaskDuration: 30, TaskDate: suite.TestingDate},
		{TaskId: 2, TaskName: "Test 2", TaskIsComplete: true, TaskDuration: 90, TaskDate: suite.TestingDate},
		{TaskId: 3, TaskName: "Test 3", TaskIsComplete: false, TaskDuration: 120, TaskDate: suite.TestingDate},
	}
	suite.FakeJSONTaskItem = "{\"task_id\": 12,\"task_name\": 456}"

	dbConn.Create(&suite.TaskItems[0])
	dbConn.Create(&suite.TaskItems[1])
}

// after all the suite
func (suite *TaskRouteAPITests) TearDownSuite() {

	os.Setenv("MACMAHOME_TEST", "true")

	dbConn := utils.NewDBConnection()
	err := dbConn.Migrator().DropTable(&model.Project{}, &model.Task{}, &model.Note{})
	if err != nil {
		panic("an error occured during the migration.")
	}
}

func TestTaskRoutesAPI(t *testing.T) {
	suite.Run(t, new(TaskRouteAPITests))
}

func (suite *TaskRouteAPITests) TestPostRoute() {

	postRequestData, _ := json.Marshal(suite.TaskItems[2])
	w := ServeTestHTTPRequest("POST", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	assert.Equal(suite.T(), 201, w.Code)
	assert.Equal(suite.T(), task, suite.TaskItems[2])

	w = ServeTestHTTPRequest("POST", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)

	assert.Equal(suite.T(), 500, w.Code)
	suite.TestDbConn.Delete(&task)

	w = ServeTestHTTPRequest("POST", "/task", bytes.NewBufferString(suite.FakeJSONTaskItem), suite.TestRouter)
	assert.Equal(suite.T(), 400, w.Code)
}

func (suite *TaskRouteAPITests) TestGetRoute() {

	w := ServeTestHTTPRequest("GET", "/task/1", nil, suite.TestRouter)
	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), task, suite.TaskItems[0])

	w = ServeTestHTTPRequest("GET", "/task/4", nil, suite.TestRouter)

	assert.Equal(suite.T(), 404, w.Code)

	var tasks []model.Task

	w = ServeTestHTTPRequest("GET", "/task", nil, suite.TestRouter)

	json.Unmarshal(w.Body.Bytes(), &tasks)

	assert.Equal(suite.T(), 2, len(tasks))
}

func (suite *TaskRouteAPITests) TestDeleteRoute() {

	w := ServeTestHTTPRequest("DELETE", "/task/2", nil, suite.TestRouter)

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), task, suite.TaskItems[1])
	suite.TestDbConn.Create(&task)

	w = ServeTestHTTPRequest("DELETE", "/task/4", nil, suite.TestRouter)

	assert.Equal(suite.T(), 404, w.Code)
}

func (suite *TaskRouteAPITests) TestPatchRoute() {

	newTaskName := "Updated task name"
	postRequestData, _ := json.Marshal(model.Task{
		TaskId:         2,
		TaskName:       newTaskName,
		TaskIsComplete: true,
		TaskDuration:   90,
		TaskDate:       suite.TestingDate,
	})
	w := ServeTestHTTPRequest("PATCH", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)

	var task model.Task
	json.Unmarshal(w.Body.Bytes(), &task)

	assert.Equal(suite.T(), 200, w.Code)
	assert.Equal(suite.T(), task.TaskName, newTaskName)

	postRequestData, _ = json.Marshal(suite.TaskItems[2])
	w = ServeTestHTTPRequest("PATCH", "/task", bytes.NewBuffer(postRequestData), suite.TestRouter)

	assert.Equal(suite.T(), 500, w.Code)
	suite.TestDbConn.Save(&task)
}
