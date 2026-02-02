package module_manager

import (
	"reflect"
	"unsafe"
)

func scanDependencies(mod Module) []string {
	seen := make(map[string]bool)
	var deps []string

	scanValue(reflect.ValueOf(mod), mod.Name(), seen, &deps, 0)

	return deps
}

func scanValue(
	val reflect.Value,
	selfName string,
	seen map[string]bool,
	deps *[]string,
	depth int,
) {
	if depth > 10 {
		return
	}

	if val.Kind() == reflect.Ptr {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if val.Kind() == reflect.Interface {
		if val.IsNil() {
			return
		}
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if !field.CanInterface() {
			field = forceInterface(field)
		}

		if !field.CanInterface() {
			continue
		}

		iface := field.Interface()

		if depMod, ok := iface.(Module); ok && depMod != nil {
			depName := depMod.Name()
			if depName != selfName && !seen[depName] {
				seen[depName] = true
				*deps = append(*deps, depName)
			}
			continue
		}

		switch field.Kind() {
		case reflect.Ptr:
			if !field.IsNil() {
				scanValue(field, selfName, seen, deps, depth+1)
			}
		case reflect.Struct:
			scanValue(field, selfName, seen, deps, depth+1)
		}
	}
}

//nolint:gosec
func forceInterface(field reflect.Value) reflect.Value {
	return reflect.NewAt(field.Type(), unsafe.Pointer(field.UnsafeAddr())).Elem()
}
