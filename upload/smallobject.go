package uploader

import (
	"fmt"
	"log"
	"os"

	"bitbucket.org/divyam_bhutani/file_uploader/color"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func createAWSSession() *session.Session {
	// accessId := os.Getenv("AWS_ACCESS_KEY_ID")

	// secretKey := os.Getenv("AWS_SECRET_ACCESS_KEY")
	region := os.Getenv("AWS_REGION")
	// token := os.Getenv("AWS_ACCESS_TOKEN")

	session, err := session.NewSession(
		&aws.Config{
			Region: aws.String(region),
			// Credentials: credentials.NewStaticCredentials(accessId, secretKey, ""),
		},
	)
	if err != nil {
		panic(err)
	}
	return session
}

func UploadSmallObject(file *os.File, dirName string) {

	fileinfo, err := file.Stat()

	if err != nil {
		log.Printf(color.RED, err)
	}

	session := createAWSSession()
	bucketName := os.Getenv("S3_BUCKET_NAME")
	printMessage := fmt.Sprintf("File <%s> started to upload... ", fileinfo.Name())
	log.Printf(color.LIGHTBLUE, printMessage)

	uploader := s3manager.NewUploader(session)
	_, err = uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(dirName + fileinfo.Name()),
		Body:   file,
	})
	if err != nil {
		log.Printf(color.RED, err)
		return
	}
	printMessage = fmt.Sprintf(fileinfo.Name() + " uploaded successfully")
	log.Printf(color.GREEN, printMessage)

}
