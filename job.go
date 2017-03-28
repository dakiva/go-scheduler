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
	"time"

	"github.com/pborman/uuid"
)

const (
	SCHEDULED_STATUS = "scheduled"
	RUNNING_STATUS   = "running"
	COMPLETED_STATUS = "completed"
	FAILED_STATUS    = "failed"
)

type Job struct {
	ID          string
	Status      string
	ScheduledOn time.Time
	StartedOn   *time.Time
	CompletedOn *time.Time
	Task        *Task
	RunErr      error
}

func NewJob(task *Task, scheduledOn time.Time) Job {
	return Job{
		ID:          uuid.New(),
		Task:        task,
		Status:      SCHEDULED_STATUS,
		ScheduledOn: scheduledOn,
	}
}

func (j Job) Start() Job {
	j.Status = RUNNING_STATUS
	now := time.Now()
	j.StartedOn = &now
	return j
}

func (j Job) Complete() Job {
	j.Status = COMPLETED_STATUS
	now := time.Now()
	j.CompletedOn = &now
	return j
}

func (j Job) Fail(err error) Job {
	j.Status = FAILED_STATUS
	j.RunErr = err
	now := time.Now()
	j.CompletedOn = &now
	return j
}

func (j Job) Valid() bool {
	return j.ID != "" && j.Status != "" && j.Task != nil
}
