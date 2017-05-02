package fibonacci

// F returns the next fibonacci number
func F() func() int {
	return fiboClosure()
}

func fiboRecursive(n int) int {
	if n < 0 {
		return 0
	}
	if n == 1 {
		return 1
	}
	return fiboRecursive(n-1) + fiboRecursive(n-2)
}

func fiboClosure() func() int {
	a, b := 0, 1
	return func() int {
		a, b = b, a+b
		return a
	}
}

func fiboConcurrent(n int) int {
	c := make(chan int, n+1)
	go func(n int, c chan int) {
		x, y := 0, 1
		for i := 0; i <= n; i++ {
			c <- x
			x, y = y, x+y
		}
		close(c)
	}(n, c)
	var i int
	for {
		v, _ := <-c
		if i == n {
			return v
		}
		i++
	}
}
