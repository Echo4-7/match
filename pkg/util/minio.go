package util

import (
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var MinioClient *minio.Client

func InitMinio(endpoint, accessKey, secretKey string) (err error) {
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false, // 如果使用 HTTPS，则设为 true
	})
	fmt.Printf("%v\n", err)
	if err != nil {
		log.Fatalf("Failed to initialize MinIO client: %v", err)
		return err
	}
	log.Println("MinIO client initialized successfully")
	return nil
}
