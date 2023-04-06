package storage

type Storage interface {
	Get(id string) (string, bool)
	Set(id, link string)
	Delete(id string)
	List() map[string]string
	Count() int
}

type LinkStorage struct {
	links map[string]string
}

func NewLinkStorage() *LinkStorage {
	return &LinkStorage{
		links: make(map[string]string),
	}
}

func (s *LinkStorage) Get(id string) (string, bool) {
	link, ok := s.links[id]
	return link, ok
}

func (s *LinkStorage) Set(id, link string) {
	s.links[id] = link
}

func (s *LinkStorage) Delete(id string) {
	delete(s.links, id)
}

func (s *LinkStorage) List() map[string]string {
	return s.links
}

func (s *LinkStorage) Count() int {
	return len(s.links)
}
