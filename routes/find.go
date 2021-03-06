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
var companyCollection *mongo.Collection = OpenCollection(Client, "companies") 

var (
	verifier = emailVerifier.NewVerifier()
)

const MAX_UPLOAD_SIZE = 1024 * 1024 * 1024

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
	fmt.Printf("File upload in progress : %f \n ", percentage)
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

// find domain/website and name from name
func GetDomain(searchTerm string) (string, string){
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	
	opts := domainfinder.SearchOptions{
		Limit: 100,
	}
	res, err := domainfinder.Search(ctx, searchTerm, opts)
	if err != nil {
		fmt.Println(err)
	} 
	defer cancel()

	var domain_link string = res[0].URL 
	var domain_title string = res[0].Title 

	return domain_link, domain_title 
}


/*
Connection PAGE
*/

func ConnectionPage(c *gin.Context){
	c.JSON(http.StatusOK, gin.H{"message":" xadhrit was here, it's OCT 2????21" })
}

// find an email

func FindEmail(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var query models.EmailQuery
    
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error":"Request method not allowed!"} )
		return 
	} 
	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		fmt.Println(err)
		return
	}
    
	validationErr := validate.Struct(query)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Validation Error error": validationErr.Error()})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		fmt.Println("error happing", insertErr)
		return
	}

	defer cancel() 
	c.JSON(http.StatusOK, result)
}

/*
  Email Validation 
*/
// add context
func GetEmailValidation(c *gin.Context){
    ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
    fmt.Println(ctx)
	
	var email models.EmailValidation
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error":"Request method not allowed!"} )
		return 
	} 

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
    fmt.Println(e)
	result, err := verifier.Verify(e)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Internal Server Error ":err.Error()})
		fmt.Println("error happing" , err.Error())
		return
	}
	
	defer cancel()
	c.JSON(http.StatusOK, result)
}

/*
   File Upload Handler
*/

func HandleUpload(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
	fmt.Println(ctx)

	file, fileHeader, err := c.Request.FormFile("fileop")
	
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
    fmt.Println(file)
	if fileHeader.Size > MAX_UPLOAD_SIZE {
		c.JSON(http.StatusBadRequest, gin.H{"message": "File size is too big. Please upload file of 1 MB"})
		return
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	
	filetype := http.DetectContentType(buffer)
	if filetype != "application/octet-stream" && filetype != "application/zip" && filetype != "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet" {
	 	c.JSON(http.StatusBadRequest, gin.H{"message": "file type is invalid. Please upload .txt, .xls file."})
	 	return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	err = os.MkdirAll("./uploads", os.ModePerm)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	extension := filepath.Ext(fileHeader.Filename)
	newFilename := uuid.New().String() + extension
	dst, err := os.Create(fmt.Sprintf("./uploads/%d%s", time.Now().UnixNano(), newFilename))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error ": err.Error()})
		return
	}

	defer dst.Close()

	pr := &Progress{
		TotalSize: float64(fileHeader.Size),
	}
	_, err = io.Copy(dst, io.TeeReader(file, pr))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"Validation Error ":err.Error()})
		return
	}
	defer cancel()
	c.JSON(http.StatusOK, gin.H{"message": "Your file has been successfully uploaded"})
}

/*
 Find available Domain from name from godaddy, hostinger and etc sites
*/

func FindDomain(c *gin.Context) {
   /*
   input -> Allied Infoline
   output -> www.allied-infoline.com,   
   */
   
    ctx, cancel :=  context.WithTimeout(context.Background(), 100*time.Second)
    fmt.Println(ctx)
    var company models.CompanyQuery 
    if err := c.BindJSON(&company); err != nil {
	    c.JSON(http.StatusBadRequest, gin.H{"JSON Binding Error" : err.Error()})
	    return
    }

	validationErr := validate.Struct(company)
	if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Validation Error": validationErr.Error()})
		return
	}

	var domain_res models.CompanyResult
	e := company.Name

	domain, company_name :=  GetDomain(e)
    fmt.Println(company_name)
	domain_res.Domain = domain 
	domain_res.Name = e 
    
	result := models.CompanyResult {
       Domain: domain_res.Domain,
	   Name: domain_res.Name,
	}

	defer cancel()

    c.JSON(http.StatusOK, result)
}

func GetDomainValidation(c *gin.Context){
	//ctx, cancel :=  context.WithTimeout(context.Background(), 100*time.Second)
    var domain models.DomainValidation
   
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error":"Request method not allowed!"} )
	    return 
	}

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
	//context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)

	var name models.CompanyQuery
    
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error":"Request method not allowed!"} )
	    return 
	}
	
	if err := c.BindJSON(&name); err != nil {
		c.JSON(http.StatusBadRequest , gin.H{"message":err.Error()})
		return
	}

	validationErr := validate.Struct(name)	
    if validationErr != nil {
		c.JSON(http.StatusBadRequest, gin.H{"Validation Error":validationErr.Error()})
		return 
	} 
    
	// Create Company Profile 
	var company models.CompanyProfile 
	company.ID = primitive.NewObjectID()
    
	e := name.Name 
	//find linkedin url 
	url, title := GetLinkedin(e)
	fmt.Println(title)

	//find domain + name (in mannered form) / website url
    domain, company_name :=  GetDomain(e)

	company.Domain = domain 
	company.Name = company_name  
	company.Linkedin_Url = url
	
	data := models.CompanyProfile {
		ID:             company.ID,
		Name :          company.Name,  
 		Domain :        company.Domain, 
		Linkedin_Url:   company.Linkedin_Url,
	}
    
	result, insertErr := companyCollection.InsertOne(ctx, data)
	
	if insertErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": insertErr.Error()})
		fmt.Println("error happing", insertErr)
		return
	}
	defer cancel()

	c.JSON(http.StatusOK, result)
}

// func  FindProfile(c *gin.Context){
//    c.JSON(http.StatusOK, gin.H{"msg":"supply chain"})
// } 