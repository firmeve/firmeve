package http

import (
	"context"
	"github.com/stretchr/testify/assert"
	"sync"
	"time"

	//testing2 "github.com/firmeve/firmeve/testing"
	"github.com/firmeve/firmeve/testing/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func init() {
	//runtime.GOMAXPROCS(1)
}

var lock sync.Mutex

func TestNewServer(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	router := mock.NewMockHttpRouter(mockCtrl)
	server1 := NewServer(router, map[string]interface{}{
		`host`: `0.0.0.0:11211`,
	})
	ctx := context.Background()
	go func(t2 *testing.T) {
		server1.Start(ctx)
	}(t)
	time.Sleep(time.Second * 1)
	err := server1.Stop(ctx)
	assert.Nil(t, err)
}

func TestNewServerHttps(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()
	router := mock.NewMockHttpRouter(mockCtrl)
	server1 := NewServer(router, map[string]interface{}{
		`host`:      `0.0.0.0:11222`,
		`cert-file`: `../testdata/ssl/server.crt`,
		`key-file`:  `../testdata/ssl/server.key`,
	})

	go func(t2 *testing.T) {
		server1.Start(ctx)
		//assert.Nil(t2, err)
	}(t)
	time.Sleep(time.Second * 2)
	err := server1.Stop(ctx)
	assert.Nil(t, err)
}

func TestNewServerHttp2(t *testing.T) {
	lock.Lock()
	defer lock.Unlock()
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	ctx := context.Background()
	router := mock.NewMockHttpRouter(mockCtrl)
	server1 := NewServer(router, map[string]interface{}{
		`host`:      `0.0.0.0:11223`,
		`cert-file`: `../testdata/ssl/server.crt`,
		`key-file`:  `../testdata/ssl/server.key`,
		`http2`:     true,
	})

	go func(t2 *testing.T) {
		server1.Start(ctx)
	}(t)
	time.Sleep(time.Second * 2)
	err := server1.Stop(ctx)
	assert.Nil(t, err)
}
