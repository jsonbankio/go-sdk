package jsonbank

import (
	"errors"
	"fmt"
	"github.com/jsonbankio/go-sdk/types"
	"testing"
)

type TestFile struct {
	Id   string
	Path string
}

func TestAuthenticatedMethods(t *testing.T) {
	var jsb = Init(Config{
		Host: "http://localhost:2221",
		Keys: Keys{
			Public:  "pub_wSef-7nVXxvW07hT9tw0_IaHTfepODYNKAqRQCibd7zypIntuzb2hy3r",
			Private: "prv_XuQ8y_ycmO53dLy7JWL0bu-aj_4k2Bi2pW0coVBGoRd0fZxU6WJ26Kaa",
		},
	})

	const project = "sdk-test"
	var testFile = TestFile{"", fmt.Sprintf("%v/index.json", project)}
	const testFileContent = `{
    	"name": "JsonBank SDK Test File",
    	"author": "jsonbank"
	}`

	// Get test file Id
	meta, err := jsb.GetOwnDocumentMeta(testFile.Path)
	if err != nil {
		if err.Code == "notFound" {
			t.Error(errors.New("Test document not found. Please create a document with the content below at {" + testFile.Path + "} before running tests."))
		} else {
			t.Error(err)
		}
		return
	}

	testFile.Id = meta.Id

	t.Run("Authenticate", func(t *testing.T) {
		authenticate, err := jsb.Authenticate()
		if err != nil {
			t.Error(err)
		}
		fmt.Println("Authenticated as: ", authenticate.Username)
	})

	t.Run("Authenticated", func(t *testing.T) {
		if !jsb.Authenticated() {
			t.Error("User is not authenticated")
		}
	})

	t.Run("GetOwnContent", func(t *testing.T) {
		content, err := jsb.GetOwnContent(testFile.Id)
		if err != nil {
			t.Error(err)
		}
		formattedContent := content.(map[string]interface{})

		if formattedContent["author"] != "jsonbank" {
			t.Error("Content does not match")
		}
	})

	t.Run("GetOwnContentByPath", func(t *testing.T) {
		content, err := jsb.GetOwnContentByPath(testFile.Path)
		if err != nil {
			t.Error(err)
		}

		formattedContent := content.(map[string]interface{})

		if formattedContent["author"] != "jsonbank" {
			t.Error("Content does not match")
		}
	})

	t.Run("GetOwnDocumentMeta", func(t *testing.T) {
		meta, err := jsb.GetOwnDocumentMeta(testFile.Id)
		if err != nil {
			t.Error(err)
		}

		if meta.Id != testFile.Id {
			t.Error("Meta does not match")
		}
	})

	t.Run("GetOwnDocumentMetaByPath", func(t *testing.T) {
		meta, err := jsb.GetOwnDocumentMetaByPath(testFile.Path)
		if err != nil {
			t.Error(err)
		}

		if meta.Id != testFile.Id {
			t.Error("Meta does not match")
		}
	})

	t.Run("DeleteDocument And CreateDocument", func(t *testing.T) {
		t.Run("DeleteDocument", func(t *testing.T) {
			_, err := jsb.DeleteDocument(testFile.Path)
			if err != nil {
				t.Error(err)
			}
		})

		t.Run("CreateDocument", func(t *testing.T) {
			_, err = jsb.CreateDocument(types.CreateDocumentBody{
				Name:    "index.json",
				Content: testFileContent,
				Project: project,
			})

			if err != nil {
				t.Error(err)
				return
			}

			testFile.Id = meta.Id
		})
	})

	t.Run("CreateDocumentIfNotExists", func(t *testing.T) {
		document, err := jsb.CreateDocumentIfNotExists(types.CreateDocumentBody{
			Name:    "index.json",
			Content: testFileContent,
			Project: project,
		})

		if err != nil {
			t.Error(err)
			return
		}

		if document.Id != testFile.Id {
			testFile.Id = document.Id
		}
	})

	t.Run("HasOwnDocument", func(t *testing.T) {
		// try by id
		exists := jsb.HasOwnDocument(testFile.Id)
		if !exists {
			t.Error("Document does not exist")
		}

		// try by path
		exists = jsb.HasOwnDocument(testFile.Path)
		if !exists {
			t.Error("Document does not exist")
		}
	})

	t.Run("UpdateOwnDocument", func(t *testing.T) {
		res, err := jsb.UpdateOwnDocument(testFile.Id, `{
    		"name": "JsonBank SDK Test File",
    		"author": "jsonbank", 
			"updated": true
		}`)

		if err != nil {
			t.Error(err)
		}

		if res.Changed != true {
			t.Error("Document was not updated")
		}

		// revert changes
		_, _ = jsb.UpdateOwnDocument(testFile.Id, testFileContent)
	})

	t.Run("CreateFolder", func(t *testing.T) {
		folder, err := jsb.CreateFolder(types.CreatedFolderBody{
			Name:    "folder",
			Project: project,
		})

		if err != nil {
			if err.Code == "name.exists" {
				fmt.Println(err.Error())
			} else {
				t.Error(err)
			}
			return
		}

		if folder.Name != "folder" || folder.Project != project {
			t.Error("New folder data mismatch")
		}
	})

	t.Run("UploadDocument", func(t *testing.T) {
		// delete test file.
		_, _ = jsb.DeleteDocument("sdk-test/upload.json")

		filePath := "./tests/upload.json"
		document, err := jsb.UploadDocument(types.UploadDocumentBody{
			FilePath: filePath,
			Project:  project,
		})

		if err != nil {
			t.Error(err)
			return
		}

		if document.Path != "upload.json" {
			t.Error("Document name mismatch")
		}
	})
}
