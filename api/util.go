package api

import (
	"github.com/pborman/uuid"
)

func randomID() string {
	return uuid.NewRandom().String()
}
