package uploader

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"time"

	"bitbucket.org/divyam_bhutani/file_uploader/color"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func createS3Session() *s3.S3 {
	session := s3.New(session.Must(session.NewSession(&aws.Config{
		Region: aws.String(os.Getenv("AWS_REGION")),
	})))
	return session
}

func UploadLargeObject(file *os.File, dirName string) {
	// read file chunk by chunk

	fileinfo, _ := file.Stat()
	session := createS3Session()

	expiryDate := time.Now().AddDate(0, 0, 1)

	createRes, err := session.CreateMultipartUpload(&s3.CreateMultipartUploadInput{
		Bucket:  aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:     aws.String(dirName + fileinfo.Name()),
		Expires: &expiryDate,
	})

	if err != nil {
		panic(err)
	}
	printMessage := fmt.Sprintf("File <%s> started to upload... ", fileinfo.Name())
	log.Printf(color.LIGHTBLUE, printMessage)

	// setting chuncksize and len
	chunkSize := 1024 * 1024 * 5 // 5mb
	buf := make([]byte, chunkSize)
	totalChuncks := math.Ceil(float64(fileinfo.Size()) / float64(chunkSize))

	fmt.Println("Number of chuncks - ", totalChuncks)

	chunckCount := 1

	var completedParts []*s3.CompletedPart

	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			log.Printf(color.RED, err)
		}
		if err == io.EOF {
			break
		}
		completed, err := uploadChunck(session, createRes, buf[:n], chunckCount)
		if err != nil {
			// that means we have to abort
			_, err := session.AbortMultipartUpload(&s3.AbortMultipartUploadInput{
				UploadId: createRes.UploadId,
				Bucket:   createRes.Bucket,
				Key:      createRes.Key,
			})
			if err != nil {
				log.Printf(color.RED, err)
				return
			}
		}

		completedParts = append(completedParts, completed)
		percentProgress := math.Ceil(float64(chunckCount) * 100 / float64(totalChuncks))
		fmt.Printf("...<%s>..Progress => %v%%  \n", fileinfo.Name(), percentProgress)
		chunckCount++
	}
	_, err = session.CompleteMultipartUpload(&s3.CompleteMultipartUploadInput{
		Bucket:   createRes.Bucket,
		Key:      createRes.Key,
		UploadId: createRes.UploadId,
		MultipartUpload: &s3.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})
	if err != nil {
		log.Printf(color.RED, err)
		return
	}
	// log.Println(completeRes)
	printMessage = fmt.Sprintf(fileinfo.Name() + " uploaded")
	log.Printf(color.GREEN, printMessage)

}

func uploadChunck(session *s3.S3, createRes *s3.CreateMultipartUploadOutput, part []byte, partCount int) (*s3.CompletedPart, error) {
	uploadRes, err := session.UploadPart(&s3.UploadPartInput{
		Body:          bytes.NewReader(part),
		Bucket:        createRes.Bucket,
		Key:           createRes.Key,
		PartNumber:    aws.Int64(int64(partCount)),
		UploadId:      createRes.UploadId,
		ContentLength: aws.Int64(int64(len(part))),
	})
	if err != nil {
		return nil, err
	}
	completedPart := &s3.CompletedPart{
		ETag:       uploadRes.ETag,
		PartNumber: aws.Int64(int64(partCount)),
	}
	return completedPart, nil
}
