package main

func main() {
	HandlePanic()
	TestTracedError()
	CreateFiles("tmp", 1000)
	RecoverPanicInGoroutine()
}
