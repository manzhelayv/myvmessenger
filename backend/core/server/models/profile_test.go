package models

import (
	"github.com/stretchr/testify/assert"
	"gitlab.com/myvmessenger/client/f3-client/protobuf"
	"testing"
)

func TestFilesProto(t *testing.T) {
	type args struct {
		tdid    string
		profile *Profile
	}
	tests := map[string]struct {
		args args
		want *protobuf.Files
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные": {
			args: args{
				tdid: "1111",
				profile: &Profile{
					Image: "111.png",
				},
			},
			want: &protobuf.Files{
				Files: []*protobuf.File{
					{
						Tdid: "1111",
						Name: "111.png",
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FilesProto(tt.args.tdid, tt.args.profile), "FilesProto(%v, %v)", tt.args.tdid, tt.args.profile)
		})
	}
}

func TestFileProto(t *testing.T) {
	type args struct {
		f *File
	}
	tests := map[string]struct {
		args args
		want *protobuf.File
	}{
		"Пустые данные": {
			args: args{},
			want: nil,
		},
		"Заполненные данные": {
			args: args{
				f: &File{
					Image: "111.png",
				},
			},
			want: &protobuf.File{
				Name:    PROFILE_IMAGE,
				Content: []byte("111.png"),
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, tt.want, FileProto(tt.args.f), "FileProto(%v)", tt.args.f)
		})
	}
}
