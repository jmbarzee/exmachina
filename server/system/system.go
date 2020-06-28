package system

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

type (
	System struct {
		// halt is used to kill all goroutines in a call to Domain.panic()
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

	// time.Sleep(time.Duration(1 * time.Second))

	pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
	fmt.Printf("######################################\n\n\n")
}

func (s *System) Panic(err error) {
	s.Logf("System Panic!!! %v", err)
	s.Halt()

	// kill -pgid (-pid)
	// ends all child processes. 90% certain (negative means groupID?)
	// TODO @jmbarzee consider killing services individually after saving command objects
	pgid, sysErr := syscall.Getpgid(syscall.Getpid())
	if sysErr != nil {
		panic(sysErr)
	}
	syscall.Kill(-pgid, syscall.SIGKILL)

	panic(err)
}

func (s *System) Logf(fmt string, args ...interface{}) {
	s.Log.Printf(fmt, args...)
}
