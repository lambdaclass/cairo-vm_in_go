package types

type ExecutionScopes struct {
	data []map[string]interface{}
}

func NewExecutionScopes() *ExecutionScopes {
	return &ExecutionScopes{
		data: make([]map[string]interface{}, 0),
	}
}

func (es *ExecutionScopes) enterScope(newScopeLocals map[string]interface{}) {
	es.data = append(es.data, newScopeLocals)
}

func (es *ExecutionScopes) exitScope() error {
	if len(es.data) == 1 {
		return ExecutionScopesError("Cannot exit main scope.")
	}
	es.data = es.data[:len(es.data)-1]
	return nil
}

func (es *ExecutionScopes) getLocalVariablesMut() (*map[string]interface{}, error) {
	if len(es.data) > 0 {
		return &es.data[len(es.data)-1], nil
	}
	return nil, ExecutionScopesError("Every enter_scope() requires a corresponding exit_scope().")
}

func (es *ExecutionScopes) getLocalVariables() (map[string]interface{}, error) {
	if len(es.data) > 0 {
		return es.data[len(es.data)-1], nil
	}
	return nil, ExecutionScopesError("Every enter_scope() requires a corresponding exit_scope().")
}

func (es *ExecutionScopes) deleteVariable(varName string) {
	locals, err := es.getLocalVariablesMut()
	if err != nil {
		return
	}
	delete(*locals, varName)

}

func (es *ExecutionScopes) assignOrUpdateVariable(varName string, varValue interface{}) {
	locals, err := es.getLocalVariablesMut()
	if err != nil {
		return
	}
	(*locals)[varName] = varValue
}
