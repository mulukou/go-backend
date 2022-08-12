package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type MinioHandler struct {
	minioInstance *minio.Client
}

var minioLock = &sync.Mutex{}

var minioInstance *MinioHandler

var minioUser = os.Getenv("MINIO_USER")
var minioPass = os.Getenv("MINIO_PASSWORD")

func getMinio() *MinioHandler {
	if minioInstance == nil {
		minioLock.Lock()
		defer minioLock.Unlock()
		if minioInstance == nil {
			fmt.Println("Creating connection to Minio.")
			ctx := context.Background()
			// MinIO Connection
			minioClient, err := minio.New(minioURL, &minio.Options{
				Creds: credentials.NewStaticV4(minioUser, minioPass, ""),
			})
			if err != nil {
				log.Fatalln(err)
			}
			err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
			if err != nil {
				exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
				if errBucketExists == nil && exists {
					log.Printf("We already own %s\n", bucketName)
				} else {
					log.Println(err)
				}
			} else {
				minioClient.SetBucketPolicy(ctx, bucketName, `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket","s3:ListBucketMultipartUploads"],"Resource":["arn:aws:s3:::react-shop"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject","s3:ListMultipartUploadParts","s3:PutObject","s3:AbortMultipartUpload","s3:DeleteObject"],"Resource":["arn:aws:s3:::react-shop/*"]}]}`)
				log.Printf("Successfully created %s\n", bucketName)
			}
			if err != nil {
				log.Println(err)
			}
			minioInstance = &MinioHandler{minioClient}
		} else {
			fmt.Println("Minio instance already created.")
		}
	} else {
		fmt.Println("Minio instance already created.")
	}

	return minioInstance
}
