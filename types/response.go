package types

type AuthenticatedKey struct {
	Title    string `json:"title"`
	Projects string `json:"projects"`
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
