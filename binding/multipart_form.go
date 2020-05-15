package binding

import (
	md52 "crypto/md5"
	"encoding/hex"
	"github.com/firmeve/firmeve/kernel/contract"
	"github.com/firmeve/firmeve/support/path"
	reflect2 "github.com/firmeve/firmeve/support/reflect"
	strings2 "github.com/firmeve/firmeve/support/strings"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	path2 "path"
	filepath2 "path/filepath"
	"reflect"
	"time"
)

type (
	multipartForm struct {
	}

	MultipartFiles map[string][]*multipart.FileHeader
)

var (
	MultipartForm = multipartForm{}
)

func (multipartForm) Protocol(protocol contract.Protocol, v interface{}) error {
	// 普通数据解析
	if err := Form.Protocol(protocol, v); err != nil {
		return err
	}

	// 上传文件解析
	if reflect2.KindElemType(reflect.TypeOf(v)) == reflect.Struct {
		reflectValue := reflect.Indirect(reflect.ValueOf(v))
		var bindingFile contract.BindingFile
		reflect2.CallFieldType(reflect2.IndirectType(reflect.TypeOf(v)), func(i int, field reflect.StructField) interface{} {
			// 这里 &bindingFile 是指向contract.BindingFile 类型
			// 如果只是bindingFile 则是一个nil 编译器并不知道类型，所以要加&指针引用
			if field.Type.Implements(reflect.TypeOf(&bindingFile).Elem()) {
				if reflectValue.Field(i).CanSet() {
					reflectValue.Field(i).Set(reflect.ValueOf(
						protocol.(contract.HttpProtocol).Request().MultipartForm.File),
					)
				}
			}

			return nil
		})
	}

	return nil
}

func (m MultipartFiles) Files(key string) []*multipart.FileHeader {
	return m[key]
}

func (m MultipartFiles) File(key string) *multipart.FileHeader {
	return m[key][0]
}

func (m MultipartFiles) Save(key string, option *contract.UploadOption) ([]*contract.FileInfo, error) {
	var filesInfo []*contract.FileInfo
	files := m.Files(key)
	for i := range files {
		fileInfo, err := m.save(files[i], option)
		if err != nil {
			return nil, err
		}
		filesInfo = append(filesInfo, fileInfo)
	}

	return filesInfo, nil
}

func (m MultipartFiles) save(file *multipart.FileHeader, option *contract.UploadOption) (*contract.FileInfo, error) {
	originalFile, err1 := file.Open()
	if err1 != nil {
		return nil, err1
	}
	defer originalFile.Close()
	var (
		newName  string
		filepath string
		fullpath string
		//mimeLength int64
		fileMime string
	)
	newName = FileName(file, option)
	filepath = Filepath(option)
	// 自动创建目录
	if !path.Exists(filepath) {
		if err := os.MkdirAll(filepath, 0755); err != nil {
			return nil, err
		}
	}
	fullpath, err := filepath2.Abs(path2.Join(filepath, newName))
	if err != nil {
		return nil, err
	}
	newFile, err2 := os.Create(fullpath)
	if err2 != nil {
		return nil, err2
	}
	defer newFile.Close()
	_, err = io.Copy(newFile, originalFile)

	//mime判断
	fileMime, err = FileMime(file)
	if err != nil {
		return nil, err
	}

	return &contract.FileInfo{
		Size:         file.Size,
		OriginalName: file.Filename,
		FullPath:     fullpath,
		Path:         filepath,
		Name:         newName,
		Extension:    path2.Ext(newName)[1:],
		Mime:         fileMime,
		MD5:          FileHash(originalFile),
	}, err
}

func (m MultipartFiles) SaveAll(option *contract.UploadOption) (map[string][]*contract.FileInfo, error) {
	return nil, nil
}

func FileName(file *multipart.FileHeader, option *contract.UploadOption) string {
	if option.Rename {
		hash := md52.New()
		hash.Write([]byte(time.Now().String() + file.Filename))
		return hex.EncodeToString(hash.Sum(nil)) + path2.Ext(file.Filename)
	}

	return file.Filename
}

func FileHash(file multipart.File) string {
	md5 := md52.New()
	io.Copy(md5, file)
	return hex.EncodeToString(md5.Sum(nil))
}

func Filepath(option *contract.UploadOption) string {
	var (
		path      string
		directory string
	)
	if option.Path != `` {
		path = option.Path
	} else {
		path = "."
	}
	if option.Grading {
		directory = strings2.Join("/", time.Now().Format(`2006-01`), time.Now().Format(`02`))
	} else {
		directory = ""
	}
	return path2.Join(path, directory)
}

func FileMime(file *multipart.FileHeader) (string, error) {
	//mime判断 length最多取512
	mimeLength := file.Size
	if mimeLength > 512 {
		mimeLength = 512
	}
	mimeBytes := make([]byte, mimeLength)
	f, err := file.Open()
	if err != nil {
		return ``, err
	}
	defer f.Close()
	//f.Seek(0, 0)
	_, err = f.Read(mimeBytes)
	return http.DetectContentType(mimeBytes), err
}
