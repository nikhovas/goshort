package AGSKeyValueDrivers

import (
	"github.com/mediocregopher/radix/v3"
	kernelErrors "goshort/src/types/errors"
)

type Redis struct {
	pool *radix.Pool
}

func (r *Redis) Get(key string) ([]byte, error) {
	var data []byte
	err := r.pool.Do(radix.Cmd(&data, "GET", "url "+key))
	if err == nil && len(data) == 0 {
		err = kernelErrors.NotFoundError
	}
	return data, err
}

func (r *Redis) SetIfNotExists(key string, data []byte) error {
	var data2 []byte
	return r.pool.Do(radix.Cmd(&data2, "SET", key, string(data), "NX"))
}

func (r *Redis) GetNewUniqueNumber() (int, error) {
	var counter int
	err := r.pool.Do(radix.Cmd(&counter, "INCR", "counter"))
	return counter, err
}
