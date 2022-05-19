package sflag

import (
	"strconv"

	"github.com/gookit/goutil/arrutil"
)

// ValueGetter struct
type ValueGetter struct {
	// value store parsed argument data. (type: string, []string)
	Value interface{}
	// is array
	Arrayed bool
}

// Int argument value to int
func (v *ValueGetter) Int(defVal ...int) int {
	def := 0
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if v.Value == nil || v.Arrayed {
		return def
	}

	if str, ok := v.Value.(string); ok {
		val, err := strconv.Atoi(str)
		if err != nil {
			return val
		}
	}

	return def
}

// String argument value to string
func (v *ValueGetter) String(defVal ...string) string {
	def := ""
	if len(defVal) == 1 {
		def = defVal[0]
	}

	if v.Value == nil || v.Arrayed {
		return def
	}

	if str, ok := v.Value.(string); ok {
		return str
	}

	return def
}

// Ints value to int slice
func (v *ValueGetter) Ints() (ints []int) {
	ints, _ = arrutil.StringsToInts(v.Strings())
	return
}

// Strings value to string slice, if argument isArray = true.
func (v *ValueGetter) Strings() (ss []string) {
	if v.Value != nil && v.Arrayed {
		ss = v.Value.([]string)
	}

	return
}

// Array alias of the Strings()
func (v *ValueGetter) Array() (ss []string) {
	return v.Strings()
}

// HasValue value is empty
func (v *ValueGetter) HasValue() bool {
	return v.Value != nil
}
