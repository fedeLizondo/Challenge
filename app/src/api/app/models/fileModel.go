package models

type File struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"titulo"`
	Description string `json:"descripcion"`
}
