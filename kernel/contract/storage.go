package contract

import (
	"context"
	"github.com/firmeve/firmeve/support/hash"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

type (
	Storage interface {
		Url(key string, params url.Values) string

		FormUpload(ctx context.Context, metadata *StorageMetadata) (*StorageInfo, error)

		PartUpload(ctx context.Context, metadata *StorageMetadata, notify Notify) (*StorageInfo, error)

		Token() string

		Callback(req *http.Request) (*StorageInfo, error)
	}

	StorageInfo struct {
		Driver       string
		Hash         string
		Name         string
		OriginalName string
		Mime         string
		Extension    string
		Size         int64
		Path         string
		Extra        map[string]interface{}
	}

	// 得到的上传元数据
	StorageMetadata struct {
		Path string
		Name string
		Mime string
		Data []byte
		Size int64
		File io.Reader
	}

	NotifyInfo struct {
		Block int
		Size  int
		Extra interface{}
	}

	Notify func(blockCount int, info *NotifyInfo, err error)
)

func (a *StorageMetadata) GenerateNewName() string {
	return hash.Sha1String(a.Name+strconv.Itoa(int(time.Now().Unix()))) + path.Ext(a.Name)
}
