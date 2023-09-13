package types

import (
	"reflect"

	"github.com/pkg/errors"
)

type ExecutionScopes struct {
	data []map[string]interface{}
}

var ErrCannotExitMainScop error = ExecutionScopesError(errors.Errorf("Cannot exit main scope."))

func ExecutionScopesError(err error) error {
	return errors.Wrapf(err, "Execution scopes error")
}

func ErrVariableNotInScope(varName string) error {
	return ExecutionScopesError(errors.Errorf("Variable %s not in scope", varName))
}

func NewExecutionScopes() *ExecutionScopes {
	data := make([]map[string]interface{}, 1)
	data[0] = make(map[string]interface{})
	return &ExecutionScopes{data}
}

func (es *ExecutionScopes) Data() []map[string]interface{} {
	return es.data
}

func (es *ExecutionScopes) EnterScope(newScopeLocals map[string]interface{}) {
	es.data = append(es.data, newScopeLocals)

}

func (es *ExecutionScopes) ExitScope() error {
	if len(es.Data()) < 2 {
		return ErrCannotExitMainScop
	}
	i := len(es.Data()) - 1
	es.data = append(es.Data()[:i], es.Data()[i+1:]...)

	return nil
}

func (es *ExecutionScopes) getLocalVariablesMut() (*map[string]interface{}, error) {
	if len(es.data) > 0 {
		return &es.data[len(es.data)-1], nil
	}
	return nil, ExecutionScopesError(errors.Errorf("Every enter_scope() requires a corresponding exit_scope()."))
}

func (es *ExecutionScopes) GetLocalVariables() (map[string]interface{}, error) {
	if len(es.data) > 0 {
		return es.data[len(es.data)-1], nil
	}
	return nil, ExecutionScopesError(errors.Errorf("Every enter_scope() requires a corresponding exit_scope()."))
}

func (es *ExecutionScopes) DeleteVariable(varName string) {
	locals, err := es.getLocalVariablesMut()
	if err != nil {
		return
	}
	delete(*locals, varName)

}

func (es *ExecutionScopes) AssignOrUpdateVariable(varName string, varValue interface{}) {
	locals, err := es.getLocalVariablesMut()
	if err != nil {
		return
	}
	(*locals)[varName] = varValue
}

func (es *ExecutionScopes) Get(varName string) (interface{}, error) {
	locals, err := es.GetLocalVariables()
	if err != nil {
		return nil, err
	}
	val, prs := locals[varName]
	if !prs {
		return nil, ErrVariableNotInScope(varName)
	}
	return val, nil
}

func (es *ExecutionScopes) GetRef(varName string) (*interface{}, error) {
	locals, err := es.GetLocalVariables()
	if err != nil {
		return nil, err
	}
	val, prs := locals[varName]
	if !prs {
		return nil, ErrVariableNotInScope(varName)
	}
	return &val, nil
}

// This error has been used for testing purposes
// We should not give information about exists a list with different types on the scope
// On production should be removed
func ErrListTypeNotEqual(varName string, expectedType string, resultType string) error {
	return ExecutionScopesError(errors.Errorf("List %s types does not match, expected type: %s, result type: %s", varName, expectedType, resultType))
}

func (es *ExecutionScopes) GetList(T reflect.Type, varName string) (interface{}, error) {
	maybeList, err := es.Get(varName)
	if err != nil {
		return nil, ErrVariableNotInScope(varName)
	}

	// Correct code
	if reflect.TypeOf(maybeList) != T {
		return nil, ErrVariableNotInScope(varName)
	}

	// If uncommented, should comment correct code
	// Extra code
	// if reflect.TypeOf(maybeList) != T {
	// 	return nil, ErrListTypeNotEqual(varName, T.String(), reflect.TypeOf(maybeList).String())
	// }

	return maybeList, nil

}
