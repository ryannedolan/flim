# flim
FP in Go

e.g.


    fl.Range(1, 100).Filter(` x < 50 `).Map(` x*x `).Fold(0, ` x + y `).Force()
    -> []interface{}{40425}


You can also apply normal Go functions: [example](samples/evens.go)
