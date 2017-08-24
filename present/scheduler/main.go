package main

import (
	"errors"
	"fmt"
	"time"

	scheduler "github.com/dakiva/go-scheduler"
)

func main() {
	fmt.Println("Starting test", time.Now())
	s := scheduler.NewScheduler()
	defer s.Destroy()
	future := time.Now().Add(3 * time.Second)
	futureJobIDs := make([]string, 0)
	currentJobIDs := make([]string, 0)
	for i := 0; i < 5; i++ {
		cp := i
		testTask, _ := scheduler.NewTask("test task", func() {
			fmt.Println("Running a future task #", cp)
		}, nil)
		job, _ := s.Schedule(testTask, future)
		futureJobIDs = append(futureJobIDs, job.ID)
	}
	for i := 5; i < 15; i++ {
		var job scheduler.Job
		cp := i
		if i%4 != 0 {
			testTask, _ := scheduler.NewTask("test task", func() {
				fmt.Println("Running a task #", cp)
			}, nil)
			job, _ = s.ScheduleNow(testTask)
		} else {
			failedTestTask, _ := scheduler.NewTask("test failed task", func() error {
				fmt.Println("Running a failed task #", cp)
				return errors.New("Error")
			}, nil)
			job, _ = s.ScheduleNow(failedTestTask)
		}
		currentJobIDs = append(currentJobIDs, job.ID)
	}
	time.Sleep(1)
	printJobs(currentJobIDs, s)
	printJobs(futureJobIDs, s)
	time.Sleep(4 * time.Second)
	fmt.Println("Future Jobs after 4s")
	printJobs(futureJobIDs, s)
	fmt.Println("Ending test", time.Now())
}

func printJobs(jobIDs []string, s *scheduler.Scheduler) {
	for _, id := range jobIDs {
		j, _ := s.JobStatus(id)
		fmt.Printf("[%s] %s\n", j.ID, j.Status)
	}
}
