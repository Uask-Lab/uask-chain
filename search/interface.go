package search

type Search interface {
	AddDoc(interface{}) error
	SearchDoc(query string) ([]interface{}, error)
	DeleteDoc(id string) error
}
