package types

import (
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

func ErrVariableHasWrongType(varName string) error {
	return ExecutionScopesError(errors.Errorf("Scope variable %s has wrong type", varName))
}

func NewExecutionScopes() *ExecutionScopes {
	data := make([]map[string]interface{}, 1)
	data[0] = make(map[string]interface{})
	return &ExecutionScopes{data}
}

func (es *ExecutionScopes) EnterScope(newScopeLocals map[string]interface{}) {
	es.data = append(es.data, newScopeLocals)

}

func (es *ExecutionScopes) ExitScope() error {
	if len(es.data) < 2 {
		return ErrCannotExitMainScop
	}
	i := len(es.data) - 1
	es.data = es.data[:i]

	return nil
}

func (es *ExecutionScopes) getLocalVariablesMut() (*map[string]interface{}, error) {
	locals, err := es.GetLocalVariables()
	if err != nil {
		return nil, err
	}
	return &locals, nil
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

// Generic version of ExecutionScopes.Get which also handles casting
// Also works if the scope variable has type *T instead of T
func FetchScopeVar[T interface{}](varName string, scopes *ExecutionScopes) (T, error) {
	locals, err := scopes.GetLocalVariables()
	if err != nil {
		return *new(T), err
	}
	valAny, prs := locals[varName]
	if !prs {
		return *new(T), ErrVariableNotInScope(varName)
	}
	val, ok := valAny.(T)
	if !ok {
		val, ok := valAny.(*T)
		if ok {
			return *val, nil
		}
		return *new(T), ErrVariableHasWrongType(varName)
	}
	return val, nil
}
