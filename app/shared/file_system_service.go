package shared

type IFileSystemService interface {
	Delete(file IFile) error
	Create(file IFile, content []byte) error
}
