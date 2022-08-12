package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/minio/minio-go/v7"
)

type jwtCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

const minioURL = "localhost:9000"
const minioUser = "ROOTUSERNAME"
const minioPass = "MEMEPASSWORD"
const bucketName = "react-shop"

func upload(c echo.Context) error {
	ctx := context.Background()

	// MinIO Connection
	minioClient := getMinio().minioInstance

	// PostgreSQL Connection
	db := getPostgres().postgresInstance

	title := c.FormValue("title")
	description := c.FormValue("description")

	price, err := strconv.ParseFloat(c.FormValue("price"), 32)
	if err != nil {
		fmt.Println("ParseFloat: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	file, err := c.FormFile("picture")
	if err != nil {
		fmt.Println("File: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		fmt.Println("Src: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer src.Close()
	dst, err := os.Create("./" + file.Filename)
	if err != nil {
		fmt.Println("Dst: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		fmt.Println("Cpy: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	minioClient.FPutObject(ctx, bucketName, file.Filename, "./"+file.Filename, minio.PutObjectOptions{})
	if err != nil {
		fmt.Println("Minio: ", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	os.Remove("./" + file.Filename)

	db.Create(&Item{Title: title, Description: description, Price: price, Picture: "http://" + minioURL + "/" + bucketName + "/" + file.Filename})

	return c.String(http.StatusOK, "Item Uploaded Successfully")

}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	r := e.Group("/upload")
	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte("windows_sucks"),
	}
	r.Use(middleware.JWTWithConfig(config))

	// Routes
	e.GET("/items", func(c echo.Context) error {
		var items []Item
		db := getPostgres().postgresInstance
		db.Find(&items)
		a, err := json.Marshal(items)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}
		n := len(a)
		s := string(a[:n])
		return c.String(http.StatusOK, s)
	})
	r.POST("", upload)

	// Server Start
	e.Logger.Fatal(e.Start(":1323"))
}
