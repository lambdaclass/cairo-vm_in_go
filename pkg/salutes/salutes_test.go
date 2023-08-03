package salutes_test

import (
	"testing"

	"github.com/lambdaclass/cairo-vm.go/pkg/salutes"
)

func TestHello(t *testing.T) {
	got := salutes.Hello()
	expected := "Hello, world! Here's your Felt: {[0 0 0 42]}"
	if got != expected {
		t.Errorf("We should have '%s' as the salute, got '%s'", expected, got)
	}
}
