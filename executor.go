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

type Executor func(jobChan <-chan Job, jobStatusChan chan<- Job)

func DefaultExecutor(jobChan <-chan Job, jobStatusChan chan<- Job) {
	for job := range jobChan {
		startedJob := job.Start()
		jobStatusChan <- startedJob
		if err := job.Task.Execute(); err != nil {
			jobStatusChan <- startedJob.Fail(err)
		} else {
			jobStatusChan <- startedJob.Complete()
		}
	}
}
