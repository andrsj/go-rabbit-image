package file

type FileRepository interface {
	CreateImage(data []byte, id string, level string) error
	GetImage(id string, level string) ([]byte, error)
}
