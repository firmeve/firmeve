package redis

import (
	testing2 "github.com/firmeve/firmeve/testing"
	redis2 "github.com/go-redis/redis"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	redis *Redis
)

func TestMain(m *testing.M) {
	// setup
	testing2.TestingApplication.Register(new(Provider), true)
	redis = testing2.TestingApplication.Resolve(`redis.client`).(*Redis)

	m.Run()

	//teardown
}

func TestRedis_Connection(t *testing.T) {
	client := redis.Connection(`default`)
	assert.IsType(t, &redis2.Client{}, client)
}

func TestRedis_Cluster(t *testing.T) {
	client := redis.Cluster(`default`)
	assert.IsType(t, &redis2.ClusterClient{}, client)
}
