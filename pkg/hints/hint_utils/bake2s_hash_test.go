package hint_utils_test

import (
	"reflect"
	"testing"

	. "github.com/lambdaclass/cairo-vm.go/pkg/hints/hint_utils"
)

func TestBlake2sCompressA(t *testing.T) {
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedState := []uint32{412110711, 3234706100, 3894970767, 982912411, 937789635, 742982576, 3942558313, 1407547065}
	newState := Blake2sCompress(h, message, 2, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}
