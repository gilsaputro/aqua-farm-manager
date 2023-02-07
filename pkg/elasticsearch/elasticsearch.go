package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v7"
)

// ElasticSearchMethod is list all available method for elastic search
type ElasticSearchMethod interface {
}

// Client is a wrapper for Elastic Search client
type Client struct {
	es *elasticsearch.Client
}

// NewPostgresClient is func to create ES client
func CreateESClient(address string) (ElasticSearchMethod, error) {
	// Connect to Elasticsearch
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{address},
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		es: es,
	}, nil
}
