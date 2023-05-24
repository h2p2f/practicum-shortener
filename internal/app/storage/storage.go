package storage

import (
	"encoding/json"
	"fmt"
)

type jsonLinks struct {
	UUID   int    `json:"uuid"`
	Short  string `json:"short_url"`
	Origin string `json:"original_url"`
}

// storage, key - generated phrase, value - link
type LinkStorage struct {
	links map[string]string
}

// constructor
func NewLinkStorage() *LinkStorage {
	return &LinkStorage{
		links: make(map[string]string),
	}
}

// get link by key
func (s *LinkStorage) Get(id string) (string, bool) {
	link, ok := s.links[id]
	return link, ok
}

// set link by key
func (s *LinkStorage) Set(id, link string) {
	s.links[id] = link
}

// delete link by key
func (s *LinkStorage) Delete(id string) {
	delete(s.links, id)
}

// get all links
func (s *LinkStorage) List() map[string]string {
	return s.links
}

// get count of links
func (s *LinkStorage) Count() int {
	return len(s.links)
}

func (s *LinkStorage) GetAllSliced() [][]byte {
	var result [][]byte
	id := 1
	for k, v := range s.links {
		jsonData := jsonLinks{
			UUID:   id,
			Short:  k,
			Origin: v,
		}
		out, err := json.Marshal(jsonData)
		if err != nil {
			fmt.Println(err)
		}
		result = append(result, out)
		id++
	}
	return result
}

func (s *LinkStorage) LoadAll(data [][]byte) {
	var jsonData jsonLinks
	for _, v := range data {
		err := json.Unmarshal(v, &jsonData)
		if err != nil {
			fmt.Println(err)
		}
		s.links[jsonData.Short] = jsonData.Origin
	}
}
