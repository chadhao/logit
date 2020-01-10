package config

import (
	"context"
	"strings"
	"time"

	"go.etcd.io/etcd/clientv3"
)

var KEY_PREFIX = "/logit"
var CLEINT_ENDPOINTS = []string{
	// "dev.logit.co.nz:2379",
	"localhost:2379",
}

type (
	Config interface {
		LoadConfig() error
		Get(string) (string, bool)
		LoadModuleConfig(string) map[string]string
	}
	config struct {
		client *clientv3.Client
		config map[string]string
	}
)

func (c *config) connect() error {
	var err error

	c.client, err = clientv3.New(clientv3.Config{
		Endpoints:   CLEINT_ENDPOINTS,
		DialTimeout: 5 * time.Second,
	})

	if err != nil {
		return err
	}

	return nil
}

func (c *config) LoadConfig() error {
	if err := c.connect(); err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := c.client.Get(ctx, KEY_PREFIX, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend))
	cancel()

	if err != nil {
		return err
	}

	for _, kv := range resp.Kvs {
		ks := strings.Split(string(kv.Key), "/")
		c.config[strings.Join(ks[2:], ".")] = string(kv.Value)
	}

	return c.client.Close()
}

func (c *config) Get(k string) (v string, ok bool) {
	v, ok = c.config[k]
	return
}

func (c *config) LoadModuleConfig(m string) map[string]string {
	moduleConfig := make(map[string]string)
	for k, v := range c.config {
		if strings.Index(k, m+".") == 0 {
			moduleConfig[k] = v
		}
	}
	return moduleConfig
}

func New() Config {
	return &config{
		client: nil,
		config: make(map[string]string),
	}
}
