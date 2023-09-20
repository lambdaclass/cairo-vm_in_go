package types

func NewExecutionScopesWithInitValue(key string, val interface{}) *ExecutionScopes {
	scopes := NewExecutionScopes()
	scopes.EnterScope(map[string]interface{}{key: val})
	return scopes
}
