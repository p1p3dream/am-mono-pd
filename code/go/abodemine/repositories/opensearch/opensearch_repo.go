// Package opensearch provides a repository interface for OpenSearch.
// This facilitates interaction with OpenSearch's HTTP API and prevents
// the use of reflection, which improves readability and performance.
//
// An OpenSearch repository serves a single config key. If multiple OpenSearch
// sources are needed, create multiple repository instances with different
// config keys.
//
// DO NOT add non-OpenSearch-specific types or methods to this package.
package opensearch

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/opensearch-project/opensearch-go/v2/opensearchapi"

	"abodemine/domains/arc"
	"abodemine/lib/errors"
)

type Repository interface {
	ConfigKey() string

	IndexExists(r *arc.Request, in *IndexExistsInput) (*IndexExistsOutput, error)
	CreateIndex(r *arc.Request, in *CreateIndexInput) (*CreateIndexOutput, error)

	PutDocument(r *arc.Request, in *PutDocumentInput) (*PutDocumentOutput, error)
}

type repository struct {
	configKey string
}

type NewRepositoryInput struct {
	ConfigKey string
}

func NewRepository(in *NewRepositoryInput) Repository {
	return &repository{
		configKey: in.ConfigKey,
	}
}

func (repo *repository) ConfigKey() string {
	return repo.configKey
}

type IndexExistsInput struct {
	IndexName string
}

type IndexExistsOutput struct {
	Exists bool
}

func (repo *repository) IndexExists(r *arc.Request, in *IndexExistsInput) (*IndexExistsOutput, error) {
	cli, err := r.Dom().SelectOpenSearch(repo.configKey)
	if err != nil {
		return nil, errors.Forward(err, "9a86bba5-4bf8-41e0-b965-8d7b1b115353")
	}

	req := opensearchapi.IndicesExistsRequest{
		Index: []string{in.IndexName},
	}

	rsp, err := req.Do(r.Context(), cli)
	if err != nil {
		return nil, &errors.Object{
			Id:     "c4fbab7a-9d37-4959-ac95-c6a7c84fd813",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to query opensearch.",
			Cause:  err.Error(),
		}
	}

	defer rsp.Body.Close()
	out := new(IndexExistsOutput)

	switch rsp.StatusCode {
	case http.StatusOK:
		out.Exists = true
	case http.StatusNotFound:
		// Do nothing.
	case http.StatusTooManyRequests:
		return nil, &errors.Object{
			Id:   "e0d0c176-dc79-4814-b15b-97a3904654af",
			Code: errors.Code_RESOURCE_EXHAUSTED,
		}
	default:
		return nil, &errors.Object{
			Id:     "891a8daf-4b5e-4de1-8676-59301586244d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Unknown response from opensearch.",
			Meta: map[string]any{
				"response": rsp.String(),
			},
		}
	}

	return out, nil
}

type CreateIndexInput struct {
	IndexName string
	Body      *IndexCreateBody
}

type CreateIndexOutput struct{}

func (repo *repository) CreateIndex(r *arc.Request, in *CreateIndexInput) (*CreateIndexOutput, error) {
	cli, err := r.Dom().SelectOpenSearch(repo.configKey)
	if err != nil {
		return nil, errors.Forward(err, "fdd2d805-1c1f-4e29-aa00-d085b0a23a66")
	}

	req := opensearchapi.IndicesCreateRequest{
		Index: in.IndexName,
	}

	if in.Body != nil {
		b, err := json.Marshal(in.Body)
		if err != nil {
			return nil, &errors.Object{
				Id:     "09deb0c0-3786-4c90-9285-1f7ee9cec12e",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to encode index create body to JSON.",
				Cause:  err.Error(),
			}
		}

		req.Body = bytes.NewBuffer(b)
	}

	rsp, err := req.Do(r.Context(), cli)
	if err != nil {
		return nil, &errors.Object{
			Id:     "8f7a9b2d-6c3e-4d5a-b1f8-e9d2c4f3a5b6",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to create index in opensearch.",
			Cause:  err.Error(),
		}
	}

	defer rsp.Body.Close()

	if rsp.StatusCode != http.StatusOK {
		return nil, &errors.Object{
			Id:     "2e4d6f8a-1c3b-4e5d-9a7f-8b2c5e4d6f8a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Unknown response from opensearch.",
			Meta: map[string]any{
				"response": rsp.String(),
			},
		}
	}

	out := new(CreateIndexOutput)

	return out, nil
}

type PutDocumentInput struct {
	IndexName string
	Items     []Document
}

type PutDocumentOutput struct{}

func (repo *repository) PutDocument(r *arc.Request, in *PutDocumentInput) (*PutDocumentOutput, error) {
	cli, err := r.Dom().SelectOpenSearch(repo.configKey)
	if err != nil {
		return nil, errors.Forward(err, "65e84bbd-c461-4016-b861-c348067c28ea")
	}

	buf := new(bytes.Buffer)
	jsonEncoder := json.NewEncoder(buf)

	for _, item := range in.Items {
		if err := jsonEncoder.Encode(&opensearchActionMetadata{
			Index: &opensearchActionMetadataIndex{
				Index: in.IndexName,
				Id:    item.OpenSearchId(),
			},
		}); err != nil {
			return nil, &errors.Object{
				Id:     "c0aa0235-fec8-4b83-b30c-eb818e58bc2d",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to encode action metadata to JSON.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"document": item,
				},
			}
		}

		if err := jsonEncoder.Encode(item); err != nil {
			return nil, &errors.Object{
				Id:     "a44ad686-cf36-4e11-8e67-37c2cf250cc9",
				Code:   errors.Code_UNKNOWN,
				Detail: "Failed to encode document to JSON.",
				Cause:  err.Error(),
				Meta: map[string]any{
					"document": item,
				},
			}
		}
	}

	req := opensearchapi.BulkRequest{
		Index:   in.IndexName,
		Body:    buf,
		Refresh: "true",
	}

	rsp, err := req.Do(r.Context(), cli)
	if err != nil {
		return nil, &errors.Object{
			Id:     "7d9e5f3a-2b1c-4d8e-9f6a-3c8b4d5e6f7a",
			Code:   errors.Code_UNKNOWN,
			Detail: "Failed to bulk insert documents in opensearch.",
			Cause:  err.Error(),
		}
	}

	defer rsp.Body.Close()

	if rsp.IsError() || rsp.HasWarnings() || rsp.StatusCode != http.StatusOK {
		if rsp.StatusCode == http.StatusTooManyRequests {
			return nil, &errors.Object{
				Id:     "ba849718-c79d-46b7-8eef-dec39fb5e0d6",
				Code:   errors.Code_RESOURCE_EXHAUSTED,
				Detail: "Too many requests to opensearch.",
			}
		}

		return nil, &errors.Object{
			Id:     "1a2b3c4d-5e6f-7a8b-9c0d-1e2f3a4b5c6d",
			Code:   errors.Code_UNKNOWN,
			Detail: "Unknown response from opensearch.",
			Meta: map[string]any{
				"response": rsp.String(),
			},
		}
	}

	out := new(PutDocumentOutput)

	return out, nil
}
