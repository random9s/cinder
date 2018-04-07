### Database 
A simple database wrapper


### Basic Usage:
```go
import (
    "fyu.se/database"
    "fyu.se/database/config"
)

func main() {
    //Read db config json from file
	conf, err := ioutil.ReadFile("conf.json")
	if err != nil {
		log.Fatal(err)
	}

    //Register conf as mysql conf
	var sqlConf = &config.MySQL{}
	err = config.Register(conf, sqlConf)
	if err != nil {
		log.Fatal(err)
	}

	//Create db from configuration
	var db = new(database.Mysql)
	database.Register(sqlConf, db)

    //Open connection to db
    conn, err := db.Open()
	if err != nil {
		log.Fatal(err)
	}
    defer db.Close()

    var mysqlConn = conn.(*sql.DB)
    fmt.Println(mysqlConn)
}
```

> * To add a new database you'll need to satisfy both:
    - [config](https://github.com/random9s/fyu.se/blob/master/database/config/config.go) 
    - [database](https://github.com/random9s/fyu.se/blob/master/database/database.go#L5-L11) interfaces. 
