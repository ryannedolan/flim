// functional programming primitives on Go chans
package fl

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

// limited to a small number of values to avoid map lookups
type env struct {
	X interface{}
	Y interface{}
	Z interface{}
}

// construct a env with X = x
func X(x interface{}) *env {
	return &env{X: x}
}

// construct a env with X = x, Y = y
func XY(x interface{}, y interface{}) *env {
	return &env{X: x, Y: y}
}

// construct a env with X = x, Y = y, Z = z
func XYZ(x interface{}, y interface{}, z interface{}) *env {
	return &env{X: x, Y: y, Z: z}
}

// evaluate an expression in the given env
func Eval(expr string, e *env) interface{} {
	return compileString(expr)(e)
}

func compileString(s string) func(*env) interface{} {
	e, err := parser.ParseExpr(s)
	if err != nil {
		panic(err)
	}
	return compile(e)
}

func compile(e ast.Expr) func(*env) interface{} {
	switch e.(type) {
	case *ast.BinaryExpr:
		return binaryExpr(e.(*ast.BinaryExpr))
	case *ast.BasicLit:
		return basicLit(e.(*ast.BasicLit))
	case *ast.Ident:
		return ident(e.(*ast.Ident))
	default:
		panic(fmt.Errorf("dunno how to compile %v", e))
	}
}

func asFloat64(a interface{}) float64 {
	switch a.(type) {
	case float64:
		return a.(float64)
	case float32:
		return float64(a.(float32))
	case int:
		return float64(a.(int))
	default:
		panic(fmt.Errorf("expected float64, got %v", a))
	}
}

func asBool(a interface{}) bool {
	switch a.(type) {
	case bool:
		return a.(bool)
	default:
		panic(fmt.Errorf("expected bool, got %v", a))
	}
}

func asInt(a interface{}) int {
	switch a.(type) {
	case int:
		return a.(int)
	case float64:
		coerced := int(a.(float64))
		if float64(coerced) != a.(float64) {
			panic(fmt.Errorf("can't coerce %v into an int", a))
		}
		return coerced
	default:
		panic(fmt.Errorf("expected int, got %v", a))
	}
}

func wrap(a interface{}) func(e *env) interface{} {
	return func(e *env) interface{} { return a }
}

func lookup(ident string) func(e *env) interface{} {
	switch ident {
	case "x":
		return func(e *env) interface{} { return e.X }
	case "y":
		return func(e *env) interface{} { return e.Y }
	case "z":
		return func(e *env) interface{} { return e.Z }
	default:
		panic(fmt.Errorf("unknown identifier %s", ident))
	}
}

func asString(a interface{}) string {
	return fmt.Sprintf("%v", a)
}

func basicLit(e *ast.BasicLit) func(*env) interface{} {
	switch e.Kind {
	case token.FLOAT:
		f, err := strconv.ParseFloat(e.Value, 64)
		if err != nil {
			panic(err)
		}
		return wrap(f)
	case token.INT:
		i, err := strconv.Atoi(e.Value)
		if err != nil {
			panic(err)
		}
		return wrap(i)
	case token.STRING:
		s, err := strconv.Unquote(e.Value)
		if err != nil {
			panic(err)
		}
		return wrap(s)
	default:
		panic(fmt.Errorf("unknown literal %v", e.Value))
	}
}

func ident(e *ast.Ident) func(e *env) interface{} {
	return lookup(e.Name)
}

func binaryExpr(a *ast.BinaryExpr) func(*env) interface{} {
	x := compile(a.X)
	y := compile(a.Y)
	apply := func(f func(interface{}, interface{}) interface{}) func(e *env) interface{} {
		return func(e *env) interface{} { return f(x(e), y(e)) }
	}
	switch a.Op {
	case token.ADD:
		return apply(add)
	case token.SUB:
		return apply(sub)
	case token.MUL:
		return apply(mul)
	case token.LSS:
		return apply(lss)
	case token.GTR:
		return apply(gtr)
	case token.LEQ:
		return apply(leq)
	case token.GEQ:
		return apply(geq)
	default:
		panic(fmt.Errorf("unknown operator %v", a.Op))
	}
}

func add(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) + asFloat64(b)
	case int:
		return a.(int) + asInt(b)
	case string:
		return asString(a) + asString(b)
	default:
		panic(fmt.Errorf("can't add these: %v + %v", a, b))
	}
}

func sub(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) - asFloat64(b)
	case int:
		return a.(int) - asInt(b)
	default:
		panic(fmt.Errorf("can't subtract these: %v - %v", a, b))
	}
}

func mul(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) * asFloat64(b)
	case int:
		return a.(int) * asInt(b)
	default:
		panic(fmt.Errorf("can't multiply these: %v * %v", a, b))
	}
}

func gtr(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) > asFloat64(b)
	case int:
		return a.(int) > asInt(b)
	default:
		panic(fmt.Errorf("can't compare these: %v > %v", a, b))
	}
}

func geq(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) >= asFloat64(b)
	case int:
		return a.(int) >= asInt(b)
	default:
		panic(fmt.Errorf("can't compare these: %v >= %v", a, b))
	}
}

func lss(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) < asFloat64(b)
	case int:
		return a.(int) < asInt(b)
	default:
		panic(fmt.Errorf("can't compare these: %v < %v", a, b))
	}
}

func leq(a interface{}, b interface{}) interface{} {
	switch a.(type) {
	case float64:
		return a.(float64) <= asFloat64(b)
	case int:
		return a.(int) <= asInt(b)
	default:
		panic(fmt.Errorf("can't compare these: %v <= %v", a, b))
	}
}
