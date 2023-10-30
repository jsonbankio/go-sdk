package types

type AuthenticatedKey struct {
	Title    string   `json:"title"`
	Projects []string `json:"projects,omitempty"`
}
type AuthenticatedData struct {
	Authenticated bool   `json:"authenticated"`
	Username      string `json:"username"`
	ApiKey        AuthenticatedKey
}

type NewDocument struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Project   string `json:"project"`
	CreatedAt string `json:"createdAt"`
	Exists    bool   `json:"exists"`
}

type FolderStats struct {
	Documents float64 `json:"documents"`
	Folders   float64 `json:"folders"`
}

type Folder struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Project   string `json:"project"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	// optional fields
	Stats *FolderStats `json:"stats,omitempty"`
}

// NewFolder extends Folder
type NewFolder struct {
	Folder `json:",inline"`
	Exists bool `json:"exists"`
}

type ContentSize struct {
	Number float64 `json:"number"`
	String string  `json:"string"`
}

type DocumentMeta struct {
	Id          string      `json:"id"`
	Name        string      `json:"name"`
	Project     string      `json:"project"`
	Path        string      `json:"path"`
	ContentSize ContentSize `json:"contentSize"`
	FolderId    string      `json:"folderId"`
	UpdatedAt   string      `json:"updatedAt"`
	CreatedAt   string      `json:"createdAt"`
}

func DataToDocumentMeta(data map[string]interface{}) *DocumentMeta {

	d := &DocumentMeta{
		Id:      data["id"].(string),
		Path:    data["path"].(string),
		Project: data["project"].(string),
		Name:    data["name"].(string),
		ContentSize: ContentSize{
			Number: data["contentSize"].(map[string]interface{})["number"].(float64),
			String: data["contentSize"].(map[string]interface{})["string"].(string),
		},
		CreatedAt: data["createdAt"].(string),
		UpdatedAt: data["updatedAt"].(string),
	}

	// check if folder exists
	if data["folderId"] != nil {
		d.FolderId = data["folderId"].(string)
	}

	return d
}

type UpdatedDocument struct {
	Changed bool `json:"changed"`
}

type DeletedDocument struct {
	Deleted bool `json:"deleted"`
}
