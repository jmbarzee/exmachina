package system

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"sync/atomic"
	"syscall"
	"time"
)

type (
	System struct {
		// halt is used to kill all goroutines in a call to Legionnaire.panic()
		Halt context.CancelFunc
		// log is where normal & debugging messages are dumped to
		Log *log.Logger
	}
)

func NewSystem(ctx context.Context, l *log.Logger) (context.Context, System) {
	s := System{
		Log: l,
	}
	go s.HandleSignals(ctx)
	ctx, cancel := context.WithCancel(ctx)
	s.Halt = cancel

	return ctx, s
}

var hackySignalDelay = int32(0)

func (s *System) HandleSignals(ctx context.Context) {
	// setup signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		s.Logf("System: received signal\n")
	case <-ctx.Done():
		s.Logf("System: context cancled\n")
	}
	s.Halt()
	time.Sleep(time.Duration(atomic.AddInt32(&hackySignalDelay, 2)) * time.Second)
	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	fmt.Printf("######################################\n\n\n")
}

func (s *System) Panic(err error) {
	s.Logf("System Panic!!! %v", err)
	s.Halt()
	//panic(err)
}

func (s *System) Logf(fmt string, args ...interface{}) {
	s.Log.Printf(fmt, args...)
}
