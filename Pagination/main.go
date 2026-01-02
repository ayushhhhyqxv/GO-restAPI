package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

type Product struct {
	ID uint `gorm:"primaryKey" json:"id"`
	Name string `json:"name"`
	Category string `json:"category"`
	Price float64 `json:"price"`
	Description string `json:"description"`
}

func initDB(){

	err:= godotenv.Load()
	if err!=nil{
		log.Fatal("Wrong Credentials provided !")
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

	database.AutoMigrate(&Product{})

	db = database
	
}

func seedData(c echo.Context) error{
	seed:=[]Product{
    {Name: "Iphone", Category: "Mobile", Price: 150000, Description: "Apple Smartphone"},
    {Name: "Samsung", Category: "Laptop", Price: 50000, Description: "Samsung Smartphone"},
    {Name: "Realme", Category: "Headphones", Price: 40000, Description: "Realme Smartphone"},
    {Name: "Nothing", Category: "Laptop", Price: 33000, Description: "Nothing Smartphone"},
    {Name: "Sony", Category: "Headphones", Price: 25000, Description: "Sony Noise Cancelling"},
    {Name: "OnePlus", Category: "Mobile", Price: 65000, Description: "OnePlus Flagship"},
    {Name: "Dell", Category: "Laptop", Price: 85000, Description: "Dell Business Laptop"},
    {Name: "Asus", Category: "Laptop", Price: 110000, Description: "Asus Gaming Laptop"},
    {Name: "Bose", Category: "Headphones", Price: 35000, Description: "Bose Premium Audio"},
    {Name: "Google Pixel", Category: "Mobile", Price: 75000, Description: "Google Android Phone"},
	}
	db.Create(&seed)

	return c.JSON(http.StatusOK,echo.Map{"Success":"Products Data was Successfully Seeded"})

}

func getData(c echo.Context) error {
	pageInfo:= c.QueryParam("page")
	limitParam:= c.QueryParam("limit")
	sortField:= c.QueryParam("sortField")
	sortOrder:= c.QueryParam("sortOrder")
	filter:= c.QueryParam("filter")

	pageNum,err:= strconv.Atoi(pageInfo)
	if err!=nil || pageNum<=0 {
		pageNum = 1
	}
	limit,err:=strconv.Atoi(limitParam)
	if err!=nil || limit<=0 {
		limit = 5
	}

	offset := (pageNum - 1) * limit // elem to skip for specific page 

	query := db.Model(&Product{})
	if filter!=""{
		filterPattern:="%"+ strings.ToLower(filter) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(category) LIKE ?",filterPattern,filterPattern)
	}

	if sortOrder!="" {
		order := "asc"
		if strings.ToLower(sortOrder)=="desc"{
			order = "desc"
		}
		query=query.Order(fmt.Sprintf("%s %s",sortField,order))
	}
	var total int64
	query.Count(&total)

	var product []Product 
	if err:=query.Limit(limit).Offset(offset).Find(&product).Error;err!=nil{
		return c.JSON(http.StatusInternalServerError,echo.Map{"Error":"Failed to fetch data"})
	}
	totalPages := (int(total)+limit-1)/limit

	return c.JSON(http.StatusOK,echo.Map{
		"page": pageNum,
		"limit": limit,
		"total_items": total,
		"total_pages": totalPages,
		"data": product,

	})

}

func main(){
	initDB()
	e := echo.New()
	e.Use(middleware.RequestLogger())
	e.Use(middleware.Recover())
	e.POST("/seed",seedData)
	e.GET("/products",getData)

	port := os.Getenv("SERVER_PORT")
	e.Logger.Fatal(e.Start(":" + port))
}