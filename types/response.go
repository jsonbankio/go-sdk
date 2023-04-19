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

type DocumentMeta struct {
	Id        string `json:"id"`
	Project   string `json:"project"`
	Path      string `json:"path"`
	UpdatedAt string `json:"updatedAt"`
	CreatedAt string `json:"createdAt"`
}

type UpdatedDocument struct {
	Changed bool `json:"changed"`
}

type DeletedDocument struct {
	Deleted bool `json:"deleted"`
}
