package main

import (
	_ "code.google.com/p/odbc"
	"database/sql"
	"fmt"
	"github.com/ant0ine/go-json-rest/rest"
	"log"
	"net/http"
)

var SQLSERVER_HOSTNAME = "TIGER\\SQLEXPRESS"
var USER_ID = "user"
var PASS = "pass"
var PORT = "1433"
var TABLENAME = "temp.dbo.Users2"
var COLNAME = "Name"

var HTTPREST_SERVERPORT = "8082"

type Api struct {
	DB *sql.DB
}

type User struct {
	Id   int
	Name string
}

func (api *Api) initDB() {
	driverString := "driver={sql server};server=" + SQLSERVER_HOSTNAME + ";uid=" + USER_ID + ";pwd=" + PASS + ";port=" + PORT
	db, err := sql.Open("odbc", driverString)
	if err != nil {
		fmt.Println(err)
		return
	}
	api.DB = db
}

func main() {
	api := Api{}
	api.initDB()

	handler := rest.ResourceHandler{
		EnableRelaxedContentType: true,
	}
	handler.SetRoutes(
		rest.RouteObjectMethod("GET", "/api/users", &api, "GetAllUsers"),
		rest.RouteObjectMethod("POST", "/api/users", &api, "PostUser"),
		rest.RouteObjectMethod("GET", "/api/users/:id", &api, "GetUser"),
		/*&rest.RouteObjectMethod("PUT", "/api/users/:id", &api, "PutUser"),*/
		rest.RouteObjectMethod("DELETE", "/api/users/:id", &api, "DeleteUser"),
	)

	http.ListenAndServe(":" + HTTPREST_SERVERPORT, &handler)
}

func (api *Api) GetAllUsers(w rest.ResponseWriter, r *rest.Request) {
	users := make([]User, 0)

	rows, err := api.DB.Query("SELECT * FROM " + TABLENAME)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		users = append(users, User{Id: id, Name: name})
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteJson(&users)
}

func (api *Api) GetUser(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")
	user := &User{}

	rows, err := api.DB.Query("SELECT * FROM " + TABLENAME + " WHERE ID = " + id)
	if err != nil {
		log.Fatal(err)
	}

	if rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			log.Fatal(err)
		}
		user = &User{Id: id, Name: name}
	} else {
		rest.NotFound(w, r)
		return
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteJson(&user)
}

func (api *Api) PostUser(w rest.ResponseWriter, r *rest.Request) {
	user := User{}

	err := r.DecodeJsonPayload(&user)
	if err != nil {
		rest.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lastInsertId, err := api.DB.Query("INSERT INTO " + TABLENAME + "(" + COLNAME + ") OUTPUT Inserted.ID VALUES('" + user.Name + "')")
	if err != nil {
		fmt.Println(err)
		return
	}

	if lastInsertId.Next() {
		var id int
		if err := lastInsertId.Scan(&id); err != nil {
			log.Fatal(err)
		}
		user.Id = id
	} else {
		rest.NotFound(w, r)
		return
	}
	if err := lastInsertId.Err(); err != nil {
		log.Fatal(err)
	}

	w.WriteJson(&user)
}

/*func (api *Api) PutUser(w rest.ResponseWriter, r *rest.Request) {
    id := r.PathParam("id")
    if self.Store[id] == nil {
        rest.NotFound(w, r)
        return
    }
    user := User{}
    err := r.DecodeJsonPayload(&user)
    if err != nil {
        rest.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    user.Id = id
    self.Store[id] = &user
    w.WriteJson(&user)
}*/

func (api *Api) DeleteUser(w rest.ResponseWriter, r *rest.Request) {
	id := r.PathParam("id")

	_, err := api.DB.Exec("DELETE FROM " + TABLENAME + " WHERE ID = " + id)
	if err != nil {
		fmt.Println(err)
		return
	}

}
