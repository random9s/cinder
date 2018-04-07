### Database 

### Configuring a DB:
```go
import (
    "github.com/random9s/cinder/database"
    "github.com/random9s/cinder/database/config"
)

func main() {
    var mysqlConfig = new(config.MySQL)
	var address = "localhost:3306"
	var user = "root"
	var password = "secure"

    //It's generally better to create a json file and read the configuration from that. 
    //I'm not telling you how to live though 
	err := mysqlConfig.Register([]byte(`{
		"user":"` + user + `",
		"password":"` + password + `",
		"address":"` + address + `",
		"name":"personal"
	}`))
	if err != nil {
		return nil, err
	}

	var db = new(database.Mysql)
	db.Register(mysqlConfig)

}
```

* Current Support:
    - Mysql

* To add a new database you'll need to satisfy both:
    - [config](https://github.com/random9s/fyu.se/blob/master/database/config/config.go) 
    - [database](https://github.com/random9s/fyu.se/blob/master/database/database.go#L5-L11) interfaces. 


### Using the Marshaler

Let's use the following struct as an example:
```go
    type User struct {
	    ID        int64       `json:"id" mysql:"id"`
	    Name      string      `json:"name" mysql:"name"`
	    Email     string      `json:"email mysql:"email"`
	    Username  string      `json:"username" mysql:"username"`
	    Password  string      `json:"password" mysql:"password"`
	    LastUpdated time.Time `json:"updated_on" mysql:"updated_on"`
    }

    type Users []*User
```

```go
import "github.com/random9s/cinder/database/marshal"

func GetUsers() (Users, error) {
    //Assuming this is an actual db connection
	var db = new(database.Mysql)

    //Getting multiple rows
    var sql = `SELECT * FROM user WHERE client_id=? AND deleted=0`
    //Create users slice
	var users = make(Users, 0)
    //Unmarshal into slice
	err := db.UnmarshalRows(&users, sql, cid)
	return users, err
}

func GetUser(id int64) (*User, error) {
    //Assuming this is an actual db connection
	var db = new(database.Mysql)

    //Getting multiple rows
    var sql = `SELECT * FROM user WHERE client_id=? AND deleted=0`
    //Create users slice
	var user = new(User)
    //Unmarshal into slice
	err := db.UnmarshalRow(user, sql, cid)
	return user, err
}
```
