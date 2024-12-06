package service

import (
	"Fire/config"
	"Fire/pkg/util"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/policy"
	"mime/multipart"
	"net/url"
	"time"
)

type MinioService struct {
}

func (service *MinioService) Upload(file *multipart.FileHeader) (string, error) { // TODO: 使用异步或 Goroutine 来优化文件上传的性能（如果有大量上传操作）。
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	// Create a unique object name
	objectName := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)

	// Upload file to the MinIO bucket
	_, err = util.MinioClient.PutObject(context.Background(), config.Config.Minio.BucketName, objectName, src, file.Size, minio.PutObjectOptions{
		ContentType: file.Header.Get("Content-Type"),
	})
	if err != nil {
		return "", err
	}
	return objectName, nil
}

// Preview generates a presigned URL for accessing the file
func (service *MinioService) Preview(fileName string) (string, error) {
	reqParams := url.Values{}
	urlPath, err := util.MinioClient.PresignedGetObject(context.Background(), config.Config.Minio.BucketName, fileName, time.Hour*24, reqParams)
	if err != nil {
		return "", err
	}
	return urlPath.String(), nil
}

// EnsureBucket ensures that the bucket exists and is public
//func (service *MinioService) EnsureBucket() error {
//	exists, err := service.Client.BucketExists(context.Background(), service.Config.BucketName)
//	if err != nil {
//		return fmt.Errorf("error checking bucket existence: %v", err)
//	}
//	if !exists {
//		err := service.Client.MakeBucket(context.Background(), service.Config.BucketName, minio.MakeBucketOptions{})
//		if err != nil {
//			return fmt.Errorf("error creating bucket: %v", err)
//		}
//		err = service.Client.SetBucketPolicy(context.Background(), service.Config.BucketName, string(policy.BucketPolicyReadOnly))
//		if err != nil {
//			return fmt.Errorf("error setting bucket policy: %v", err)
//		}
//	}
//	return nil
//}

func (service *MinioService) EnsureBucket() error {
	err := util.MinioClient.MakeBucket(context.Background(), config.Config.Minio.BucketName, minio.MakeBucketOptions{})
	if err != nil {
		if exists, checkErr := util.MinioClient.BucketExists(context.Background(), config.Config.Minio.BucketName); checkErr == nil && exists {
			return nil // Bucket already exists
		}
		return fmt.Errorf("error ensuring bucket existence: %v", err)
	}

	err = util.MinioClient.SetBucketPolicy(context.Background(), config.Config.Minio.BucketName, string(policy.BucketPolicyReadOnly))
	if err != nil {
		return fmt.Errorf("error setting bucket policy: %v", err)
	}
	return nil
}
