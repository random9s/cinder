### Recovery!
Middleware package for recoving panic

### Examples

```go
func handler() http.Handler {
    retrun http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	    panic("this should not happen")
    })
}

var panicHandler = func(err interface{}, stacktrace []string, r *http.Request) {
    log.Println(err)
}

func main() {
    var p = Panic(panicHandler)
    http.Handle("/panic", p(handler()))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
