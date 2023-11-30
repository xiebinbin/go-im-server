package user

import (
	"imsdk/internal/demo/pkg/imsdk"
	"sync"
)

type Options struct {
	Model       imsdk.ModelType
	Credentials *imsdk.Credentials
}

type Client struct {
	options *Options
}

var (
	once   sync.Once
	client *Client
)

func NewClient(options *Options) *Client {
	once.Do(func() {
		client = &Client{
			options: options,
		}
	})
	return client
}

