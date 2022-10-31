package jsonbank

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"jsonbank/types"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
)

type Keys struct {
	Public  string // Public Key
	Private string // Private Key
}

type Config struct {
	Host string // Server Host
	Keys Keys   // Keys
}

type Instance struct {
	config Config         // Instance Config
	memory map[string]any // Instance memory
	urls   struct {
		v1     string // v1 url
		public string // public url
	}
}

// ==========  Public Methods ==========

// Init - initializes the jsonbank instance
func Init(config Config) Instance {
	// Validate config
	// Assign default Host if not provided
	if len(config.Host) <= 0 {
		config.Host = "https://api.jsonbank.io"
	}

	// make instance
	jsb := Instance{}
	// set config
	jsb.config = config
	// set urls
	jsb.urls.v1 = config.Host + "/v1"
	jsb.urls.public = config.Host + "/"
	// set memory
	jsb.memory = make(map[string]any)

	return jsb
}

// InitWithoutKeys - initializes the jsonbank instance without Keys
//
//goland:noinspection GoUnusedExportedFunction
func InitWithoutKeys() Instance {
	return Init(Config{})
}

// ========== Private Methods ==========
// hasKey - validates the Keys
func (jsb *Instance) hasKey(key string) bool {
	if key == "public" {
		// check if Public key is set
		return len(jsb.config.Keys.Public) > 0
	} else if key == "private" {
		// check if Private key is set
		return len(jsb.config.Keys.Private) > 0
	} else {
		return false
	}
}

// MakePostRequest - make a request with only Public api key
func (jsb *Instance) makeRequest(method string, url string, data io.Reader) (*http.Request, *RequestError) {
	// check if Public key is set
	if !jsb.hasKey("public") {
		return nil, &RequestError{"bad_request", "Public key is not set"}
	}

	req, _ := http.NewRequest(method, url, data)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("jsb-pub-key", jsb.config.Keys.Public)
	return req, nil
}

// MakePrivatePostRequest - make a request with both Public && Private api Keys
func (jsb *Instance) makePrivateRequest(method string, url string, data io.Reader) (*http.Request, *RequestError) {
	req, err := jsb.makeRequest(method, url, data)
	if err != nil {
		return nil, err
	}

	// check if private key is set
	if !jsb.hasKey("private") {
		return nil, &RequestError{"bad_request", "Private key is not set"}
	}
	req.Header.Add("jsb-prv-key", jsb.config.Keys.Private)

	return req, nil
}

func (jsb *Instance) sendRequest(req *http.Request) (any, *RequestError) {
	// make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &RequestError{"request_error", err.Error()}
	}

	// convert response to json
	var data map[string]any
	jsonError := json.NewDecoder(res.Body).Decode(&data)

	if jsonError != nil {
		return nil, &RequestError{"json_error", jsonError.Error()}
	}

	// check if request was successful
	if res.StatusCode != 200 {
		if data["error"] != nil {
			dataError := data["error"]
			// check if dataError is a map
			if reflect.TypeOf(dataError).Kind() == reflect.String {
				return nil, &RequestError{"request_error", dataError.(string)}
			} else if reflect.TypeOf(dataError).Kind() == reflect.Map {
				dataError := dataError.(map[string]any)
				return nil, &RequestError{dataError["code"].(string), dataError["message"].(string)}
			} else {
				return nil, &RequestError{"request_error", "Request was not successful"}
			}
		} else {
			return nil, &RequestError{"request_error", "Request was not successful"}
		}
	}

	return data, nil
}

// ==========  Public Instance Methods ==========

// Authenticate - authenticates the jsonbank instance
func (jsb *Instance) Authenticate() (*types.AuthenticatedData, *RequestError) {
	url := jsb.urls.v1 + "/authenticate"
	req, _ := jsb.makeRequest("POST", url, nil)

	// make request
	d, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	jsb.memory["authenticated"] = true
	data := d.(map[string]any)

	authenticatedData := types.AuthenticatedData{
		Authenticated: data["authenticated"].(bool),
		Username:      data["username"].(string),
		ApiKey: types.AuthenticatedKey{
			Title:    data["apiKey"].(map[string]interface{})["title"].(string),
			Projects: data["apiKey"].(map[string]interface{})["projects"].(string),
		},
	}

	jsb.memory["authenticatedData"] = authenticatedData

	return &authenticatedData, nil
}

// Authenticated - checks if the jsonbank instance is authenticated
func (jsb *Instance) Authenticated() bool {
	return jsb.memory["authenticated"] == true
}

// GetUsername - gets the username of the authenticated user
func (jsb *Instance) GetUsername() string {
	if !jsb.Authenticated() {
		return ""
	}
	return jsb.memory["authenticatedData"].(types.AuthenticatedData).Username
}

// GetOwnContent - gets the content of a document owned by the authenticated user
func (jsb *Instance) GetOwnContent(idOrPath string) (any, *RequestError) {
	req, err := jsb.makeRequest("GET", jsb.urls.v1+"/file/"+idOrPath, nil)
	if err != nil {
		return nil, err
	}

	// make request
	data, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GetOwnContentByPath - gets the content (by path) of a document owned by the authenticated user
// This is only but a syntactic sugar for GetOwnContent by path
func (jsb *Instance) GetOwnContentByPath(path string) (any, *RequestError) {
	return jsb.GetOwnContent(path)
}

// GetOwnDocumentMeta - gets the content meta of the authenticated user
func (jsb *Instance) GetOwnDocumentMeta(idOrPath string) (*types.DocumentMeta, *RequestError) {
	req, err := jsb.makeRequest("GET", jsb.urls.v1+"/meta/file/"+idOrPath, nil)
	if err != nil {
		return nil, err
	}

	// make request
	d, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	data := d.(map[string]any)

	return &types.DocumentMeta{
		Id:        data["id"].(string),
		Project:   data["project"].(string),
		Path:      data["path"].(string),
		UpdatedAt: data["updatedAt"].(string),
		CreatedAt: data["createdAt"].(string),
	}, nil
}

// GetOwnDocumentMetaByPath - gets the content meta (by path) of the authenticated user
// This is only but a syntactic sugar for GetOwnDocumentMeta by path
func (jsb *Instance) GetOwnDocumentMetaByPath(path string) (*types.DocumentMeta, *RequestError) {
	return jsb.GetOwnDocumentMeta(path)
}

// CreateDocument - creates a document
func (jsb *Instance) CreateDocument(document types.CreateDocumentBody) (*types.NewDocument, *RequestError) {
	// project is required
	if document.Project == "" {
		return nil, &RequestError{"bad_request", "Project is required"}
	}
	// name is required
	if document.Name == "" {
		return nil, &RequestError{"bad_request", "Name is required"}
	}

	url := fmt.Sprintf("/project/%s/document", document.Project)

	// check if content is a valid json string
	if !IsValidJsonString(document.Content) {
		return nil, &InvalidJsonError
	}

	// convert document to reader
	body, _ := json.Marshal(document)

	// send request
	req, err := jsb.makePrivateRequest("POST", jsb.urls.v1+url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	d, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	data := d.(map[string]any)

	return &types.NewDocument{
		Id:        data["id"].(string),
		Name:      data["name"].(string),
		Path:      data["path"].(string),
		Project:   data["project"].(string),
		CreatedAt: data["createdAt"].(string),
		Exists:    data["exists"] == true,
	}, nil
}

// UploadDocument - uploads a json document
func (jsb *Instance) UploadDocument(document types.UploadDocumentBody) (*types.NewDocument, *RequestError) {
	// project is required
	if document.Project == "" {
		return nil, &RequestError{"bad_request", "Project is required"}
	}

	// check if file exists
	if _, err := os.Stat(document.FilePath); os.IsNotExist(err) {
		return nil, &RequestError{"file_not_found", "File does not exist"}
	}

	// get content of file
	content, err := os.ReadFile(document.FilePath)
	if err != nil {
		return nil, &RequestError{"invalid_file", "Could not read file"}
	}

	// check if content is a valid json string
	if !IsValidJsonString(string(content)) {
		return nil, &InvalidJsonError
	}

	// set name if not set
	if document.Name == "" {
		document.Name = filepath.Base(document.FilePath)
	}

	// create document
	return jsb.CreateDocument(types.CreateDocumentBody{
		Project: document.Project,
		Name:    document.Name,
		Content: string(content),
		Folder:  document.Folder,
	})
}

// CreateDocumentIfNotExists - creates a document if it does not exist
func (jsb *Instance) CreateDocumentIfNotExists(document types.CreateDocumentBody) (*types.NewDocument, *RequestError) {
	data, err := jsb.CreateDocument(document)
	if err != nil {
		// if code is "name.exists" then fetch content meta
		if err.Code == "name.exists" {
			meta, err := jsb.GetOwnDocumentMeta(MakeDocumentPath(document))
			if err != nil {
				return nil, err
			}

			return &types.NewDocument{
				Id:        meta.Id,
				Name:      document.Name,
				Path:      meta.Path,
				Project:   meta.Project,
				CreatedAt: meta.CreatedAt,
				Exists:    true,
			}, nil
		} else {
			return nil, err
		}
	}

	return data, nil
}

// HasOwnDocument - tries to get the content then returns true if it exists
func (jsb *Instance) HasOwnDocument(idOrPath string) bool {
	_, err := jsb.GetOwnDocumentMeta(idOrPath)
	return err == nil
}

// UpdateOwnDocument - Update document owned by the authenticated user
func (jsb *Instance) UpdateOwnDocument(idOrPath string, content string) (*types.UpdatedDocument, *RequestError) {
	// check if content is a valid json string
	if !IsValidJsonString(content) {
		return nil, &InvalidJsonError
	}

	body := JsonToReader(struct {
		Content string `json:"content"`
	}{
		Content: content,
	})

	req, err := jsb.makePrivateRequest("POST", jsb.urls.v1+"/file/"+idOrPath, body)
	if err != nil {
		return nil, err
	}

	// send request
	data, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	// convert to map
	d := data.(map[string]interface{})

	return &types.UpdatedDocument{
		Changed: d["changed"].(bool),
	}, nil
}

// DeleteDocument - deletes a document
func (jsb *Instance) DeleteDocument(idOrPath string) (*types.DeletedDocument, *RequestError) {
	req, err := jsb.makePrivateRequest("DELETE", jsb.urls.v1+"/file/"+idOrPath, nil)
	if err != nil {
		return nil, err
	}

	// send request
	data, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	// convert to map
	d := data.(map[string]interface{})

	return &types.DeletedDocument{
		Deleted: d["deleted"].(bool),
	}, nil
}

// CreateFolder - creates a folder
func (jsb *Instance) CreateFolder(body types.CreatedFolderBody) (*types.NewFolder, *RequestError) {
	// project is required
	if body.Project == "" {
		return nil, &RequestError{"bad_request", "Project is required"}
	}
	// name is required
	if body.Name == "" {
		return nil, &RequestError{"bad_request", "Name is required"}
	}

	url := fmt.Sprintf("/project/%s/folder", body.Project)

	// make request
	req, err := jsb.makePrivateRequest("POST", jsb.urls.v1+url, JsonToReader(body))
	if err != nil {
		return nil, err
	}

	// send request
	data, err := jsb.sendRequest(req)
	if err != nil {
		return nil, err
	}

	// convert to map
	d := data.(map[string]interface{})

	return &types.NewFolder{
		Id:      d["id"].(string),
		Name:    d["name"].(string),
		Path:    d["path"].(string),
		Project: d["project"].(string),
	}, nil
}
