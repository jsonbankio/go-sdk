package types

type CreateDocumentBody struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Folder  string `json:"folder"`
	Content string `json:"content"`
}

type CreatedFolderBody struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Folder  string `json:"folder"`
}

type UploadDocumentBody struct {
	FilePath string `json:"file"`
	Project  string `json:"project"`
	Name     string `json:"name"`
	Folder   string `json:"folder"`
}
