package storage

//storage, key - generated phrase, value - link
type LinkStorage struct {
	links map[string]string
}

//constructor
func NewLinkStorage() *LinkStorage {
	return &LinkStorage{
		links: make(map[string]string),
	}
}

//get link by key
func (s *LinkStorage) Get(id string) (string, bool) {
	link, ok := s.links[id]
	return link, ok
}

//set link by key
func (s *LinkStorage) Set(id, link string) {
	s.links[id] = link
}

//delete link by key
func (s *LinkStorage) Delete(id string) {
	delete(s.links, id)
}

//get all links
func (s *LinkStorage) List() map[string]string {
	return s.links
}

//get count of links
func (s *LinkStorage) Count() int {
	return len(s.links)
}
