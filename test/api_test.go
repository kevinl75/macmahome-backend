package test

import (
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinl75/macmahome-backend/api"
	"github.com/kevinl75/macmahome-backend/model"
	"github.com/kevinl75/macmahome-backend/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type RouteAPITests struct {
	suite.Suite
	TestDbConn          *gorm.DB
	NoteItems           []model.Note
	TaskItems           []model.Task
	ProjectItems        []model.Project
	TestingDate         time.Time
	TestRouter          *gin.Engine
	FakeJSONProjectItem string
	FakeJSONTaskItem    string
	FakeJSONNoteItem    string
}

// before all the suite
func (suite *RouteAPITests) SetupSuite() {

	os.Setenv("MACMAHOME_TEST", "true")

	dbConn := utils.NewDBConnection()
	err := dbConn.AutoMigrate(&model.Project{}, &model.Note{}, &model.Task{})
	if err != nil {
		panic("an error occured during the migration.")
	}
	suite.TestDbConn = dbConn
	suite.TestRouter = api.NewRouter()
	suite.TestingDate, _ = time.Parse("2006-01-02", "2022-01-01")
	suite.NoteItems = []model.Note{
		{NoteId: 1, NoteName: "Test 1", NoteContent: "Note content 1", NoteDate: suite.TestingDate, ProjectId: 1},
		{NoteId: 2, NoteName: "Test 2", NoteContent: "Note content 2", NoteDate: suite.TestingDate},
		{NoteId: 3, NoteName: "Test 3", NoteContent: "Note content 3", NoteDate: suite.TestingDate},
		{NoteId: 4, NoteName: "Test 4", NoteContent: "Note content 4", NoteDate: suite.TestingDate, ProjectId: 3},
	}
	suite.TaskItems = []model.Task{
		{TaskId: 1, TaskName: "Test 1", TaskIsComplete: false, TaskDuration: 30, TaskDate: suite.TestingDate, ProjectId: 1},
		{TaskId: 2, TaskName: "Test 2", TaskIsComplete: true, TaskDuration: 90, TaskDate: suite.TestingDate},
		{TaskId: 3, TaskName: "Test 3", TaskIsComplete: false, TaskDuration: 120, TaskDate: suite.TestingDate},
		{TaskId: 4, TaskName: "Test 4", TaskIsComplete: false, TaskDuration: 120, TaskDate: suite.TestingDate},
	}
	suite.ProjectItems = []model.Project{
		{ProjectId: 1, ProjectName: "Test 1", ProjectDescription: "I'm a project description.", Tasks: []model.Task{}, Notes: []model.Note{}},
		{ProjectId: 2, ProjectName: "Test 2", ProjectDescription: "I'm a project description.", Tasks: []model.Task{}, Notes: []model.Note{}},
		{ProjectId: 3, ProjectName: "Test 3", ProjectDescription: "I'm a project description."},
	}
	suite.FakeJSONNoteItem = "{\"note_id\": 12,\"note_name\": 456}"
	suite.FakeJSONTaskItem = "{\"task_id\": 12,\"task_name\": 456}"
	suite.FakeJSONProjectItem = "{\"project_id\": 12,\"project_name\": 456}"

	dbConn.Omit(clause.Associations).Create(&suite.ProjectItems[0])
	dbConn.Omit(clause.Associations).Create(&suite.ProjectItems[1])
	dbConn.Create(&suite.TaskItems[0])
	dbConn.Omit("ProjectId").Create(&suite.TaskItems[1])
	dbConn.Create(&suite.NoteItems[0])
	dbConn.Omit("ProjectId").Create(&suite.NoteItems[1])
}

// after all the suite
func (suite *RouteAPITests) TearDownSuite() {

	os.Setenv("MACMAHOME_TEST", "true")

	dbConn := utils.NewDBConnection()
	err := dbConn.Migrator().DropTable(&model.Project{}, &model.Task{}, &model.Note{})
	if err != nil {
		panic("an error occured during the migration.")
	}
}

func TestRoutesAPI(t *testing.T) {
	suite.Run(t, new(RouteAPITests))
}
