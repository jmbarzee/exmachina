package domain

import (
	"errors"
)

type (
	Service struct {
		Port          int
		ServiceConfig ServiceConfig
	}

	ServiceConfig struct {
		Name     string
		Priority Priority
		Depends  []string
		Traits   []string
	}
)

type Priority string

const (
	Required   Priority = "required"
	Dependency Priority = "dependency"
)

func (s *Service) Status() error {
	return errors.New("Unimplemented!")
}

func (s ServiceConfig) String() string {
	msg := "(" + s.Name + ", " +
		string(s.Priority) + ", "

	msg += "["
	for _, dependency := range s.Depends {
		msg += dependency + ", "
	}
	msg += "]"

	msg += ", "

	msg += "["
	for _, trait := range s.Traits {
		msg += trait + ", "
	}
	msg += "]"

	msg += ")"
	return msg
}
