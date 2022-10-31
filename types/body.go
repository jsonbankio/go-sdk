package types

type CreateDocumentBody struct {
	Name    string `json:"name"`
	Project string `json:"project"`
	Folder  string `json:"folder"`
	Content string `json:"content"`
}
