package main

import (
	"errors"
	"fmt"
	"time"
)

type f func()

func Def(fn f) {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("panic")
			}
		}()

		fn()
	}()
}

func RecoverPanicInGoroutine() {
	Def(foo)
	time.Sleep(5 * time.Second)
}

func foo() {
	panic(errors.New("error"))
}

//func Run() (err error) {
//	defer func() {
//		if r := recover(); r != nil {
//			switch t := r.(type) {
//			case string:
//				err = error.New(t)
//			case error:
//				err = t
//			}
//		}
//	}()
//	panic("A panic")
//}
