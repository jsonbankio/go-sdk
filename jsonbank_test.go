package jsonbank

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/jsonbankio/go-sdk/types"
	"os"
	"testing"
)

type TestFile struct {
	Id   string
	Path string
}

const testFileContent = `{
	"name": "JsonBank SDK Test File",
	"author": "jsonbank"
}`

// read .env file
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		// if CI is true, we are running tests on github actions, so we don't need to panic
		if os.Getenv("CI") != "true" {
			panic(err)
		}
	}
}

func TestNotAuthenticated(t *testing.T) {
	loadEnv()

	var jsb = InitWithoutKeys()
	jsb.SetHost(os.Getenv("JSB_HOST"))

	const project = "jsonbank/sdk-test"
	var testFile = TestFile{"", fmt.Sprintf("%v/index.json", project)}

	// Get test file Id
	meta, err := jsb.GetDocumentMetaByPath(testFile.Path)
	if err != nil {
		if err.Code == "notFound" {
			t.Error("Test document not found. Please create a document with the content below at {" + testFile.Path + "} before running tests.\n" + testFileContent)
		} else {
			t.Error(err)
		}
		return
	}
	testFile.Id = meta.Id

	t.Run("GetContent", func(t *testing.T) {
		document, err := jsb.GetContent(testFile.Id)
		if err != nil {
			t.Error(err)
			return
		}

		// convert to map
		data := document.(map[string]interface{})

		if data["author"] != "jsonbank" {
			t.Error(errors.New("GetContent should return a valid document"))
		}
	})

	t.Run("GetContentByPath", func(t *testing.T) {
		document, err := jsb.GetContentByPath(testFile.Path)

		if err != nil {
			t.Error(err)
			return
		}

		// convert to map
		data := document.(map[string]interface{})

		if data["author"] != "jsonbank" {
			t.Error(errors.New("GetContentByPath should return a valid document"))
		}
	})

	t.Run("GetGithubContent", func(t *testing.T) {
		data, err := jsb.GetGithubContent("jsonbankio/jsonbank-js/package.json")

		if err != nil {
			t.Error(err)
			return
		}

		// convert to map
		pkg := data.(map[string]interface{})

		if pkg["name"] != "jsonbank" {
			t.Error("GetGithubContent should return a valid document")

		} else if pkg["author"] != "jsonbankio" {
			t.Error("GetGithubContent should return a valid document")
		}
	})

}

func TestAuthenticated(t *testing.T) {
	loadEnv()

	var jsb = Init(Config{
		Host: os.Getenv("JSB_HOST"),
		Keys: Keys{
			Public:  os.Getenv("JSB_PUBLIC_KEY"),
			Private: os.Getenv("JSB_PRIVATE_KEY"),
		},
	})

	const project = "sdk-test"
	var testFile = TestFile{"", fmt.Sprintf("%v/index.json", project)}

	// Get test file Id
	meta, err := jsb.GetOwnDocumentMeta(testFile.Path)
	if err != nil {
		if err.Code == "notFound" {
			t.Error("Test document not found. Please create a document with the content below at {" + testFile.Path + "} before running tests.\n" + testFileContent)
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
		folder, err := jsb.CreateFolder(types.CreateFolderBody{
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

	t.Run("CreateFolderIfNotExists", func(t *testing.T) {
		name := "folder"
		folder, err := jsb.CreateFolderIfNotExists(types.CreateFolderBody{
			Name:    name,
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

		if folder.Name != name || folder.Project != project {
			t.Error("New folder data mismatch")
		}
	})

	t.Run("GetFolder", func(t *testing.T) {
		folderPath := project + "/folder"
		folder, err := jsb.GetFolder(folderPath)
		if err != nil {
			t.Error(err)
		}

		// project must match
		if folder.Project != project {
			t.Error("Folder project mismatch")
		}

		// try to get folder by id
		folder, err = jsb.GetFolder(folder.Id)
		if err != nil {
			t.Error(err)
		}

		// project must match
		if folder.Project != project {
			t.Error("Folder project mismatch")
		}

	})

	t.Run("GetFolderWithStats", func(t *testing.T) {
		folderPath := project + "/folder"
		folder, err := jsb.GetFolderWithStats(folderPath)
		if err != nil {
			t.Error(err)
		}

		// project must match
		if folder.Project != project {
			t.Error("Folder project mismatch")
		}

		// check that stats exist
		if folder.Stats == nil {
			t.Error("Folder stats are nil")
		}

		// try to get folder by id
		folder, err = jsb.GetFolderWithStats(folder.Id)
		if err != nil {
			t.Error(err)
		}

		// project must match
		if folder.Project != project {
			t.Error("Folder project mismatch")
		}

		// check that stats exist
		if folder.Stats == nil {
			t.Error("Folder stats are nil")
		}
	})

	t.Run("UploadDocument To New Folder", func(t *testing.T) {
		// delete test file.
		_, _ = jsb.DeleteDocument("sdk-test/folder/upload.json")

		// Upload file to new folder
		filePath := "./tests/upload.json"
		document, err := jsb.UploadDocument(types.UploadDocumentBody{
			FilePath: filePath,
			Project:  project,
			Folder:   "folder",
		})

		if err != nil {
			t.Error(err)
			return
		}

		if document.Path != "folder/upload.json" {
			t.Error("Document name mismatch")
		}
	})
}
