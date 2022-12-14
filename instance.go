package jsonbank

import (
	"github.com/jsonbankio/go-sdk/types"
)

type Instance struct {
	config Config         // Instance Config
	memory map[string]any // Instance memory
	urls   struct {
		v1     string // v1 url
		public string // public url
	}
}

// SetHost - switch the host of the instance
func (jsb *Instance) SetHost(host string) {
	jsb.urls.v1 = host + "/v1"
	jsb.urls.public = host
	jsb.config.Host = host
}

// GetContent - get public content from jsonbank
func (jsb *Instance) GetContent(idOrPath string) (any, *RequestError) {
	req, err := jsb.makePublicRequest("GET", jsb.urls.public+"/f/"+idOrPath, nil)
	if err != nil {
		return nil, err
	}

	// send request
	data, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetContentByPath - get public content  by path
// This is only but a syntactic sugar for GetContent by path
func (jsb *Instance) GetContentByPath(path string) (any, *RequestError) {
	return jsb.GetContent(path)
}

// GetDocumentMeta - get public document meta
func (jsb *Instance) GetDocumentMeta(idOrPath string) (*types.DocumentMeta, *RequestError) {
	req, err := jsb.makePublicRequest("GET", jsb.urls.public+"/meta/f/"+idOrPath, nil)
	if err != nil {
		return nil, err
	}

	// send request
	d, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	data := d.(map[string]any)

	return &types.DocumentMeta{
		Id:        data["id"].(string),
		Path:      data["path"].(string),
		Project:   data["project"].(string),
		CreatedAt: data["createdAt"].(string),
		UpdatedAt: data["updatedAt"].(string),
	}, nil
}

// GetDocumentMetaByPath - get public document meta by path
// This is only but a syntactic sugar for GetDocumentMeta by path
func (jsb *Instance) GetDocumentMetaByPath(path string) (*types.DocumentMeta, *RequestError) {
	return jsb.GetDocumentMeta(path)
}

// GetGithubContent - get public content from GitHub
func (jsb *Instance) GetGithubContent(idOrPath string) (any, *RequestError) {
	req, err := jsb.makePublicRequest("GET", jsb.urls.public+"/gh/"+idOrPath, nil)
	if err != nil {
		return nil, err
	}

	// send request
	data, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}
