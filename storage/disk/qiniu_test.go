package disk

import (
	"context"
	"fmt"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/stretchr/testify/assert"
	"testing"
)

var baseSetting = struct {
	AccessKey string
	SecretKey string
	Bucket    string
}{
	AccessKey: `BCAWdfbevaB7UutUJC152_MTNvUjrhN330G5oWkK`,
	SecretKey: `4obFiUCsIx5P3xMaWrv7rjRAwyL2s8XMNaU0qtsW`,
	Bucket:    `firmeve-blog`,
}

func TestQiNiu_FormUpload2(t *testing.T) {
	qiniu := NewQiNiu(&QiNiuConfiguration{
		AccessKey: baseSetting.AccessKey,
		SecretKey: baseSetting.SecretKey,
		Bucket:    baseSetting.Bucket,
		Domain:    `http://attachment.crcms.cn`,
		Region:    `z2`,
		IsPrivate: true,
		Expires:   100,
		UseHttps:  false,
		UseCdn:    false,
	})
	v, err := qiniu.FormUpload(context.Background(), &contract.StorageMetadata{
		Path: "../../../test/testdata/1.jpg",
		Name: "1.jpg",
		Mime: "image/jpeg",
		Data: nil,
		Size: 0,
	})
	assert.Nil(t, err)
	assert.NotNil(t, v)
	fmt.Printf("%#v\n", v)
}

func TestQiNiu_Url(t *testing.T) {
	qiniu := NewQiNiu(&QiNiuConfiguration{
		AccessKey: baseSetting.AccessKey,
		SecretKey: baseSetting.SecretKey,
		Bucket:    baseSetting.Bucket,
		Domain:    `http://attachment.crcms.cn`,
		Region:    `z2`,
		IsPrivate: true,
		Expires:   100,
		UseHttps:  false,
		UseCdn:    false,
	})

	items, _, _, _, err := qiniu.(*QiNiu).bucket.ListFiles(`firmeve-blog`, ``, ``, ``, 1)
	assert.Nil(t, err)
	fmt.Println(qiniu.Url(items[0].Key, map[string][]string{
		`w`: {`200`},
		`h`: {`200`},
	}))
}

func TestQiNiu_PartUpload(t *testing.T) {
	qiniu := NewQiNiu(&QiNiuConfiguration{
		AccessKey: baseSetting.AccessKey,
		SecretKey: baseSetting.SecretKey,
		Bucket:    baseSetting.Bucket,
		Domain:    `http://attachment.crcms.cn`,
		Region:    `z2`,
		IsPrivate: true,
		Expires:   100,
		UseHttps:  false,
		UseCdn:    false,
	})
	var blockSum = 1
	result, err := qiniu.PartUpload(context.Background(), &contract.StorageMetadata{
		Path: "../../../test/testdata/1.mp4",
		Name: "1.mp4",
		Mime: "image/jpeg",
		Data: nil,
		Size: 0,
	}, func(blockCount int, info *contract.NotifyInfo, err error) {
		if err != nil {
			panic(err)
		}
		fmt.Printf("当前上传%.2f\n", (float64(blockSum)/float64(blockCount))*100)
		blockSum += 1
	})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	fmt.Printf("%#v\n", result)
}
