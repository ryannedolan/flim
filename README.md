# flim
FP in Go

Flim turns strings like "x+y" into lambdas:

    foo := fl.Lambda(`x + y`)
    foo(1, 2)
    -> 3

You can apply lambdas to chans:

    ch := fl.Range(1, 100).Chan()
    fl.Iter(ch).Filter(`x < 50`).Map(`x*x`).Fold(0, `x + y`).List()
    -> List(40425)

You can also apply normal Go functions:

    items := []int{1, 2, 3, 4}
    fl.Iter(items).Map(func (i int) int { return i + 1 }).Array()
    -> {2, 3, 4, 5}


    
