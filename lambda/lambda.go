// convenience aliases for creating lambdas
package lambda

import (
	"github.com/ryannedolan/flim/fl"
)

func X(expr string) func(interface{}) interface{} {
	f := fl.Lambda(expr)
	return func(x interface{}) interface{} {
		return f(fl.X(x))
	}
}

func XY(expr string) func(interface{}, interface{}) interface{} {
	f := fl.Lambda(expr)
	return func(x interface{}, y interface{}) interface{} {
		return f(fl.XY(x, y))
	}
}

func XYZ(expr string) func(interface{}, interface{}, interface{}) interface{} {
	f := fl.Lambda(expr)
	return func(x interface{}, y interface{}, z interface{}) interface{} {
		return f(fl.XYZ(x, y, z))
	}
}
