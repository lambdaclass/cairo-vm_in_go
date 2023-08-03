package vm

import "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"

type Program struct {
	Data []memory.MaybeRelocatable
}
