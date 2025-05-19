package shared

type IDownloaderService interface {
	Download(file IFile) error
}