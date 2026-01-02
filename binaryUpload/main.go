package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type FileUpload struct {
	gorm.Model
	Filename string `json:"file_name,omitempty"`
	Filetype string `json:"file_type,omitempty"`
	Filedata []byte `json:"-" gorm:"type:bytea"`
}

var db *gorm.DB

func initDB(){
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error reading .env file !")
	}

	dsn:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
	os.Getenv("DB_HOST"),
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),
	os.Getenv("DB_PORT"),
		)

	database,err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err!=nil{
		log.Fatal("Wrong credentials of DB provided ! ")
	}

	database.AutoMigrate(&FileUpload{})

	db = database
}

func uploadFile(c echo.Context) error {
	file,err:= c.FormFile("file")
	if err!=nil{
		return c.JSON(http.StatusBadRequest,echo.Map{"Error":"File format is required"})
	}
	src,err:=file.Open()
	if err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Cannot open the file"})
	}
	defer src.Close()

	fileBytes,err := io.ReadAll(src)
	if err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"While Reading the file"})
	}
	uploadfile := FileUpload{
		Filename: file.Filename,
		Filetype: file.Header.Get("Content-Type"),
		Filedata: fileBytes,
	}
	if err:=db.Create(&uploadfile).Error;err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Refused to store data ! "})
	}
	return c.JSON(http.StatusOK,echo.Map{
		"Success":"File Saved Successfully",
		"file_name":uploadfile.Filename,
		"file_type":uploadfile.Filetype,
	})
}

func getFile(c echo.Context) error {
	id := c.Param("id")
	var file FileUpload

	if err:=db.First(&file,id).Error;err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"File Not Found !"})
	}
	return c.Stream(http.StatusOK,file.Filetype,bytes.NewReader(file.Filedata))
	
}

func main(){
	initDB()
	e:= echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())

	e.POST("/upload", uploadFile)
	e.GET("/file/:id",getFile)
	
	e.Logger.Fatal(e.Start(":8082"))
}