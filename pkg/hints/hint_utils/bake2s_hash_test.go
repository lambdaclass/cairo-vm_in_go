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

func TestBlake2sCompressB(t *testing.T) {
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{456710651, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedState := []uint32{1061041453, 3663967611, 2158760218, 836165556, 3696892209, 3887053585, 2675134684, 2201582556}
	newState := Blake2sCompress(h, message, 2, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}

func TestBlake2sCompressC(t *testing.T) {
	//Hashing "Hello World"
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{1819043144, 1870078063, 6581362, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedState := []uint32{939893662, 3935214984, 1704819782, 3912812968, 4211807320, 3760278243, 674188535, 2642110762}
	newState := Blake2sCompress(h, message, 9, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}

func TestBlake2sCompressD(t *testing.T) {
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{1819043144, 1870078063, 6581362, 274628678, 715791845, 175498643, 871587583, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	expectedState := []uint32{3980510537, 3982966407, 1593299263, 2666882356, 3288094120, 2682988286, 1666615862, 378086837}
	newState := Blake2sCompress(h, message, 28, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}

func TestBlake2sCompressE(t *testing.T) {
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{1819043144, 1870078063, 6581362, 274628678, 715791845, 175498643, 871587583, 635963558, 557369694, 1576875962, 215769785, 0, 0, 0, 0, 0}
	expectedState := []uint32{3251785223, 1946079609, 2665255093, 3508191500, 3630835628, 3067307230, 3623370123, 656151356}
	newState := Blake2sCompress(h, message, 44, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}

func TestBlake2sCompressF(t *testing.T) {
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{1819043144, 1870078063, 6581362, 274628678, 715791845, 175498643, 871587583, 635963558, 557369694, 1576875962, 215769785, 152379578, 585849303, 764739320, 437383930, 74833930}
	expectedState := []uint32{2593218707, 3238077801, 914875393, 3462286058, 4028447058, 3174734057, 2001070146, 3741410512}
	newState := Blake2sCompress(h, message, 64, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}

func TestBlake2sCompressG(t *testing.T) {
	h := [8]uint32{1795745351, 3144134277, 1013904242, 2773480762, 1359893119, 2600822924, 528734635, 1541459225}
	message := [16]uint32{11563522, 43535528, 653255322, 274628678, 73471943, 17549868, 87158958, 635963558, 343656565, 1576875962, 215769785, 152379578, 585849303, 76473202, 437253230, 74833930}
	expectedState := []uint32{3496615692, 3252241979, 3771521549, 2125493093, 3240605752, 2885407061, 3962009872, 3845288240}
	newState := Blake2sCompress(h, message, 64, 0, 4294967295, 0)
	if !reflect.DeepEqual(expectedState, newState) {
		t.Error("Wrong state returned by Blake2sCompress")
	}
}
