/*
-> file_uploader helps to upload large files to s3

-> Instructions
  just enter the file/directory path when prompted
 	you can modify the bucket name and aws region in config.env
 	if you want to quit just enter 'q'
-> Specs -
 - this program has 2 types of upload functions
 - for objects < 5mb , it will upload the file in one go
 - for larger objects, it uses multipart upload by calculating the number of chuncks instantly
*/
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("config.env")
}

func main() {
	firstTime := true
	for {
		if !firstTime {
			fmt.Println(`*
*
*`)
		}
		fmt.Println("Enter file/folder path or enter 'q' to quit -")
		var fPath string
		fmt.Scan(&fPath)
		if fPath == "q" {
			break
		}
		selectFile(fPath)
		firstTime = false
	}
}

/*Upload file function no args , you call this function and it will prompt for a path
 -> Instructions :-
 		- you can enter a valid path to a file or folder
    - incase the path is invalid then it will log the error
 		- this func keeps running and prompts you to enter path again after finishing previous uploads
		- enter 'q' to quit and stop execution
*/
// func UploadFile() {
// 	firstTime := true
// 	for {
// 		if !firstTime {
// 			fmt.Println(`*
// *
// *`)
// 		}
// 		fmt.Println("Enter file/folder path or enter 'q' to quit -")
// 		var fPath string
// 		fmt.Scan(&fPath)
// 		if fPath == "q" {
// 			break
// 		}
// 		selectFile(fPath)
// 		firstTime = false
// 	}
// }

// selectFile
func selectFile(fPath string) {

	// open this file
	f, err := os.Open(fPath)

	if err != nil {
		log.Println(err)
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	fileInfo, err := f.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		// if file is a directory then upload file by file
		selectUploaderDirSample(f, fPath)

	} else {
		// if it is a single file then just upload

		selectUploader(f, "")
	}

}
