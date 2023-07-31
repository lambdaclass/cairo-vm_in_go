package salutes

import (
	"fmt"

	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
)

func Hello() string {
	lambdaworks.From(uint64(5))
	return fmt.Sprintf("Hello, world! Here's your number: %d", lambdaworks.Number())

}
