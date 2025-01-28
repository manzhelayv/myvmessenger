package manager

import (
	"context"
	"crypto/md5"
	"f3/database/drivers"
	"fmt"
	"strings"
)

type FileService interface {
	DownloadFile(ctx context.Context, filePath string) ([]byte, error)
	UploadFile(ctx context.Context, filePath string, content []byte) (string, error)
}

type FileManager struct {
	repo drivers.FilesRepository
}

func NewFileStorageManager(repo drivers.FilesRepository) FileService {
	return &FileManager{repo: repo}
}

func (im *FileManager) UploadFile(ctx context.Context, filePath string, content []byte) (string, error) {
	folder, fileName := folderAndFileName(filePath)

	hash := md5.Sum(content)

	hashFile := fmt.Sprintf("%x", hash)

	fileName = fmt.Sprintf(`%s_%s`, fileName, hashFile)

	err := im.repo.UploadFile(ctx, fileName, folder, content)
	if err != nil {
		return "", err
	}

	return fileName, nil
}

func (im *FileManager) DownloadFile(ctx context.Context, filePath string) ([]byte, error) {
	folder, fileName := folderAndFileName(filePath)

	file, err := im.repo.DownloadFile(ctx, fileName, folder)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func folderAndFileName(path string) (string, string) {
	folder, fileName := "", path

	if fullPath := strings.Split(path, "/"); len(fullPath) > 1 {
		folder, fileName = fullPath[0], fullPath[1]
	}

	return folder, fileName
}
