package fl

type fn func(...interface{}) interface{}

func Compose(f1 fn, f2 fn) fn {
	return func(args ...interface{}) interface{} {
		return f2(f1(args...))
	}
}

func Curry(f fn, args ...interface{}) fn {
	return func(args2 ...interface{}) interface{} {
		return f(append(args, args2...)...)
	}
}
