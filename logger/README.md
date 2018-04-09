### Logger
A simple log file

### Basic Usage:
```go
    import(
      "github.com/random9s/cinder/logger"
    )

    func main() {
        var l, err = logger.New("testlog")
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
