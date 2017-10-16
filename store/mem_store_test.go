// Copyright 2017 Northern.tech AS
//
//    Licensed under the Apache License, Version 2.0 {the "License"};
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package store

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStoreOne(t *testing.T) {
	s := InitStore()
	assert.NotNil(t, s)

	s.Store("test", 1)
	assert.Len(t, s.data, 1)

	s.Store("test", 2)
	assert.Len(t, s.data, 1)

	all := s.GetAll("test")
	assert.Len(t, all, 2)

	last := s.GetLast("test")
	assert.Equal(t, 2, last.Value)

	empty := s.GetLast("non-existing")
	assert.Equal(t, Store{}, empty)

	allEmpty := s.GetAll("non-existing")
	assert.Nil(t, allEmpty)
}

func TestRotate(t *testing.T) {
	s := InitStore()
	assert.NotNil(t, s)

	for i := 0; i < MAX_SIZE+1; i++ {
		s.Store("test", i)
	}
	all := s.GetAll("test")
	assert.Len(t, all, MAX_SIZE)
	assert.Equal(t, 1, all[0].Value)

	last := s.GetLast("test")
	assert.Equal(t, MAX_SIZE, last.Value)
}

func TestEmpty(t *testing.T) {
	// do not init store
	s := MemStore{}
	empty := s.GetLast("non-existing")
	assert.Equal(t, Store{}, empty)

	allEmpty := s.GetAll("non-existing")
	assert.Nil(t, allEmpty)
}

func TestMultiple(t *testing.T) {
	s := InitStore()
	assert.NotNil(t, s)

	s.Store("test", 1)
	assert.Len(t, s.data, 1)

	s.Store("test-1", 10)
	assert.Len(t, s.data, 2)

	last := s.GetLast("test")
	assert.Equal(t, 1, last.Value)

	last = s.GetLast("test-1")
	assert.Equal(t, 10, last.Value)
}
