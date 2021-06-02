package KeyValueKernels

import "github.com/mediocregopher/radix/v3"

type RedisKeyValueKernel struct {
	pool     *radix.Pool
	ip       string
	poolSize int
}
