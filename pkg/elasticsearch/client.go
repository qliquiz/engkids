package elasticsearch

/*
import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
)

// Client обертка для клиента Elasticsearch
type Client struct {
	es *elasticsearch.Client
}

// NewClient создает новый клиент для работы с Elasticsearch
func NewClient() (*Client, error) {
	esURL := os.Getenv("ELASTICSEARCH_URL")
	if esURL == "" {
		esURL = "http://elasticsearch:9200"
	}

	cfg := elasticsearch.Config{
		Addresses: []string{esURL},
	}

	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	res, err := es.Info()
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// Проверяем статус ответа
	if res.IsError() {
		return nil, fmt.Errorf("elasticsearch error: %s", res.String())
	}

	return &Client{es: es}, nil
}

// Search выполняет поиск в Elasticsearch
func (c *Client) Search(index string, query map[string]interface{}) (map[string]interface{}, error) {
	var buf strings.Builder

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return nil, err
	}

	res, err := c.es.Search(
		c.es.Search.WithContext(context.Background()),
		c.es.Search.WithIndex(index),
		c.es.Search.WithBody(strings.NewReader(buf.String())),
		c.es.Search.WithTrackTotalHits(true),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return nil, fmt.Errorf("search error: %s", res.String())
	}

	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// IndexDocument индексирует документ в Elasticsearch
func (c *Client) IndexDocument(index string, document map[string]interface{}) error {
	var buf strings.Builder
	if err := json.NewEncoder(&buf).Encode(document); err != nil {
		return err
	}

	res, err := c.es.Index(
		index,
		strings.NewReader(buf.String()),
		c.es.Index.WithContext(context.Background()),
		c.es.Index.WithRefresh("true"),
	)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.IsError() {
		return fmt.Errorf("index error: %s", res.String())
	}

	return nil
}
*/
