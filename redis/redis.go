package redis

import (
	redis2 "github.com/go-redis/redis/v8"
)

type (
	Redis struct {
		clusters map[string]*redis2.ClusterClient
		clients  map[string]*redis2.Client
		config   *Configuration
	}

	Configuration struct {
		Clients map[string]struct {
			Addr     string `json:"addr" yaml:"addr"`
			Password string `json:"password" yaml:"password"`
			DB       int    `json:"db" yaml:"db"`
		} `json:"clients" yaml:"clients"`

		Clusters map[string]struct {
			Addrs []string `json:"addrs" yaml:"addrs"`
		} `json:"clusters" yaml:"clusters"`
	}
)

func New(config *Configuration) *Redis {
	return &Redis{
		config:   config,
		clusters: make(map[string]*redis2.ClusterClient, 1),
		clients:  make(map[string]*redis2.Client, 1),
	}
}

func (r *Redis) Cluster(cluster string) *redis2.ClusterClient {
	if v, ok := r.clusters[cluster]; ok {
		return v
	}

	client := redis2.NewClusterClient(&redis2.ClusterOptions{
		Addrs: r.config.Clusters[cluster].Addrs,
	})

	r.clusters[cluster] = client

	return client
}

func (r *Redis) Client(connection string) *redis2.Client {
	if v, ok := r.clients[connection]; ok {
		return v
	}

	clientConfig, ok := r.config.Clients[connection]
	if !ok {
		panic(`the connection client not found`)
	}

	client := redis2.NewClient(&redis2.Options{
		Addr:     clientConfig.Addr,
		DB:       clientConfig.DB,
		Password: clientConfig.Password,
	})

	r.clients[connection] = client

	return client
}
