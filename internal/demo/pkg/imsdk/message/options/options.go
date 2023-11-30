package options

import "imsdk/internal/demo/pkg/imsdk"

type Options struct {
	Model       imsdk.ModelType
	Credentials *imsdk.Credentials
}

type Option func(opt *Options)
