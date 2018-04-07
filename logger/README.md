### Logger
A simple log file

### Basic Usage:
```go
    import(
      "fyu.se/logfile"
    )

    func main() {
        var l, err = logfile.New("testlog")
	    if err != nil {
		    log.Fatal(err)
	    }
        defer l.Close()

        //General logging use
        l.Trace("Trace level message")
        l.Info("Info level message")
        l.Warning("Warning level message")
        l.Error("Error level message")

        //Stop flow of control and log
        l.Fatal("Fatal level message")
        l.Panic("Panic level message")
    }
```

### Middleware to add logs to context
```go
    import(
      "net/http"

      "fyu.se/logfile"
    )

    type Logger string
    const Key Logger = "access:log"

    func main() {
        var l, err = logfile.New("testlog")
	    if err != nil {
		    log.Fatal(err)
	    }
        defer l.Close()

        copy := CopyLogger(l)
        r := router.New()
        l.Fatal(http.ListenAndServe(":8080", copy(r)))
    }

    //Copy copies logger to request contexts
    func CopyLogger(val logfile.Logger) func(http.Handler) http.Handler {
	    return func(h http.Handler) http.Handler {
		    return func(w http.ResponseWriter, r *http.Request) {
                ctx := r.Context()
	            newCtx := context.WithValue(ctx, Key, val)
	            r = r.WithContext(newCtx)       
            }
        }
    }
```
