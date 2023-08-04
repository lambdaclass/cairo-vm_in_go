package salutes

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

func Hello() string {
	return fmt.Sprintf("Hello, world! Here's your Felt: %d", lambdaworks.FeltFromUint64(42))
}
