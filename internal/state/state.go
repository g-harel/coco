package state

import (
	"encoding/json"
	"io/ioutil"
	"os"

	"github.com/g-harel/coco/internal/log"
)

type State struct {
	path     string
	contents map[string]map[string]string
}

func NewFromFile(path string) State {
	data := map[string]map[string]string{}
	if path != "" {
		file, err := ioutil.ReadFile(path)
		if err == nil {
			json.Unmarshal(file, &data)
		}
	}
	return State{path: path, contents: data}
}

func (s *State) Save() {
	if s.path == "" {
		return
	}

	f, err := os.Create(s.path)
	if err != nil {
		log.Error("create state file: %s\n", err)
		return
	}
	defer f.Close()

	data, _ := json.Marshal(s.contents)
	_, err = f.Write(data)
	if err != nil {
		log.Error("save state file: %s\n", err)
	}
}

func (s *State) Read(namespace, key string) (value string, ok bool) {
	value, ok = s.contents[namespace][key]
	return
}

func (s *State) Write(namespace, key, value string) {
	if _, ok := s.contents[namespace]; !ok {
		s.contents[namespace] = map[string]string{}
	}
	s.contents[namespace][key] = value
}
