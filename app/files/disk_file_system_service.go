/*
*

---------ATTENTION----------
Code non utilisé pour le moment

*
*/
package files

import (
	"os"
	"peer-to-peer/app/shared"
)

type DiskFileSystemService struct {}

func (d *DiskFileSystemService) Delete(file shared.IFile) error {
	err := os.Remove(file.GetPath())
	if err != nil {
		println("Error deleting file:", err)
		panic(err)
	}

	return nil
}

func (d *DiskFileSystemService) Create(file shared.IFile, content []byte) error {
	f, err := os.Create(file.GetPath())
	if err != nil {
		println("Error creating file:", err)
		panic(err)
	}
	defer f.Close()

	_, err = f.Write(content)
	if err != nil {
		println("Error writing to file:", err)
		panic(err)
	}

	return nil
}