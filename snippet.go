package main

// Snippet struct
type Snippet struct {
	ID             string      `json:"id"`
	Type           string      `json:"type"`
	ProjectID      string      `json:"project_id"`
	Name           string      `json:"name"`
	Description    interface{} `json:"description"`
	Mergecode      string      `json:"mergecode"`
	Languages      []string    `json:"languages"`
	ContentType    interface{} `json:"content_type"`
	CurrentVersion struct {
		En string `json:"en"`
	} `json:"current_version"`
	Visibility    string      `json:"visibility"`
	ReaderRoles   interface{} `json:"reader_roles"`
	DateCreated   string      `json:"date_created"`
	DateModified  string      `json:"date_modified"`
	DatePublished interface{} `json:"date_published"`
	DateDeleted   interface{} `json:"date_deleted"`
	Status        string      `json:"status"`
	PdfHide       string      `json:"pdf_hide"`
	Blurb         interface{} `json:"blurb"`
}

// APISnippetResponse struct
type APISnippetResponse struct {
	Valid     bool `json:"valid"`
	PageStats struct {
		TotalRecords int `json:"total_records"`
		TotalPages   int `json:"total_pages"`
	} `json:"page_stats"`
	Data []Snippet `json:"data"`
}
