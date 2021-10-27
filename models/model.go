package models 

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type CompanyQuery struct {
	Name   string   `json:"name"`
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

type DomainValidation struct {
	Domain string `json:"domain"`
}


type Emails struct {
	Email  string   `json:"email"`
}

type EmailFromDomain struct {
	Domain      string   `json:"domain"`
}

type EmailValidation struct {
	Email  string  `json:"email"`
}

type EmailQuery struct {
	Firstname    *string      `json:"firstname"`
	Lastname     *string      `json:"lastname"`
	Domain       *string      `json:"domain"`    
}


type Profile struct { 
 	ID               primitive.ObjectID  `bson:"_id"`
 	Email            string             `json:"email"`
 	Linkedin_Url     string             `json:"linkedin_url"`
 	Linkedin_Title   string             `json:"linkedin_title"` 
 	EmailQuery
}
