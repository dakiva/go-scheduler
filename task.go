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
	"reflect"
)

type Task struct {
	name   string
	f      interface{}
	params []interface{}
}

func NewTask(name string, f interface{}, params []interface{}) (*Task, error) {
	if name == "" {
		return nil, errors.New("Task name cannot be empty")
	}
	t := reflect.TypeOf(f)
	if t.Kind() != reflect.Func {
		return nil, errors.New("Task function must be a function.")
	}
	// TODO check params length and types against function param length and types
	return &Task{
		name:   name,
		f:      f,
		params: params,
	}, nil
}

func (t *Task) Name() string {
	return t.name
}

func (t *Task) Execute() error {
	exec := reflect.ValueOf(t.f)
	in := make([]reflect.Value, len(t.params))
	for k, param := range t.params {
		in[k] = reflect.ValueOf(param)
	}
	result := exec.Call(in)
	return findError(result)
}

func findError(result []reflect.Value) error {
	for _, v := range result {
		if v.CanInterface() {
			err, ok := v.Interface().(error)
			if ok {
				return err
			}
		}
	}
	return nil
}
