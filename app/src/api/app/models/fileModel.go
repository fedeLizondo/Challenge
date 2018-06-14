package models

type File struct {
	ID          string `json:"id,omitempty"`
	Title       string `json:"titulo"`
	Description string `json:"descripcion"`
}

// type FileServiceInterface interface {
// 	File(id string) (*File, error)
// 	Files() ([]*File, error)
// 	CreateFile(i *File) error
// 	DeleteFile(id string) error
// }
