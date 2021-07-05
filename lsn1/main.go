package main

func main() {
	HandlePanic()
	CreateFiles("tmp", 10000)
	RecoverPanicInGoroutine()
}
