package kernel

import (
	errors2 "goshort/types/errors"
	"os"
	"os/signal"
	"syscall"
)

type SignalKernel struct {
	Kernel *Kernel
}

func (signalKernel *SignalKernel) Run() error {
	signals := make(chan os.Signal, 1)

	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	_ = signalKernel.Kernel.Logger.Send(&errors2.GenericLog{Name: "Kernel.Signal.StartWaiting", IsError: false})

	select {
	case systemSignal := <-signals:
		_ = signalKernel.Kernel.Logger.Send(&SignalGotLog{Signal: systemSignal.String()})
		break
	}

	return nil
}

func (signalKernel *SignalKernel) Signal(_ interface{}) error {
	return nil
}
