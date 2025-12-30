package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq"
)

type User struct {  // to fetch data from database ! 
	ID int `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Age int `json:"age"`
}

func main(){
	e := echo.New()

	e.Use(middleware.RequestLogger())

	dsn := "host=localhost port=5432 user=postgres password=test@123 dbname=newdb password=test@123 sslmode=disable"

	db,err:= sql.Open("postgres",dsn)
	if err!=nil {
		log.Fatal("Error while connecting to P4S: ",err)
	}
	defer db.Close()

	if err:=db.Ping();err!=nil{
		log.Fatal("Database Ping failed: ",err)
	}

	createTable := `
		CREATE TABLE IF NOT EXISTS USERS (
		id serial primary key,
		name text,
		email text unique,
		age int
		);
	`
	if _,err:=db.Exec(createTable);err!=nil {
		log.Fatal("Failed to create the Table: ",err)
	}

	e.POST("/users",func (c echo.Context) error { 
		 // echo.Context holds all the parameters,req body and stuff regarding to that http request and response

		u:=new(User) 

		// creates a pointer to a new, empty User struct. It acts as a container to hold the data coming from the HTTP request body.

		if err:=c.Bind(u);err!=nil { 
			// Bind binds path params, query params and the request body into provided type user struct!
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid JSON format"})
		}

		var id int 

		err:= db.QueryRow(
			"Insert into users(name,email,age) VALUES ($1,$2,$3) Returning id",u.Name,u.Email,u.Age,
		).Scan(&id)

		if err!= nil {
			return c.JSON(http.StatusInternalServerError,map[string]string{"Error" : "Internal Server Error !"})
		}

		u.ID = id 

		return c.JSON(http.StatusCreated,u)
	})

	e.GET("/users",func (c echo.Context) error {
		rows,err:=db.Query("Select id,name,email,age from users")
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"Error" : "Internal Server Error !"})
		}
		defer rows.Close()
		var usersFetch []User // array of user's data for multiple users data is stored here ! 

		for rows.Next(){
			var u User 
			if err:=rows.Scan(&u.ID,&u.Name,&u.Email,&u.Age);err!=nil{
				return c.JSON(http.StatusInternalServerError,map[string]string{"Error" : "Internal Server Error !"})
			}
			usersFetch = append(usersFetch, u)
		}

		return c.JSON(http.StatusOK,usersFetch)
	})

	e.GET("/users/:id",func (c echo.Context) error {
		id,_ := strconv.Atoi(c.Param("id"))
		var u User 

		err:= db.QueryRow("Select id,name,email,age from users where id=$1",id).Scan(&u.ID,&u.Name,&u.Email,&u.Age)
		if err!=nil{
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID provided "})
		}
		return c.JSON(http.StatusOK,u)
	})

	e.PUT("/users/:id",func (c echo.Context) error { // unlike patch here we need every info
		id,_ := strconv.Atoi(c.Param("id"))
		u:= new(User)

		if err:= c.Bind(u);err!=nil{
			return c.JSON(http.StatusBadRequest,map[string]string{"Error" : "Invalid Request Body !"})
		}

		result,err:=db.Exec("Update users set name=$1,email=$2,age=$3 where id=$4",u.Name,u.Email,u.Age,id)
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"Error" : "Internal Server Error !"})
		}

		rowsAffected,_ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound,map[string]string{"Error" : "Internal Server Error !"})
		}
		u.ID = id 
		return c.JSON(http.StatusOK,u)
	})

	e.PATCH("/users/:id",func (c echo.Context) error { // unlike patch here we need every info
		id,_ := strconv.Atoi(c.Param("id"))
		u:= new(User)

		if err:= c.Bind(u);err!=nil{
			return c.JSON(http.StatusBadRequest,map[string]string{"Error" : "Invalid Request Body !"})
		}

		result,err:=db.Exec("Update users set name=$1 where id=$2",u.Name,id)
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"Error" : "Internal Server Error !"})
		}

		rowsAffected,_ := result.RowsAffected()
		if rowsAffected == 0 {
			return c.JSON(http.StatusNotFound,map[string]string{"Error" : "Internal Server Error !"})
		}
		u.ID = id 
		return c.JSON(http.StatusOK,u)
	})

	e.DELETE("/users/:id",func (c echo.Context) error{
		id,_:=strconv.Atoi(c.Param("id"))
		result,err:=db.Exec("Delete from users where id=$1",id)
		if err!=nil{
			return c.JSON(http.StatusBadRequest,map[string]string{"Error" : "Invalid Delete Request !"})
		}
		rowsAffected,_:=result.RowsAffected()
		if rowsAffected==0{
			return c.JSON(http.StatusNotFound,map[string]string{"Error" : " User not found !"})
		}

		return c.JSON(http.StatusOK,map[string]string{"Success" : " User Deleted Successfully !"})

	})

	e.Logger.Fatal(e.Start(":8090"))
}