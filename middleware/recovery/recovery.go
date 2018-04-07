package recovery

import (
	"fmt"
	"net/http"
	"runtime"
)

//OnPanic ...
type OnPanic func(err interface{}, stacktrace []string, r *http.Request)

type recoverHandler struct {
	h       http.Handler
	onPanic OnPanic
}

//Panic ...
func Panic(p OnPanic) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &recoverHandler{
			h:       h,
			onPanic: p,
		}
	}
}

func (h *recoverHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		err := recover()
		if err != nil {
			var stack []string

			for skip := 1; ; skip++ {
				pc, file, line, ok := runtime.Caller(skip)
				if !ok {
					break
				}

				if file[len(file)-1] != 'c' {
					f := runtime.FuncForPC(pc)
					s := fmt.Sprintf("%s:%d %s()", file, line, f.Name())
					stack = append(stack, s)
				}
			}

			h.onPanic(err, stack, r)
			return
		}
	}()

	h.h.ServeHTTP(w, r)
}
