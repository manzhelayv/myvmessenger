package file

import (
	"bytes"
	"context"
	"fmt"
	minio "github.com/minio/minio-go/v7"
	"os"
	"strings"
)

type StaticFileStorageRepository struct {
	defaultBucket string
	buckets       map[string]string
	client        *minio.Client
}

func (s *StaticFileStorageRepository) DownloadFile(ctx context.Context, fileName, folder string) ([]byte, error) {
	key := fileName

	if folder != "" {
		key = fmt.Sprintf("%s/%s", folder, fileName)
	}

	response, err := s.client.GetObject(ctx, s.defineBucket(fileName), key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}

	defer response.Close()

	buf := new(bytes.Buffer)

	_, err = buf.ReadFrom(response)

	errResp, ok := err.(minio.ErrorResponse)
	if !ok {
		return buf.Bytes(), nil
	}

	if errResp.Code == "NoSuchKey" {
		return nil, os.ErrNotExist
	}

	return nil, err
}

func (s *StaticFileStorageRepository) UploadFile(ctx context.Context, fileName, folder string, data []byte) error {
	key := fileName

	if folder != "" {
		key = fmt.Sprintf("%s/%s", folder, fileName)
	}

	_, err := s.client.PutObject(ctx, s.defineBucket(fileName), key, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return err
}

func (s *StaticFileStorageRepository) defineBucket(imageName string) string {
	fileParts := strings.Split(imageName, "_")
	if len(fileParts) == 0 {
		return s.defaultBucket
	}

	if _, ok := s.buckets[fileParts[0]]; ok {
		return fileParts[0]
	}

	return s.defaultBucket
}
