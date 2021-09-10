package disk

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/qiniu/api.v7/v7/auth/qbox"
	"github.com/qiniu/api.v7/v7/storage"
	"net/http"
	url2 "net/url"
	"os"
	"strconv"
	"time"
)

const QiNiuDisk = `qiniu`

type (
	QiNiu struct {
		config *QiNiuConfiguration
		bucket *storage.BucketManager
	}

	QiNiuConfiguration struct {
		AccessKey string `json:"access_key" yaml:"access_key" mapstructure:"access_key"`
		SecretKey string `json:"secret_key" yaml:"secret_key" mapstructure:"secret_key"`
		Bucket    string `json:"bucket" yaml:"bucket"`
		IsPrivate bool   `json:"is_private" yaml:"is_private" mapstructure:"is_private"`
		Region    string `json:"region" yaml:"region"` // z0:Huadong z1:Huabei z2:Huanan na0:NorthAmerica as0:Singapore
		Domain    string `json:"domain" yaml:"domain"`
		Expires   uint64 `json:"expires" yaml:"expires"` // 资源公开有效期，单位：s
		UseHttps  bool   `json:"use_https" yaml:"use_https" mapstructure:"use_https"`
		UseCdn    bool   `json:"use_cdn" yaml"use_cdn" mapstructure:"use_cdn"`
		NotifyUrl string `json:"notify_url" yaml:"notify_url" mapstructure:"notify_url"`
	}

	ReturnResult struct {
		Hash         string `json:"hash"`
		PersistentID string `json:"persistent_id"`
		Key          string `json:"key"`
		OriginalName string `json:"original_name"`
		Ext          string `json:"ext"`
		Bucket       string `json:"bucket"`
		ImageInfo    struct {
			Width  int `json:"width"`
			Height int `json:"height"`
		} `json:"image_info"`
		Exif     map[string]interface{} `json:"exif"`
		AvInfo   map[string]interface{} `json:"avinfo"`
		MimeType string                 `json:"mime_type"`
		Size     int                    `json:"size"`
	}
)

// 一个新的bucket是一个新的Server
func NewQiNiu(config *QiNiuConfiguration) contract.Storage {
	return &QiNiu{
		config: config,
		bucket: createBucket(config),
	}
}

// 创建所有桶
func createBucket(config *QiNiuConfiguration) *storage.BucketManager {
	region, _ := storage.GetRegionByID(storage.RegionID(config.Region))
	return storage.NewBucketManager(qbox.NewMac(config.AccessKey, config.SecretKey), &storage.Config{
		Zone:          &region,         // 空间对应的机房
		UseHTTPS:      config.UseHttps, // 是否使用https域名
		UseCdnDomains: config.UseCdn,   // 上传是否使用CDN上传加速
	})
}

// 返回临时的一次性令牌
func (q QiNiu) Token() string {
	return q.uploadToken()
}

// 获取资源访问URL
func (q QiNiu) Url(key string, params url2.Values) string {
	var url string
	prefix := `?`
	if params != nil{
		w ,ok := params[`w`]
		h ,ok2 := params[`h`]
		if ok && ok2 {
			prefix += fmt.Sprintf(`imageView2/1/w/%s/h/%s/q/%s`,w[0],h[0],`75`) + `&`
		}

		key = key + prefix + params.Encode()
	}

	if q.config.IsPrivate {
		deadline := time.Now().Add(time.Second * time.Duration(q.config.Expires)).Unix()
		url = storage.MakePrivateURL(q.bucket.Mac, q.config.Domain, key, deadline)
	} else {
		url = storage.MakePublicURL(q.config.Domain, key)
	}

	return url
}

// 简单表单上传
func (q QiNiu) FormUpload(ctx context.Context, metadata *contract.StorageMetadata) (*contract.StorageInfo, error) {
	var (
		localFile = metadata.Path
		byteData = metadata.Data
		fileIO = metadata.File
		newName   = metadata.GenerateNewName()
		path      = q.config.Bucket + `/` + newName
		result = new(ReturnResult)
		extra = &storage.PutExtra{
			Params: map[string]string{
				`x:original_name`: metadata.Name,
			},
		}
		uploader = storage.NewFormUploader(q.bucket.Cfg)
		err error
	)

	// 构建表单上传的对象和获取的处理结果
	if localFile != ``{
		err = uploader.PutFile(ctx, result, q.uploadToken(), path, localFile, extra)
	} else if fileIO != nil{
		err = uploader.Put(ctx, result, q.uploadToken(), path, fileIO, metadata.Size, extra)
	} else if byteData != nil {
		err = uploader.Put(ctx, result, q.uploadToken(), path, bytes.NewReader(byteData), metadata.Size,extra)
	} else {
		err = errors.New(`upload source error`)
	}

	if err != nil {
		return nil, err
	}

	return q.storageInfo(newName, path, metadata, result), nil
}

// 分块上传
func (q QiNiu) PartUpload(ctx context.Context, metadata *contract.StorageMetadata, notify contract.Notify) (*contract.StorageInfo, error) {
	var (
		localFile = metadata.Path
		newName   = metadata.GenerateNewName()
		path      = q.config.Bucket + `/` + newName
	)

	// 获取当前文件总块数
	localFileInfo, err := os.Stat(localFile)
	if err != nil {
		return nil, err
	}
	blockCount := storage.BlockCount(localFileInfo.Size())

	// 分块上传回调申明
	var (
		putNotify    func(blkIdx int, blkSize int, ret *storage.BlkputRet)
		putNotifyErr func(blkIdx int, blkSize int, err error)
	)
	if notify != nil {
		putNotify = func(blkIdx int, blkSize int, ret *storage.BlkputRet) {
			notify(blockCount, &contract.NotifyInfo{
				Block: blkIdx,
				Size:  blkSize,
				Extra: ret,
			}, nil)
		}
		putNotifyErr = func(blkIdx int, blkSize int, err error) {
			notify(blockCount, &contract.NotifyInfo{
				Block: blkIdx,
				Size:  blkSize,
				Extra: nil,
			}, err)
		}
	}

	result := new(ReturnResult)
	resumeUploader := storage.NewResumeUploader(q.bucket.Cfg)
	err = resumeUploader.PutFile(ctx, &result, q.uploadToken(), path, localFile, &storage.RputExtra{
		Params: map[string]string{
			`x:original_name`: metadata.Name,
		},
		Notify:    putNotify,
		NotifyErr: putNotifyErr,
	})
	if err != nil {
		return nil, err
	}

	return q.storageInfo(newName, path, metadata, result), err
}

// 客户端上传成功后的异步回调
// TODO: 此功能需要测试
func (q QiNiu) Callback(req *http.Request) (*contract.StorageInfo, error) {
	status, err := qbox.VerifyCallback(q.bucket.Mac, req)
	if !status || err != nil {
		return nil, err
	}

	req.ParseForm()
	form := req.Form

	size, err := strconv.Atoi(form.Get(`size`))
	if err != nil {
		return nil, err
	}
	return &contract.StorageInfo{
		Driver:       QiNiuDisk,
		Hash:         form.Get(`hash`),
		Name:         form.Get(`fname`),
		OriginalName: form.Get(`original_name`),
		Mime:         form.Get(`mime`),
		Extension:    form.Get(`ext`),
		Size:         int64(size),
		Path:         form.Get(`key`),//q.config.Bucket + `/` + form.Get(`key`),
		Extra: map[string]interface{}{
			`image`:         form.Get(`image_info`),
			`exif`:          form.Get(`exif`),
			`audio`:         form.Get(`audio`),
			`video`:         form.Get(`video`),
			`avinfo`:        form.Get(`avinfo`),
			`persistent_id`: form.Get(`persistent_id`),
		},
	}, nil
}


// 桶对象
func (q QiNiu) Bucket() *storage.BucketManager {
	return q.bucket
}

// 获取上传token
func (q QiNiu) uploadToken() string {
	// 这里是JSON格式有些是int类型是不用加引号("")的
	var body = `{"key":"$(key)","hash":"$(hash)","ext":"$(ext)","persistent_id":"$(persistentId)","exif":$(exif),"avinfo":$(avinfo),"fsize":$(fsize),"image_info":$(imageInfo),"bucket":"$(bucket)","fname":"$(fname)","mime_type":"$(mimeType)","original_name":"$(x:original_name)"}`
	policy := &storage.PutPolicy{
		Scope:        q.config.Bucket,
		CallbackURL:  q.config.NotifyUrl,
		CallbackBody: body,
		ReturnBody:   body,
		Expires:      q.config.Expires,
	}
	return policy.UploadToken(q.bucket.Mac)
}

// 构建存储信息
func (q QiNiu) storageInfo(name, path string, metadata *contract.StorageMetadata, result *ReturnResult) *contract.StorageInfo {
	return &contract.StorageInfo{
		Driver:       QiNiuDisk,
		Name:         name,
		OriginalName: metadata.Name,
		Mime:         result.MimeType,
		Extension:    result.Ext,
		Size:         int64(result.Size),
		Path:         path,
		Hash:         result.Hash,
		Extra: map[string]interface{}{
			`image`:         result.ImageInfo,
			`exif`:          result.Exif,
			`audio`:         result.AvInfo[`audio`],
			`video`:         result.AvInfo[`video`],
			`avinfo`:        result.AvInfo,
			`persistent_id`: result.PersistentID,
		},
	}
}
