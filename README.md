# multipart-multifile-fileuploader

## Description

A simple file uploader for s3 . You can upload files ranges from small files to large file in GB's. For larger files it uses Multipart Upload and uploades files chunck by chunk (chunck size = 5mb).
Can also upload folders and uses goroutines to upload them in parallel . (Max 5 at a time)
Cannot upload folders within folder and ignores them (logs error)

---

## How to run

First get all dependencies , therefore run command in root of repo
`go get`

**Edit config.env and enter you credentials like aws_region and aws_bucket_name**

_Also for aws_session credentials like token,key and access_id , see [aws docs](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html)_

Now to run this project use -
`go run github.com/divyambhutani/multipart-multifile-uploader`
and enter the files you want to upload


<img src="img/Screenshot from 2022-03-28 11-52-12.png"  alt="News Screen" width="600"  />
<img src="img/Screenshot from 2022-03-28 11-54-10.png"  alt="News Screen" width="600"  />
<img src="img/Screenshot from 2022-03-28 11-56-15.png"  alt="News Screen" width="400"  />
