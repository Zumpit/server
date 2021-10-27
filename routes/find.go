package routes

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/Zumpit/server/models"
	emailVerifier "github.com/AfterShip/email-verifier"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	//"go.mongodb.org/mongo-driver/bson"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/Zumpit/googlesearch"
	"github.com/Zumpit/domainfinder"
	
)

//ensure the data we are receving from client is correct
var validate = validator.New()
var profileCollection *mongo.Collection = OpenCollection(Client, "profiles")


var (
	verifier = emailVerifier.NewVerifier()
)

const MAX_UPLOAD_SIZE = 1024 * 1024

// marker for progress of file upload
type Progress struct {
	TotalSize float64
	BytesRead float64
}

func (pr *Progress) Write(p []byte) (n int, err error) {
	n, err = len(p), nil
	pr.BytesRead += float64(n)
	pr.Display()
	return
}

func (pr *Progress) Display() {
	if pr.BytesRead == pr.TotalSize {
		fmt.Println("100%")
		return
	}

	percentage := (pr.BytesRead / pr.TotalSize) * 100
	fmt.Printf("File upload in progress : %f\n", percentage)
}

// find linkedin url and title for email address
func GetLinkedin(email string) (string, string) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	res, err := googlesearch.Search(ctx, email)
	if err != nil {
		fmt.Println(err)
	}
	defer cancel()

	var link_url string = res[0].URL
	var link_title string = res[0].Title
    
	return link_url, link_title
}

func ConnectionPage(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"message":" xadhrit was here 🛑" })
}

// find an email

func FindEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var query models.EmailQuery

	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}

	validationErr := validate.Struct(query)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
		return
	}

	var profile models.Profile

	profile.ID = primitive.NewObjectID()
	profile.Email = "dow.mike@aes.com"
	url, title := GetLinkedin(profile.Email)

	profile.Linkedin_Url = url
	profile.Linkedin_Title = title
	data := models.Profile{
		ID:             profile.ID,
		Email:          profile.Email,
		Linkedin_Url:   profile.Linkedin_Url,
		Linkedin_Title: profile.Linkedin_Title,
		EmailQuery:     query,
	}

	result, insertErr := profileCollection.InsertOne(ctx, data)

	if insertErr != nil {
		msg := fmt.Sprintf("profile was not created")
		c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
		fmt.Println("error happing", insertErr)
		return
	}

	defer cancel() 
	c.JSON(http.StatusOK, result)
}

/*
  Email Validation 
*/

func GetEmailValidation(c *gin.Context){

    var email models.EmailValidation
	
	if err := c.BindJSON(&email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Request Error " : err.Error()})
        fmt.Println(err)
        return
	}
	validationErr := validate.Struct(email)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{" Syntax Error": validationErr.Error()})
		return
	}
    e := email.Email 
	result, err := verifier.Verify(e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error ":err.Error()})
		fmt.Println("error happing" , err.Error())
		return
	}
	
	
	c.JSON(http.StatusOK, result)
}


func HandleUpload(c *gin.Context) {
	file, fileHeader, err := c.Request.FormFile("file")

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if fileHeader.Size > MAX_UPLOAD_SIZE {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File size is too big. Please upload file of 1 MB"})
		return
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "error while processing file"})
		return
	}
	// filetype := http.DetectContentType(buffer)
	// if filetype != "" && filetype != "" && filetype != "" {
	// 	c.JSON(http.StatusBadRequest, gin.H{"message": "file type is invalid. Please upload .txt, .lxsl file."})
	// 	return
	// }

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server error"})
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error while creating uploads folter"})
		return
	}

	// file size restiction

	/* file type restriction
	   buffer := make([]byte, 512)
	   _, err = file.Read(buffer)
	*/

	extension := filepath.Ext(fileHeader.Filename)
	newFilename := uuid.New().String() + extension
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), newFilename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error while destination"})
		return
	}

	defer dst.Close()

	pr := &Progress{
		TotalSize: float64(fileHeader.Size),
	}
	_, err = io.Copy(dst, io.TeeReader(file, pr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "error while copying file in server"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Your file has been successfully uploaded"})
}

func FindDomain(c *gin.Context) {
   /*
   input -> Allied Infoline
   output -> www.domain.com,   
   */
    ctx, cancel :=  context.WithTimeout(context.Background(), 100*time.Second)
    var company models.CompanyName
   
    if err := c.BindJSON(&company); err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"JSON Binding Error" : err.Error()})
	    return
    }

	validationErr := validate.Struct(company)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Validation Error": validationErr.Error()})
		return
	}
    searchTerm := company.Name
	domain, err := domainfinder.Search(ctx, searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Finding Error : " : err.Error()})
		return
	}
	defer cancel()
   
   c.JSON(http.StatusOK, domain)

}

func GetDomainValidation(c *gin.Context){
	//ctx, cancel :=  context.WithTimeout(context.Background(), 100*time.Second)
    var domain models.DomainValidation
   
    if err := c.BindJSON(&domain); err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"JSON Binding Error" : err.Error()})
	    return
    }

	validationErr := validate.Struct(domain)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Validation Error": validationErr.Error()})
		return
	}
    searchTerm := domain.Domain 
	result, err := domainfinder.DomainValidation(searchTerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Finding Error : " : err.Error()})
		return
	}
	//defer cancel()
   
   c.JSON(http.StatusOK, result)
	
}

func GetEmailFromDomain(c *gin.Context){
    // ctx, cancel := context.WithTimeout(context.Background(), 100 * time.Second)

	// var domain models.EmailFromDomain

	// if err := c.BindJSON(&domain); err != nil{
	// 	c.JSON(http.StatusBadRequest, gin.H{"JSON Binding Error : ", err.Error()})
	// 	return 
	// }

	// validationErr := validate.Struct(domain)
	// if validationErr != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"Validation Error": validationErr.Error()})
	// 	return
	// }
	// searchTerm := domain.Domain 

	// result, err :=  
    // if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"Finding Error :": err.Error()})
	// 	return 
	// }
	// defer cancel()
	c.JSON(http.StatusOK, gin.H{"message":"get email from domain name"})
} 

func FindCompany(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"write your website.com , link":"Pending working on it"})
}
