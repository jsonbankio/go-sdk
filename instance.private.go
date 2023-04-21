package jsonbank

import (
	"encoding/json"
	"io"
	"net/http"
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

func (jsb *Instance) makePublicRequest(method string, url string, data io.Reader) (*http.Request, *RequestError) {
	req, _ := http.NewRequest(method, url, data)
	req.Header.Add("Content-Type", "application/json")
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

// sendRequestAsText - send request and return response as text
func (jsb *Instance) sendRequestAsText(req *http.Request) (*string, *RequestError) {
	// make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, &RequestError{"request_error", err.Error()}
	}

	// check if request was successful
	if res.StatusCode != 200 {
		// convert response to json
		var data map[string]any
		jsonError := json.NewDecoder(res.Body).Decode(&data)

		if jsonError != nil {
			return nil, &RequestError{"json_error", jsonError.Error()}
		}

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

	// convert res.Body to string
	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, &RequestError{"request_error", err.Error()}
	}

	bodyString := string(bodyBytes)

	return &bodyString, nil
}
