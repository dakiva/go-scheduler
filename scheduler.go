// Copyright 2017 Daniel Akiva

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package scheduler

import (
	"errors"
	"time"
)

const (
	DefaultNumExecutors = 5
)

var JobNotFound error = errors.New("Job not found")

type Scheduler struct {
	repository JobRepository

	jobChan       chan Job
	jobStatusChan chan Job
}

func NewScheduler() *Scheduler {
	return NewCustomScheduler(DefaultNumExecutors, NewInMemoryJobRepository(), DefaultExecutor)
}

func NewCustomScheduler(numExecutors int, repository JobRepository, executorStrategy Executor) *Scheduler {
	if numExecutors <= 0 {
		numExecutors = DefaultNumExecutors
	}
	scheduler := &Scheduler{
		repository: repository,
	}
	scheduler.jobChan, scheduler.jobStatusChan = make(chan Job, numExecutors), make(chan Job, numExecutors)
	go func() {
		for job := range scheduler.jobStatusChan {
			scheduler.repository.Save(job)
		}
	}()

	for i := 0; i < numExecutors; i++ {
		go func() {
			executorStrategy(scheduler.jobChan, scheduler.jobStatusChan)
		}()
	}
	return scheduler
}

func (s *Scheduler) ScheduleNow(task *Task) (Job, error) {
	return s.Schedule(task, time.Now())
}

func (s *Scheduler) Schedule(task *Task, scheduledOn time.Time) (Job, error) {
	if task == nil {
		return Job{}, errors.New("Task was empty")
	}
	job := NewJob(task, scheduledOn)
	if err := s.repository.Save(job); err != nil {
		return Job{}, err
	}
	go func() {
		time.Sleep(scheduledOn.Sub(time.Now()))
		s.jobChan <- job
	}()
	return job, nil
}

func (s *Scheduler) JobStatus(jobID string) (Job, error) {
	if job, err := s.repository.Load(jobID); err == nil {
		return job, nil
	}
	return Job{}, JobNotFound
}

func (s *Scheduler) Destroy() {
	close(s.jobChan)
	close(s.jobStatusChan)
}
