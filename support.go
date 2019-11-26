package lua

import (
	"sync"
)

var (
	origMathFuncs map[string]LGFunction //< Map to hold original mathlib map
	origOsFuncs   map[string]LGFunction //< Map to hold original oslib map

	mathFuncsMu sync.Mutex //< Mutex for math function map in original mathlib
	osFuncsMu   sync.Mutex //< Mutex for os function map in original oslib

)

func init() {

	// Copy original functions map to respective orig maps

	origMathFuncs = make(map[string]LGFunction, len(mathFuncs))
	for key, val := range mathFuncs {
		origMathFuncs[key] = val
	}

	origOsFuncs = make(map[string]LGFunction, len(osFuncs))
	for key, val := range osFuncs {
		origOsFuncs[key] = val
	}
}

// ILibrary defines the interface for any library
type ILibrary interface {
	// SetFuncs sets the given functions, from among builtin functions, to
	// the library
	// If "all" parameter is passed, all the builtin functions in the
	// library are available to be used in the script
	SetFuncs(funcs ...string)
	// AddFunc adds a func with given name and f as the backend function to
	// be called in the library
	AddFunc(name string, f LGFunction)

	// Open opens the configured functions in the library, to be used in
	// the lua script
	Open(L *LState) int
}

// #region Generics

func setLibFuncs(funcsMap, origFuncs map[string]LGFunction, funcs ...string) {

	if len(funcs) == 1 {
		if funcs[0] == "all" {
			for k, v := range origFuncs {
				funcsMap[k] = v
			}
		}
	} else {
		for _, v := range funcs {
			if f, ok := origFuncs[v]; ok {
				funcsMap[v] = f
			}
		}
	}
}

func addLibFunc(funcsMap map[string]LGFunction, name string, f LGFunction) {
	funcsMap[name] = f
}

// #endregion

// #region MathLib

// MathLib implements ILibrary
type MathLib struct {
	funcsMap map[string]LGFunction
}

// NewMathLib returns a new math library instance
func NewMathLib() *MathLib {
	f := make(map[string]LGFunction)

	return &MathLib{funcsMap: f}
}

func (ml *MathLib) SetFuncs(funcs ...string) {
	setLibFuncs(ml.funcsMap, origMathFuncs, funcs...)
}

func (ml *MathLib) AddFunc(name string, f LGFunction) {
	addLibFunc(ml.funcsMap, name, f)
}

func (ml *MathLib) Open(L *LState) int {

	mathFuncsMu.Lock()
	mathFuncs = ml.funcsMap
	val := OpenMath(L)
	mathFuncsMu.Unlock()

	return val
}

// #endregion

// #region OsLib

// OsLib implements ILibrary
type OsLib struct {
	funcsMap map[string]LGFunction
}

// NewOsLib returns a new math library instance
func NewOsLib() *OsLib {
	f := make(map[string]LGFunction)

	return &OsLib{funcsMap: f}
}

func (ol *OsLib) SetFuncs(funcs ...string) {
	setLibFuncs(ol.funcsMap, origOsFuncs, funcs...)
}

func (ol *OsLib) AddFunc(name string, f LGFunction) {
	addLibFunc(ol.funcsMap, name, f)
}

func (ol *OsLib) Open(L *LState) int {

	osFuncsMu.Lock()
	osFuncs = ol.funcsMap
	val := OpenOs(L)
	osFuncsMu.Unlock()

	return val
}

// #endregion
