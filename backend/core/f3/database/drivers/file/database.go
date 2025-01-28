package file

import (
	"context"
	"f3/database/drivers"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FilesStorage struct {
	defaultBucket string
	apiAddress    string
	accessKey     string
	secretKey     string
	region        string

	client *minio.Client

	buckets                     map[string]string
	staticFileStorageRepository *StaticFileStorageRepository
}

type AddBucketFn func(*FilesStorage)

func New(apiAddress, accessKey, secretKey, region string, buckets ...AddBucketFn) drivers.FilesStorage {
	fileStorage := &FilesStorage{
		apiAddress: apiAddress,
		accessKey:  accessKey,
		secretKey:  secretKey,
		region:     region,
	}

	for _, addBucket := range buckets {
		addBucket(fileStorage)
	}

	return fileStorage
}

func (s *FilesStorage) Connect() error {
	fileStorageClient, err := minio.New(s.apiAddress, &minio.Options{
		Creds:  credentials.NewStaticV4(s.accessKey, s.secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return err
	}

	s.client = fileStorageClient

	if err = s.createBucket(context.TODO()); err != nil {
		return err
	}

	return nil
}

func (s *FilesStorage) Close(_ context.Context) error {
	return nil
}

func (s *FilesStorage) FilesRepository() drivers.FilesRepository {
	if s.staticFileStorageRepository != nil {
		return s.staticFileStorageRepository
	}

	s.staticFileStorageRepository = &StaticFileStorageRepository{
		defaultBucket: s.defaultBucket,
		buckets:       s.buckets,
		client:        s.client,
	}

	return s.staticFileStorageRepository
}

func InitBuckets(defaultBucket string, buckets map[string]string) AddBucketFn {
	return func(s *FilesStorage) {
		s.buckets = buckets

		s.defaultBucket = defaultBucket
	}
}

func (s *FilesStorage) createBucket(ctx context.Context) error {
	for _, bucket := range s.buckets {
		exists, err := s.client.BucketExists(ctx, bucket)
		if err != nil {
			return err
		}

		if exists {
			continue
		}

		if err = s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{
			Region: s.region,
		}); err != nil {
			return err
		}
	}

	return nil
}
