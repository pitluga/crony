package crony

import (
	"strings"
	"testing"
)

func TestSink(t *testing.T) {
	reader := strings.NewReader("this is a string")
	sink := NewSink()
	sink.ConsumeReader(reader)

	if sink.Content != "this is a string" {
		t.Error("sink did not properly consume reader")
	}
}
