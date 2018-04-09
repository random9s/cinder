### Responder!
Middleware package for responding to HTTP requests

> Dependencies:
> - gorilla/context

### Examples

```go
import "github.com/random9s/cinder/middleware/response"

func success() http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //Get responder from context
    	ctx := r.Context()
	    v := ctx.Value(response.Key)
	    resp, _ := v.(*response.RespWrapper)

	    //Write successful response
        var t = &struct{
            T bool `json:"T"`
        }{
            T:true,
        }

	    resp.WriteJSON(t)
    })
}

func err() http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        //Get responder from context
    	ctx := r.Context()
	    v := ctx.Value(response.Key)
	    resp, _ := v.(*response.RespWrapper)

        var err = errors.New("bad stuff")
	    resp.Error(err, http.StatusInternalServerError)
    })
}

func main() {
    var encodedWriter = response.Writer(response.GZIP)
    http.Handle("/success", encodedWriter(success()))
    http.Handle("/error", encodedWriter(err()))
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```
