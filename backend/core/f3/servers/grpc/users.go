package grpcresources

import (
	"context"
	"f3/manager"
	fileProto "gitlab.com/myvmessenger/client/f3-client/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FileStorageResource struct {
	fileMan manager.FileService
	fileProto.UnimplementedFileClientServer
}

func NewFileStorageResource(fileMan manager.FileService) *FileStorageResource {
	return &FileStorageResource{fileMan: fileMan}
}

func (m *FileStorageResource) UploadFile(ctx context.Context, f *fileProto.File) (*fileProto.ID, error) {
	id, err := m.fileMan.UploadFile(ctx, f.Name, f.Content)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &fileProto.ID{
		Id: id,
	}, nil
}

func (m *FileStorageResource) LoadFiles(ctx context.Context, fileInfo *fileProto.Files) (*fileProto.Files, error) {
	files := make([]*fileProto.File, 0, len(fileInfo.Files))
	for _, f := range fileInfo.Files {
		file, err := m.fileMan.DownloadFile(ctx, f.Name)
		if err != nil {
			return nil, status.Errorf(codes.Internal, err.Error())
		}

		files = append(files, &fileProto.File{
			Tdid:    f.Tdid,
			Name:    f.Name,
			Content: file,
		})
	}

	return &fileProto.Files{
		Files: files,
	}, nil
}
