package module_manager

func GetTypedConfig[T any](
	m *Manager,
	moduleName string,
) (T, bool) {
	var zero T
	cfg, ok := m.GetConfig(moduleName)
	if !ok {
		return zero, false
	}
	typed, ok := cfg.(T)
	if !ok {
		return zero, false
	}
	return typed, true
}
