## Server
Reduce redundancy in future Web API projects.

### Basic Usage:
```go
    import (
        "github.com/gorilla/mux"
        "github.com/random9s/cinder/server"
    )

    func main() {
        router := mux.NewRouter()
        srv := server.New(router)
        srv.Run()    
    }
```

> A server is initialized with a [configuration](https://github.com/random9s/cinder/blob/master/server/config.go#L19-L32) which is created from either a json formatted config file or from flags on the command line. 

> The database configuration will also be stored under the server config.  

### Configuration:
> * Port <string> -- Sets the port where the server will listen
> * TLSCert <string> -- Sets path to a TLS Certificate
> * TLSKey <string> -- Sets path to a TLS Key
