package storage

//import (
//	"context"
//	"github.com/firmeve/firmeve/kernel/contract"
//	testing2 "github.com/firmeve/firmeve/testing"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestProvider_Register(t *testing.T) {
//	//container := container2.New()
//	//container
//
//	testing2.ApplicationDefault(new(Provider))
//
//	provider := testing2.DefaultContainer.Make().(*Provider)
//	provider.Register()
//	assert.NotNil(t, provider.Container.Get(`storage`))
//	assert.NotNil(t, provider.Container.Get(`storage`).(*Storage).Disk(Default))
//}

//func Test_Disk_Upload(t *testing.T) {
//	provider := testing2.DefaultContainer.Make(new(Provider)).(*Provider)
//	provider.Register()
//	assert.NotNil(t, provider.Container.Get(`storage`))
//	assert.NotNil(t, provider.Container.Get(`storage`).(*Storage).Disk(Default))
//
//	s := provider.Container.Get(`storage`).(*Storage)
//	_, err := s.Disk(Default).FormUpload(context.Background(), &contract.StorageMetadata{
//		Path: "../../test/testdata/1.jpg",
//		Name: "1.jpg",
//		Mime: "image/jpeg",
//		Data: nil,
//		Size: 0,
//	})
//	assert.Nil(t, err)
//}
