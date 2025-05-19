package shared

type IDirectory interface {
	HasFile(fileName string) bool
	GetFiles() map[string]IFile
	GetFile(fileName string) (IFile, bool)
	AddFile(file IFile)
	RemoveFile(fileName string)
}

type Directory struct {
	DirectoryPath string
	Files         map[string]IFile
}

func (dirStat *Directory) HasFile(fileName string) bool {
	_, exists := dirStat.Files[fileName]
	return exists
}

func (dirStat *Directory) GetFiles() map[string]IFile {
	return dirStat.Files
}

func (dirStat *Directory) GetFile(fileName string) (IFile, bool) {
	file, exists := dirStat.Files[fileName]
	return file, exists
}

func (dirStat *Directory) AddFile(file IFile) {
	dirStat.Files[file.GetName()] = file
}

func (dirStat *Directory) RemoveFile(fileName string) {
	delete(dirStat.Files, fileName)
}

func NewDirectory(directoryPath string, files map[string]IFile) IDirectory {
	return &Directory{
		DirectoryPath: directoryPath,
		Files:         files,
	}
}