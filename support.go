package lua

import "sync"

var (
	origMathFuncs map[string]LGFunction //< Map to hold original mathlib map
	origOsFuncs   map[string]LGFunction //< Map to hold original oslib map

	mathFuncsMu sync.Mutex //< Mutex for math function map in original mathlib
	osFuncsMu   sync.Mutex //< Mutex for os function map in original oslib

)

func init() {

	// Copy original functions map to respective orig maps
	copyFuncsMap(mathFuncs, origMathFuncs)
	copyFuncsMap(osFuncs, origOsFuncs)
}

func copyFuncsMap(src, dst map[string]LGFunction) {
	dst = make(map[string]LGFunction, len(src))
	for key, val := range src {
		dst[key] = val
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

func openLib(funcsMap map[string]LGFunction,
	origFuncs map[string]LGFunction,
	open LGFunction,
	mu *sync.Mutex,
	L *LState) int {

	mu.Lock()
	origFuncs = funcsMap
	val := open(L)
	mu.Unlock()

	return val

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
	return openLib(ml.funcsMap, mathFuncs, OpenMath, &mathFuncsMu, L)
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
	return openLib(ol.funcsMap, osFuncs, OpenOs, &osFuncsMu, L)
}

// #endregion
