package hint_utils_test

// import (
// 	"testing"

// 	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
// 	"github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
// 	"github.com/lambdaclass/cairo-vm.go/pkg/parser"
// 	"github.com/lambdaclass/cairo-vm.go/pkg/vm"
// 	. "github.com/lambdaclass/cairo-vm.go/pkg/vm/memory"
// )

// func TestGetAddressFromReferenceOff2Value(t *testing.T) {
// 	vm := vm.NewVirtualMachine()
// 	reference := HintReference{
// 		Offset1: OffsetValue{
// 			ValueType: Reference,
// 			Value:     1,
// 		},
// 		Offset2: OffsetValue{
// 			ValueType: Value,
// 			Value:     2,
// 		},
// 	}
// 	// ap + 1 + 2 = (0, 3)
// 	expectedAddr := NewRelocatable(0, 3)
// 	addr, ok := GetAddressFromReference(&reference, parser.ApTrackingData{}, vm)
// 	if addr != expectedAddr || !ok {
// 		t.Errorf("GetAddressFromReference returned wrong result")
// 	}
// }

// func TestGetAddressFromReferenceOff1Immediate(t *testing.T) {
// 	vm := vm.NewVirtualMachine()
// 	reference := HintReference{
// 		Offset1: OffsetValue{
// 			ValueType: Immediate,
// 			Immediate: lambdaworks.FeltFromUint64(17),
// 		},
// 	}
// 	expectedAddr := Relocatable{}
// 	addr, ok := GetAddressFromReference(&reference, parser.ApTrackingData{}, vm)
// 	if addr != expectedAddr || ok {
// 		t.Errorf("GetAddressFromReference returned wrong result")
// 	}
// }

// func TestGetAddressFromReferenceDiffApTrackingGroup(t *testing.T) {
// 	vm := vm.NewVirtualMachine()
// 	reference := HintReference{
// 		Offset1: OffsetValue{
// 			ValueType: Reference,
// 			Immediate: lambdaworks.FeltFromUint64(17),
// 		},
// 		ApTrackingData: parser.ApTrackingData{
// 			Group: 7,
// 		},
// 	}
// 	expectedAddr := Relocatable{}
// 	addr, ok := GetAddressFromReference(&reference, parser.ApTrackingData{}, vm)
// 	if addr != expectedAddr || ok {
// 		t.Errorf("GetAddressFromReference returned wrong result")
// 	}
// }
