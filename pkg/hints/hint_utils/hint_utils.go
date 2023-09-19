package hint_utils

import (
	"strings"

	. "github.com/lambdaclass/cairo-vm.go/pkg/lambdaworks"
	"github.com/pkg/errors"
)

func GetConstantFromVarName(name string, constants *map[string]Felt) (Felt, error) {
	if constants == nil {
		return Felt{}, errors.Errorf("Caled GetConstantFromVarName with a nil constants map. Var Name: %s", name)
	}

	for key, value := range *constants {
		keySplit := strings.Split(key, ".")
		if keySplit[len(keySplit)-1] == name {
			return value, nil
		}
	}

	return Felt{}, errors.Errorf("Variable name not found in constants map. Var Name: %s", name)
}
