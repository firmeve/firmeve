package contract

import "mime/multipart"

type (
	Binding interface {
		Protocol(protocol Protocol, v interface{}) error
	}

	BindingFile interface {
		Files(key string) []*multipart.FileHeader

		File(key string) *multipart.FileHeader

		Save(key string, option *UploadOption) ([]*FileInfo, error)

		SaveAll(option *UploadOption) (map[string][]*FileInfo, error)
	}

	FileInfo struct {
		Name         string
		OriginalName string
		Mime         string
		Extension    string
		Size         int64
		Path         string
		FullPath     string
		MD5          string
	}

	UploadOption struct {
		Path            string
		AllowMimes      []string
		AllowExtensions []string
		MaxSize         int64
		MinSize         int64
		Grading         bool //多层级
		Rename          bool
	}
)
