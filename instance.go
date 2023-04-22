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

// GetContentAsString - get public content from jsonbank as string
func (jsb *Instance) GetContentAsString(idOrPath string) (string, *RequestError) {
	req, err := jsb.makePublicRequest("GET", jsb.urls.public+"/f/"+idOrPath, nil)
	if err != nil {
		return "", err
	}

	// send request
	data, err := jsb.sendRequestAsText(req)

	if err != nil {
		return "", err
	}

	return *data, nil
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

	return types.DataToDocumentMeta(data), nil
}

// GetGithubContent - get public content from GitHub
func (jsb *Instance) GetGithubContent(path string) (any, *RequestError) {
	req, err := jsb.makePublicRequest("GET", jsb.urls.public+"/gh/"+path, nil)
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

// GetGithubContentAsString - get public content from GitHub as string
func (jsb *Instance) GetGithubContentAsString(path string) (string, *RequestError) {
	req, err := jsb.makePublicRequest("GET", jsb.urls.public+"/gh/"+path, nil)
	if err != nil {
		return "", err
	}

	// send request
	data, err := jsb.sendRequestAsText(req)

	if err != nil {
		return "", err
	}

	return *data, nil
}
