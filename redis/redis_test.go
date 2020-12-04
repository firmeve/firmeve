package redis

import (
	"github.com/firmeve/firmeve/kernel/contract"
	testing2 "github.com/firmeve/firmeve/testing"
	redis2 "github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	redis *redis2.Client
	app   contract.Application
)

func TestMain(m *testing.M) {
	// setup
	app = testing2.ApplicationDefault(new(Provider))
	redis = app.Resolve(`redis.client`).(*redis2.Client)

	m.Run()

	//teardown
}

func TestRedis_Connection(t *testing.T) {
	assert.IsType(t, &redis2.Client{}, app.Resolve(`redis.client`))
}

func TestRedis_Cluster(t *testing.T) {
	assert.IsType(t, &redis2.ClusterClient{}, app.Resolve(`redis.cluster`))
}
