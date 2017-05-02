#Fibonacci [![Circle CI](https://circleci.com/bb/_paulo/fibonacci.svg?style=svg)](https://circleci.com/bb/_paulo/fibonacci)

Test and benchmark several implementations of the fibonacci numbers in Go

```
$ GOMAXPROCS=1 go test . -bench=. -benchmem
BenchmarkFiboRecursive           3000000               471 ns/op               0 B/op          0 allocs/op
BenchmarkFiboClosure            20000000                84.7 ns/op            48 B/op          3 allocs/op
BenchmarkFiboConcurrent          2000000               933 ns/op             192 B/op          1 allocs/op
```
