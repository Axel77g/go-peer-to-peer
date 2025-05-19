package file_reconcilier

import (
	"peer-to-peer/app/shared"
)

type TaskType int

const (
	TaskDownload TaskType = iota
	TaskDeleteAndUpload
	TaskDelete
)

type Task struct {
    Type     TaskType
    File  shared.IFile
}


func ExecuteTask(task Task, downloader shared.IDownloaderService, fileSystem shared.IFileSystemService) error {
	switch task.Type {
	case TaskDownload:
		return downloader.Download(task.File)
	case TaskDeleteAndUpload:
		downloader.Download(task.File)
		return fileSystem.Delete(task.File)
	case TaskDelete:
		return fileSystem.Delete(task.File)
	}
	return nil
}