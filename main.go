package main

import (
	"log"
	"os"

	//"net/http"
	"server/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	router := gin.New()
	router.Use(gin.Logger())

	// router.Use(cors.New(cors.Config{
	//     AllowOrigins:     []string{"http://127.0.0.1:3000/*"},
	//     AllowMethods:     []string{http.MethodGet, http.MethodPatch, http.MethodPost, http.MethodHead, http.MethodDelete, http.MethodOptions},
	//     AllowHeaders:     []string{"Content-Type", "Accept", "Origin", "X-Requested-With"},
	//     ExposeHeaders:    []string{"Content-Length"},
	//     AllowCredentials: true,
	// }))
	router.Use(cors.Default())
    router.GET("/", routes.ConnectionPage)
	router.POST("/findEmail", routes.FindEmail)
	router.GET("/validateEmail", routes.GetEmailValidation)
    router.GET("/emailfromDomain", routes.GetEmailFromDomain)
	router.POST("/upload", routes.HandleUpload)
	router.GET("/validateDomain",routes.GetDomainValidation)
	router.GET("/findDomain", routes.FindDomain)
	router.GET("/findCompany", routes.FindCompany)
	

	router.Run(":" + port)
}
