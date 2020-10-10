package redis

import (
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	redis *Redis
	app   contract.Application
)

func TestMain(m *testing.M) {
	// setup
	app = testing2.ApplicationDefault(new(Provider))
	redis = app.Resolve(`redis.client`).(*Redis)

	m.Run()

	//teardown
}

func TestRedis_Connection(t *testing.T) {
	client := redis.Client(`default`)
	assert.IsType(t, &redis2.Client{}, client)
}

func TestRedis_Cluster(t *testing.T) {
	client := redis.Cluster(`default`)
	assert.IsType(t, &redis2.ClusterClient{}, client)
}
