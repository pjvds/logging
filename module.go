package tidy

import (
	"errors"
	"reflect"
	"runtime"
	"strings"
)

type Module struct {
	// the package import path, e.q.: github.com/pjvds/tidy/logentries
	path string
	// the base package path, e.q.: logentries
	name string
}

func NewModule(path string) Module {
	if len(path) == 0 {
		return Module{}
	}

	if lastSlash := strings.LastIndex(path, "/"); lastSlash != -1 {
		return Module{
			path: path,
			name: path[lastSlash:],
		}
	}

	return Module{
		path: path,
	}
}

func GetModuleFromValue(value interface{}) Module {
	return NewModule(reflect.TypeOf(value).PkgPath())
}

func GetModuleFromCaller(depth int) Module {
	pc, _, _, ok := runtime.Caller(1 + depth)

	if !ok {
		panic(errors.New("failed to get caller from runtime"))
	}

	function := runtime.FuncForPC(pc)

	if function == nil {
		panic(errors.New("failed to get function from program counter"))
	}

	// The function name is the complete package path and function name without signature seperated by a dot.
	// e.q.: github.com/pjvds/tidy.GetLogger
	name := function.Name()

	lastDotIndex := strings.LastIndex(name, ".")

	return NewModule(name[0:lastDotIndex])
}

func (this Module) String() string {
	return string(this.name)
}
