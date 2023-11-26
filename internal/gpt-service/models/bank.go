package models

import (
	"github.com/search-platform/gpt-service/api/gpt"
	"github.com/search-platform/gpt-service/internal/pkg/modelhelper"
)

type Contact struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Value       string `json:"value"`
}

type BankDetails struct {
	URL         string     `json:"url"`
	Name        string     `json:"name"`
	Country     string     `json:"country"`
	LogoLink    string     `json:"logo_link"`
	FaviconLink string     `json:"favicon_link"`
	Address     string     `json:"address"`
	Contacts    []*Contact `json:"contacts"`
}

func (c *Contact) ToAPI() *gpt.Contact {
	contact := &gpt.Contact{
		Value:       c.Value,
		Description: c.Description,
	}
	if c.Type == "PHONE" {
		contact.Type = gpt.Contact_PHONE
	}
	if c.Type == "EMAIL" {
		contact.Type = gpt.Contact_EMAIL
	}
	return contact
}

func (b *BankDetails) ToAPI() *gpt.BankInfo {
	bankDetails := &gpt.BankInfo{
		Url:         b.URL,
		Name:        b.Name,
		Country:     b.Country,
		LogoLink:    b.LogoLink,
		FaviconLink: b.FaviconLink,
		Address:     b.Address,
	}
	if len(b.Contacts) > 0 {
		bankDetails.Contacts = modelhelper.APISlicePtr[*gpt.Contact](b.Contacts)
	}
	return bankDetails
}
