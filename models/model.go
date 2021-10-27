package models 

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)


type CompanyName struct {
	Name   string   `json:"name"`
}

type DomainValidation struct {
	Domain string `json:"domain"`
}

type EmailValidation struct {
	Email  string  `json:"email"`
}

type EmailQuery struct {
	Firstname    *string      `json:"firstname"`
	Lastname     *string      `json:"lastname"`
	Domain       *string      `json:"domain"`    
}

type EmailFromDomain struct {
	Domain      string   `json:"domain"`
}

type Profile struct { 
 	ID               primitive.ObjectID  `bson:"_id"`
 	Email            string             `json:"email"`
 	Linkedin_Url     string             `json:"linkedin_url"`
 	Linkedin_Title   string             `json:"linkedin_title"` 
 	EmailQuery
}
