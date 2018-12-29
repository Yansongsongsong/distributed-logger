package test

import (
	"logsystem"
	"testing"
)

func TestAdd(t *testing.T) {
	wr := logsystem.NewWorker("wr", "localhost:9999", "Test.txt", "localhost:9990")

}
