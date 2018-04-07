### Authenticate!
Middleware for Issuing and Validating JWTs 

> Dependencies:
> - dgrijalva/jwt-go
> - gorilla/context

#### JWT Issuer
```go
func main() {
    Issuer := authenticate.IssueJWT(authenticate.NewIssuer(
        "My Server", //Issuer
        "", //Audience
        "Super-Secret-Secret", //JWT Signer Secret
        86400, //Expire After 1 day
        jwt.SigningMethodHMAC, //Specify signing method to use
    ))

    r := mux.NewRouter()
    r.HandleFunc("/signin", Signin) 

    http.ListenAndServe(":8000", Issuer(r))
}

func Signin(w http.ResponseWriter, r *http.Request) {
    //Do Auth stuff

    //Once you have done auth stuff call this method 
    //to store boolean isUserValid in request context
    authenticate.User(r, isUserValid) 
    return
}
```
