## Redirect!
Middleware handler for redirecting http to https

```go
func handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, r.Proto)
		return
	})
}

func main() {
    //OnHTTPS takes bool arg to enable redirect, false in development 
    var redirect = OnHTTPS(true)
    http.Handle("/", redirect(handle()))

    go log.Fatal(http.ListenAndServe(":80", nil))
    log.Fatal(http.ListenAndServeTLS(":443", "cert.pem", "key.pem", nil))
}
```
