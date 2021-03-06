package models 

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type CompanyQuery struct {
	Name   string   `json:"name" binding:"required`
}

type CompanyProfile struct {
	ID               primitive.ObjectID  `bson:"_id"`                 
	Name             string   `json:"name"`
    Domain           string   `json:"domain"`
	Avatar           string   `json:"avatar"`
	Size             string   `json:"size"`
	Hq               string   `json:"hq"`
	Industry         string   `json:"industry"`
	Revenue          string   `json:"revenue"`
	Cruchbase_url    string   `json:"crunchbase_url"`
	Description      string   `json:"description"`
	Linkedin_Url     string   `json:"linkedin_url"`
	Phone            string   `json:"phone"`
	Emails           
}

type CompanyResult struct {
	Domain string `json:"domain" binding:"required"`
	Name   string   `json:"name" binding:"required"`
}

type DomainValidation struct {
	Domain string `json:"domain" binding:"required" `
}


type Emails struct {
	Email  string   `json:"email binding:"required"`
}

type EmailFromDomain struct {
	Domain      string   `json:"domain binding:"required"`
}

type EmailValidation struct {
	Email  string  `json:"email binding:"required"`
}

type EmailQuery struct {
	Firstname    *string      `json:"firstname" binding:"required" `
	Lastname     *string      `json:"lastname" binding:"required" `
	Domain       *string      `json:"domain" binding:"required" `    
}


type Profile struct { 
 	ID               primitive.ObjectID  `bson:"_id"`
 	Email            string             `json:"email"`
 	Linkedin_Url     string             `json:"linkedin_url"`
 	Linkedin_Title   string             `json:"linkedin_title"` 
 	EmailQuery
}
