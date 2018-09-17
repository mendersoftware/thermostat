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
	"time"
)

type Store struct {
	Value interface{} `json:"value"`
	Date  time.Time   `json:"time"`
}

type MemStore struct {
	data map[string]([]Store)
}

func InitStore() MemStore {
	data := make(map[string]([]Store), 0)
	return MemStore{data}
}

const MAX_SIZE = 10

func (s *MemStore) Store(name string, value interface{}) {
	data := s.data
	if _, ok := data[name]; ok {
		if len(data[name]) == MAX_SIZE {
			// rotate
			data[name] = data[name][1:]
		}
		data[name] = append(data[name], Store{value, time.Now()})
		return
	}

	data[name] = make([]Store, 0, 10)
	data[name] = append(data[name], Store{value, time.Now()})
}

func (s *MemStore) GetAll(name string) []Store {
	if s.data == nil {
		return nil
	}
	return s.data[name]
}

func (s *MemStore) GetLast(name string) Store {
	if s.data == nil {
		return Store{}
	}
	if bank, ok := s.data[name]; ok {
		return bank[len(bank)-1]
	}
	return Store{}
}
