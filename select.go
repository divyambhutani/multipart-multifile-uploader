package main

import (
	"io/fs"
	"log"
	"os"
	"sync"

	uploader "bitbucket.org/divyam_bhutani/fileuploader/upload"
)

// selectUploader function takes pointer of type *os.File and a directory Name
// This function calculates the size of given file using fileinfo
// It uses Uploader package to call uploadSmallObject is file is smaller than 5mb
// For larger files it calls UploadLargeObject
func selectUploader(file *os.File, dirPath string) {
	fileInfo, _ := file.Stat()
	fileSize := fileInfo.Size()
	size := fileSize / (1024 * 1024)
	if size < 5 {
		// upload as a single file
		uploader.UploadSmallObject(file, dirPath)

	} else {
		// upload chunk by chunk
		uploader.UploadLargeObject(file, dirPath)

	}

}

// SelectUploadDir is called when we have to upload a directory
// it take the pointer to *os.File for that directory and path of directory in our system
// it reads all file in that directory one by one and starts uploading by calling selectUploader
func selectUploaderDir(f *os.File, fPath string) {
	// now we know that this file is dir
	// we are going to list all files and call selectUploader
	// on each one by one
	files, err := f.ReadDir(-1)
	if err != nil {
		log.Fatal(err)
	}

	fileinfo, _ := f.Stat()
	dirName := fileinfo.Name()
	for _, file := range files {
		path := fPath + "/" + file.Name()
		// fmt.Println(path)
		curFile, err := os.Open(path)
		if err != nil {
			log.Println("Error opening-", path, ". Continuing Forward")
			continue
		}

		selectUploader(curFile, dirName+"/")
		curFile.Close()
	}
}

func selectUploaderDirSample(f *os.File, fPath string) {
	// now we know that this file is dir
	// we are going to list all files and call selectUploader
	// on each one by one
	files, err := f.ReadDir(-1)
	if err != nil {
		log.Fatal(err)
	}

	fileinfo, _ := f.Stat()
	dirName := fileinfo.Name()
	// file1 := files[0].Name()
	// file2 := files[1].Name()
	// fmt.Println(file1, file2)

	wg := &sync.WaitGroup{}
	ch := make(chan int)
	wg.Add(len(files))
	countAtInstant := 0
	const maxCount = 5
	for _, file := range files {
		if countAtInstant >= maxCount {
			countAtInstant -= (<-ch)
		}
		go func(wg *sync.WaitGroup, ch chan int, file fs.DirEntry, dirName string) {
			defer func() {
				wg.Done()
				ch <- 1
			}()
			path := fPath + "/" + file.Name()
			// fmt.Println(path)
			curFile, _ := os.Open(path)
			if err != nil {
				log.Println("Error opening-", path, ". Continuing Forward")
				return
			}
			selectUploader(curFile, dirName+"/")
			curFile.Close()
		}(wg, ch, file, dirName)
		countAtInstant++
	}
	wg.Wait()
}
