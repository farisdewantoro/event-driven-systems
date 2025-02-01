package util

import (
	"errors"
	"fmt"
	"runtime/debug"

	"github.com/labstack/gommon/log"
)

func SafeGoroutine(fn func()) {
	var err error
	go func() {
		defer func() {
			if r := recover(); r != nil {
				var ok bool
				err, ok = r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				log.Error(errors.New("goroutine panic: " + err.Error()))
				log.Error(errors.New("goroutine stacktrace: \n" + string(debug.Stack())))
			}
		}()
		fn()
	}()
}
