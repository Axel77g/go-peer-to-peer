package shared

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
)

type IFile interface {
	GetName() string
	GetPath() string
	GetChecksum() *string
}

/*
*  File implementation file_shared.IFile
 */
type File struct {
	Name string
	Path string
	Checksum string
}

func (f *File) GetPath() string {
	return f.Path
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetChecksum() *string {
	if f.Checksum == "" {
		return nil
	}
	return &f.Checksum
}

func GetFileChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
    if err != nil {
        return "", err
    }
    defer file.Close()

    hasher := md5.New()
    if _, err := io.Copy(hasher, file); err != nil {
        return "", err
    }

    return hex.EncodeToString(hasher.Sum(nil)), nil
}

func NewFile(name string, path string) (IFile , error) {
	checksum, err := GetFileChecksum(path)
	if err != nil {
		return &File{}, err
	}

	return &File{
		Name: name,
		Path: path,
		Checksum: checksum,
	}, nil
}
