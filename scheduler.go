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

import "time"

const (
	DEFAULT_NUM_EXECUTORS = 5
)

type Scheduler struct {
	repository JobRepository

	jobChan       chan Job
	jobStatusChan chan Job
}

func NewScheduler() *Scheduler {
	return NewCustomScheduler(DEFAULT_NUM_EXECUTORS, NewInMemoryJobRepository(), DefaultExecutor)
}

func NewCustomScheduler(numExecutors int, repository JobRepository, executorStrategy Executor) *Scheduler {
	if numExecutors <= 0 {
		numExecutors = DEFAULT_NUM_EXECUTORS
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

func (s *Scheduler) ScheduleNow(task *Task) *Job {
	if task == nil {
		return nil
	}
	job := NewJob(task, time.Now())
	s.repository.Save(job)
	// TODO handle error on save
	s.jobChan <- job
	return &job
}

func (s *Scheduler) Schedule(task *Task, scheduledOn time.Time) *Job {
	if task == nil || time.Now().After(scheduledOn) {
		return nil
	}
	job := NewJob(task, scheduledOn)
	s.repository.Save(job)
	// TODO handle error on save
	go func() {
		time.Sleep(scheduledOn.Sub(time.Now()))
		s.jobChan <- job
	}()
	return &job
}

func (s *Scheduler) JobStatus(jobID string) *Job {
	if job, err := s.repository.Load(jobID); err == nil {
		return &job
	}
	return nil
}

func (s *Scheduler) Destroy() {
	close(s.jobChan)
	close(s.jobStatusChan)
}
