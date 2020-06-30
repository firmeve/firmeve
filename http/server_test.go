package http

import (
	"github.com/stretchr/testify/assert"
	"time"

	//testing2 "github.com/firmeve/firmeve/testing"
	"github.com/firmeve/firmeve/testing/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestNewServer(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	router := mock.NewMockHttpRouter(mockCtrl)
	server1 := NewServer(router, map[string]interface{}{
		`host`: `0.0.0.0:11211`,
	})

	go func() {
		err := server1.Start()
		assert.Nil(t, err)
	}()
	time.Sleep(time.Second * 2)
	err := server1.Stop()
	assert.Nil(t, err)
}

func TestNewServerHttps(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	router := mock.NewMockHttpRouter(mockCtrl)
	server1 := NewServer(router, map[string]interface{}{
		`host`:      `0.0.0.0:11222`,
		`cert-file`: `../testdata/ssl/server.crt`,
		`key-file`:  `../testdata/ssl/server.key`,
	})

	go func() {
		err := server1.Start()
		assert.Nil(t, err)
	}()
	time.Sleep(time.Second * 2)
	err := server1.Stop()
	assert.Nil(t, err)
}

func TestNewServerHttp2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	router := mock.NewMockHttpRouter(mockCtrl)
	server1 := NewServer(router, map[string]interface{}{
		`host`:      `0.0.0.0:11223`,
		`cert-file`: `../testdata/ssl/server.crt`,
		`key-file`:  `../testdata/ssl/server.key`,
		`http2`:     true,
	})

	go func() {
		err := server1.Start()
		assert.Nil(t, err)
	}()
	time.Sleep(time.Second * 2)
	err := server1.Stop()
	assert.Nil(t, err)
}
