package system

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

type (
	System struct {
		// halt is used to kill all goroutines in a call to Legionnaire.panic()
		Halt context.CancelFunc
		// log is where normal & debugging messages are dumped to
		Log printable
	}

	printable interface {
		Printf(fmt string, args ...interface{})
	}
)

func NewSystem(ctx context.Context, p printable) (context.Context, System) {
	s := System{
		Log: p,
	}
	go s.HandleSignals(ctx)
	ctx, cancel := context.WithCancel(ctx)
	s.Halt = cancel

	return ctx, s
}

func (s *System) HandleSignals(ctx context.Context) {
	// setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		s.Logf("received signal\n")
		s.Halt()
		// TODO end all connections
	case <-ctx.Done():
		// Do nothing, just let the routine die
	}
}

func (s *System) Panic(err error) {
	s.Logf("System Panic!!! %v", err)
	s.Halt()
	//panic(err)
}

func (s *System) Logf(fmt string, args ...interface{}) {
	s.Log.Printf(fmt, args...)
}
