package elastic

import (
	"context"

	"github.com/olivere/elastic"
	"github.com/v-zhidu/orb/logging"
)

//Config is configuration struct for elasticsearch.
type Config struct {
	Nodes []string
}

//NewElasticClient - the factory for build the client of elasticsearch.
func NewElasticClient(esConfig *Config) (*elastic.Client, error) {
	client, err := elastic.NewClient(elastic.SetURL(esConfig.Nodes...), elastic.SetSniff(false))

	if err != nil {
		logging.Error("failed to connect elasticsearch server", logging.Fields{
			"nodes": esConfig.Nodes,
		}, err)
		return nil, err
	}

	return client, err
}

//SearchElastic - common method to search with elastic.
func SearchElastic(client *elastic.Client, index string, query elastic.Query,
	sortBy string, ascending bool, from int, size int) ([]*elastic.SearchHit, error) {
	searchResult, err := client.Search().
		Index(index).
		Query(query).
		Pretty(false).
		From(from).Size(size).
		Sort(sortBy, ascending).
		Do(context.Background())

	if err != nil {
		return nil, err
	}

	// Here's how you iterate through results with full control over each step.
	if searchResult.Hits.TotalHits > 0 {
		logging.Info("elasticsearch search hits", logging.Fields{
			"counts":      searchResult.Hits.TotalHits,
			"TookInMills": searchResult.TookInMillis,
		})

		return searchResult.Hits.Hits, nil
	}

	// No hits
	logging.Infoln("elasticsearch search no hits")
	return nil, nil
}

//CountElastic - common method to count with elastic.
func CountElastic(client *elastic.Client, index string, query elastic.Query) (int64, error) {
	countResult, err := client.Count().
		Index(index).
		Query(query).
		Do(context.Background())

	if err != nil {
		return 0, err
	}

	return countResult, nil
}
