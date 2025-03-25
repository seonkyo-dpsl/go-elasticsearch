package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensearch-project/opensearch-go/v2"
	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"
)

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source map[string]interface{} `json:"_source"`
			ID     string                 `json:"_id"`
			Score  float64                `json:"_score"`
		} `json:"hits"`
	} `json:"hits"`
}

func main() {
	// OpenSearch 클라이언트 초기화
	client, err := opensearch.NewClient(opensearch.Config{
		Addresses: []string{"https://localhost:9200/"},
		Password:  "passwd",
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	})
	if err != nil {
		fmt.Printf("Error creating the client: %s\n", err)
		return
	}

	indexName := "v3_contact"

	// 문서 검색
	search := opensearchapi.SearchRequest{
		Index: []string{indexName},
		Body: strings.NewReader(`{
            "query": {
                "match": {
                    "last_name": "jamil"
                }
            }
        }`),
	}

	searchResponse, err := search.Do(context.Background(), client)
	if err != nil {
		fmt.Printf("Error searching documents: %s\n", err)
		return
	}
	defer searchResponse.Body.Close()

	// JSON 응답 파싱
	var searchResult SearchResponse
	if err := json.NewDecoder(searchResponse.Body).Decode(&searchResult); err != nil {
		fmt.Printf("Error parsing the response body: %s\n", err)
		return
	}

	// 검색 결과 출력
	for _, hit := range searchResult.Hits.Hits {
		fmt.Printf("Document ID: %s\n", hit.ID)
		fmt.Printf("Document Source: %+v\n", hit.Source)
	}
}
