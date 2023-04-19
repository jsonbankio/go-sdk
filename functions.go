package jsonbank

import (
	"bytes"
	"encoding/json"
	"github.com/jsonbankio/go-sdk/types"
	"io"
)

// MakeDocumentPath - generate a document full path
func MakeDocumentPath(document types.CreateDocumentBody) string {
	folder := ""
	if document.Folder != "" {
		folder = document.Folder + "/"
	}
	return document.Project + "/" + folder + document.Name
}

// MakeFolderPath - generate a folder full path
func MakeFolderPath(folder types.CreateFolderBody) string {
	f := ""
	if folder.Folder != "" {
		f = folder.Folder + "/"
	}
	return folder.Project + "/" + f + folder.Name
}

// IsValidJsonString - check if a string is  a valid json string
func IsValidJsonString(s string) bool {
	return json.Valid([]byte(s))
}

// JsonToReader - convert json string to io.Reader
func JsonToReader(s any) io.Reader {
	body, _ := json.Marshal(s)
	return bytes.NewReader(body)
}
