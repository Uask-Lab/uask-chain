package search

type Search interface {
	AddDoc(interface{}) error
	SearchDoc(query string) ([]interface{}, error)
	UpdateDoc(id string, i interface{}) error
	DeleteDoc(id string) error
}

type NonSearch struct{}

func (n *NonSearch) AddDoc(i interface{}) error {
	return nil
}

func (n *NonSearch) SearchDoc(query string) ([]interface{}, error) {
	return nil, nil
}

func (n *NonSearch) UpdateDoc(id string, i interface{}) error {
	return nil
}

func (n *NonSearch) DeleteDoc(id string) error {
	return nil
}
