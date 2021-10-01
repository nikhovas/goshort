package kernel

type SignalGotLog struct {
	Signal string
}

func (log *SignalGotLog) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"name":   "Kernel.Signal.Got",
		"type":   "log",
		"signal": log.Signal,
	}
}
