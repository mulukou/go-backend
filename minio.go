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

func getMinio() *MinioHandler {
	if minioInstance == nil {
		minioLock.Lock()
		defer minioLock.Unlock()
		if minioInstance == nil {
			fmt.Println("Creating connection to Minio.")
			ctx := context.Background()
			// MinIO Connection
			minioClient, err := minio.New(os.Getenv("MINIO_URL"), &minio.Options{
				Creds: credentials.NewStaticV4(os.Getenv("MINIO_USER"), os.Getenv("MINIO_PASSWORD"), ""),
			})
			if err != nil {
				log.Fatalln(err)
			}
			err = minioClient.MakeBucket(ctx, os.Getenv("MINIO_BUCKET"), minio.MakeBucketOptions{})
			if err != nil {
				exists, errBucketExists := minioClient.BucketExists(ctx, os.Getenv("MINIO_BUCKET"))
				if errBucketExists == nil && exists {
					log.Printf("We already own %s\n", os.Getenv("MINIO_BUCKET"))
				} else {
					log.Println(err)
				}
			} else {
				minioClient.SetBucketPolicy(ctx, os.Getenv("MINIO_BUCKET"), `{"Version":"2012-10-17","Statement":[{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetBucketLocation","s3:ListBucket","s3:ListBucketMultipartUploads"],"Resource":["arn:aws:s3:::react-shop"]},{"Effect":"Allow","Principal":{"AWS":["*"]},"Action":["s3:GetObject","s3:ListMultipartUploadParts","s3:PutObject","s3:AbortMultipartUpload","s3:DeleteObject"],"Resource":["arn:aws:s3:::react-shop/*"]}]}`)
				log.Printf("Successfully created %s\n", os.Getenv("MINIO_BUCKET"))
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
