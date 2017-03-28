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
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultScheduler(t *testing.T) {
	// given
	s := NewScheduler()
	defer s.Destroy()
	testTask, err := NewTask("test task", func() {
		fmt.Println("Running a task")
	}, nil)
	assert.NoError(t, err)

	failedTestTask, err := NewTask("test failed task", func() error {
		fmt.Println("Running a failed task")
		return errors.New("Error")
	}, nil)
	assert.NoError(t, err)

	// when
	job := s.ScheduleNow(testTask)
	job2 := s.ScheduleNow(failedTestTask)

	time.Sleep(2 * time.Second)

	completedJob := s.JobStatus(job.ID)
	failedJob := s.JobStatus(job2.ID)

	// then
	assert.Equal(t, SCHEDULED_STATUS, job.Status)
	assert.Nil(t, job.StartedOn)
	assert.Equal(t, COMPLETED_STATUS, completedJob.Status)
	assert.NotNil(t, completedJob.StartedOn)
	assert.NotNil(t, completedJob.CompletedOn)
	assert.Nil(t, completedJob.RunErr)

	assert.Equal(t, SCHEDULED_STATUS, job2.Status)
	assert.Nil(t, job2.StartedOn)
	assert.Equal(t, FAILED_STATUS, failedJob.Status)
	assert.Error(t, failedJob.RunErr)
	assert.NotNil(t, failedJob.StartedOn)
	assert.NotNil(t, failedJob.CompletedOn)
}
