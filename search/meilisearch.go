package search

import "github.com/meilisearch/meilisearch-go"

type Meili struct {
	idx *meilisearch.Index
}

type MeiliCfg struct {
	// default: http://localhost:7700
	Host string `toml:"host"`
	// default: uask
	Index string `toml:"index"`
	// default: id
	PrimaryKey string `toml:"primary_key"`
}

func NewMeili(cfg *MeiliCfg) (*Meili, error) {
	cli := meilisearch.NewClient(meilisearch.ClientConfig{Host: cfg.Host})
	_, err := cli.CreateIndex(&meilisearch.IndexConfig{Uid: cfg.Index, PrimaryKey: cfg.PrimaryKey})
	if err != nil {
		return nil, err
	}
	return &Meili{cli.Index(cfg.Index)}, nil
}

func (m *Meili) AddDoc(i interface{}) error {
	_, err := m.idx.AddDocuments(i)
	return err
}

func (m *Meili) SearchDoc(query string) ([]interface{}, error) {
	resp, err := m.idx.Search(query, &meilisearch.SearchRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Hits, nil
}

func (m *Meili) UpdateDoc(id string, i interface{}) error {
	_, err := m.idx.UpdateDocuments(i, id)
	return err
}

func (m *Meili) DeleteDoc(id string) error {
	_, err := m.idx.DeleteDocument(id)
	return err
}
