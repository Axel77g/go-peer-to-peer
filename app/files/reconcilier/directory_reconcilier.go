package file_reconcilier

import "peer-to-peer/app/shared"


func Reconcile(localDirectory, wantedDirectory shared.IDirectory) []Task {
	tasks := []Task{}

    localFiles := localDirectory.GetFiles()  
    wantedFiles := wantedDirectory.GetFiles()

    for name, wantedFile := range wantedFiles {
        localFile, exists := localDirectory.GetFile(name)

        if !exists {
            tasks = append(tasks, Task{
                Type:     TaskDownload,
				File: wantedFile,
            })
        } else {
            if localFile.GetChecksum() != wantedFile.GetChecksum() {
                tasks = append(tasks, Task{
                    Type:     TaskDeleteAndUpload,
                    File: wantedFile,
                })
            }
        }
    }

    for name, localFile := range localFiles {
        if _, exists := wantedDirectory.GetFile(name); !exists {
            tasks = append(tasks, Task{
                Type:     TaskDelete,
                File: localFile,
            })
        }
    }

    return tasks
}