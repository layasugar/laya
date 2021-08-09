package gstore

import (
	"github.com/olivere/elastic/v7"
)

// olivere/elastic/v7
// this is demo
func NewEsClient(addr []string, user, pwd string) (*elastic.Client, error) {
	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(addr...),
		elastic.SetBasicAuth(user, pwd),
	)
	if err != nil {
		return nil, err
	}
	return client, nil
}
