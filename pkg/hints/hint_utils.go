package hint_utils

import (
	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
)

type HintReference struct {
	offset1          OffsetValue
	offset2          OffsetValue
	dereference      bool
	ap_tracking_data ApTrackingData
}

type ApTrackingData struct {
	group  uint
	offset uint
}

type OffsetValue struct {
	valueType   offsetValueType
	immediate   Felt
	value       uint
	register    vm.Register
	dereference bool
}

type offsetValueType uint

const (
	Immediate offsetValueType = 0
	Value     offsetValueType = 1
	Reference offsetValueType = 2
)
