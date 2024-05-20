package uuid

import (
	"github.com/google/uuid"
)

type Generator interface {
	NewString() string
}

type Implementor struct{}

func (i Implementor) NewString() string {
	return uuid.NewString()
}
