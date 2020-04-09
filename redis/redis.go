package redis

import (
	"github.com/firmeve/firmeve/kernel/contract"
	redis2 "github.com/go-redis/redis"
)

type (
	Redis struct {
		clusters    map[string]*redis2.ClusterClient
		connections map[string]*redis2.Client
		config      contract.Configuration
	}
)

func New(config contract.Configuration) *Redis {
	return &Redis{
		clusters:    make(map[string]*redis2.ClusterClient, 0),
		config:      config,
		connections: make(map[string]*redis2.Client, 0),
	}
}

func (r *Redis) Cluster(cluster string) *redis2.ClusterClient {
	if v, ok := r.clusters[cluster]; ok {
		return v
	}

	client := redis2.NewClusterClient(&redis2.ClusterOptions{
		Addrs: r.config.GetStringSlice(`clusters.` + cluster + `.hosts`),
	})

	r.clusters[cluster] = client

	return client
}

func (r *Redis) Connection(connection string) *redis2.Client {
	if v, ok := r.connections[connection]; ok {
		return v
	}

	connectionKey := `connections.` + connection + `.`

	client := redis2.NewClient(&redis2.Options{
		Addr: r.config.GetString(connectionKey + `host`),
		DB:   r.config.GetInt(connectionKey + `db`),
	})

	r.connections[connection] = client

	return client
}
