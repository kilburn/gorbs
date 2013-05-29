package log_test

import (
	"testing"

	"github.com/kilburn/gobs/log"
)

func TestWarn(t *testing.T) {
	log.GetBackend("std").SetLevel(log.WARN)
	log.Warn("It works!")
	log.Info("It doesn't work! :(")
}
