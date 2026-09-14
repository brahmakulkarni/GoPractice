package main

import (
	"fmt"
)

type Storage interface {
	Upload(file string) error
}

type S3Client struct{}

func (s *S3Client) PutObject(bucket string, key string) error {
	fmt.Printf("S3: uploading %s to bucket %s\n", key, bucket)
	return nil
}

type S3StorageAdapter struct {
	client *S3Client
	bucket string
}

func NewS3StorageAdapter(client *S3Client, bucket string) *S3StorageAdapter {
	return &S3StorageAdapter{
		client: client,
		bucket: bucket,
	}
}

func (s *S3StorageAdapter) Upload(file string) error {
	return s.client.PutObject(s.bucket, file)
}

type LocalDiskClient struct{}

func (l *LocalDiskClient) WriteFile(path string) error {
	fmt.Println("Writing file to:", path)
	return nil
}

type LocalDiskClientAdapter struct {
	client *LocalDiskClient
	path string
}

func NewLocalDiskClientAdapter(client *LocalDiskClient, path string) *LocalDiskClientAdapter {
	return &LocalDiskClientAdapter{
		client: client,
		path: path,
	}
}

func (l *LocalDiskClientAdapter) Upload(file string) error {
	return l.client.WriteFile(l.path+"/"+file)
}

type MediaUploader struct {
	storage Storage
}

func NewMediaUploader(storage Storage) *MediaUploader {
	return &MediaUploader{
		storage: storage,
	}
}

func (m *MediaUploader) Upload(file string) error {
	return m.storage.Upload(file)
}

func main() {
	localDiskClient := &LocalDiskClient{}

	storage := NewLocalDiskClientAdapter(localDiskClient, "my-uploads")

	uploader := NewMediaUploader(storage)

	if err := uploader.Upload("movie.mp4"); err != nil {
		fmt.Println(err)
		return
	}
}