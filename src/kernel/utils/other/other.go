package other

import (
	"goshort/src/types"
	errors2 "goshort/src/types/errors"
)

func InterfaceToLogWrapper(element interface{}) types.Log {
	log, ok := element.(types.Log)
	if !ok {
		err, ok := element.(error)
		if ok {
			return &errors2.SimpleErrorWrapper{Err: err}
		} else {
			return nil
		}
	} else {
		return log
	}
}
