package executor

import (
	"fmt"
	"github.com/jakecoffman/cron"
	"github.com/jakecoffman/gorunner/models"
	"github.com/nu7hatch/gouuid"
	"time"
)

var c *cron.Cron
var triggers map[string]struct{}

func init() {
	c = cron.New()
	c.Start()
	c.AddFunc("0 * * * *", func() { fmt.Println("test ran at " + time.Now().Format("2006-01-02 15:04:05")) }, "test")
}

func AddTrigger(t models.Trigger) {
	c.AddFunc(t.Schedule, func() { findAndRun(t) }, t.Name)
}

func RemoveTrigger(name string) {
	c.RemoveJob(name)
	println("Trigger has been removed")
}

// Walks through each job, seeing if the trigger who's turn it is to execute is attached. Executes those jobs.
func findAndRun(t models.Trigger) {
	jobList := models.GetJobList()
	for _, job := range jobList.GetAll() {
		j := job.(models.Job)
		for _, trigger := range j.Triggers {
			if trigger == t.ID() {
				fmt.Println("Running job " + j.Name)
				runnit(j)
				break
			}
		}
	}
}

// Gathers the tasks attached to the given job and executes them.
func runnit(j models.Job) {
	tasksList := models.GetTaskList()
	runsList := models.GetRunList()
	id, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}
	var tasks []models.Task
	for _, taskName := range j.Tasks {
		task, err := tasksList.Get(taskName)
		if err != nil {
			panic(err)
		}
		t := task.(models.Task)
		tasks = append(tasks, t)
	}
	err = runsList.AddRun(id.String(), j, tasks)
	if err != nil {
		panic(err)
	}
}
