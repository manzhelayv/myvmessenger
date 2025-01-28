package drivers

import (
	"context"
	"github.com/aws/aws-sdk-go/service/s3"
)

type ImagesRepository interface {
	UploadImage(ctx context.Context, imageName, folder string, data []byte) error
	DownloadImage(ctx context.Context, imageName, folder string) ([]byte, error)
	ImageList(ctx context.Context, prefix string, limit int64) ([]*s3.Object, error)
	DeleteImages(ctx context.Context, imageNames []string) error
}

type FilesRepository interface {
	DownloadFile(ctx context.Context, fileName, folder string) ([]byte, error)
	UploadFile(ctx context.Context, fileName, folder string, data []byte) error
}
