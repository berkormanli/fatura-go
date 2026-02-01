package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/berkormanli/fatura-go"
	"github.com/berkormanli/fatura-go/enums"
	"github.com/berkormanli/fatura-go/models"
)

// Run with: go test ./tests -v

func TestServiceAuth(t *testing.T) {
	service := fatura.NewGib(enums.DocumentTypeInvoice)
	service.SetCredentials("333333", "666666")
	
	u, p := service.GetCredentials()
	assert.Equal(t, "333333", u)
	assert.Equal(t, "666666", p)
	
	// Test Credentials (network call)
	// We might skip this if network is flaky, but PHP test does it.
	service.SetTestCredentials("", "")
	uTest, _ := service.GetCredentials()
	assert.NotEmpty(t, uTest)
}

func TestGetAll(t *testing.T) {
	service := fatura.NewGib(enums.DocumentTypeInvoice)
	err := service.SetTestCredentials("", "").Login("", "")
	if err != nil {
		t.Skip("Login failed, skipping network test")
	}
	
	now := time.Now().Format("02/01/2006")
	prev := time.Now().AddDate(0, -1, 0).Format("02/01/2006")
	
	docs, err := service.GetAll(prev, now)
	assert.Nil(t, err)
	assert.IsType(t, []map[string]interface{}{}, docs)
	
	if len(docs) > 0 {
		doc := docs[0]
		assert.Contains(t, doc, "belgeNumarasi")
		assert.Contains(t, doc, "ettn")
	}
}

func TestCreateDraft(t *testing.T) {
	service := fatura.NewGib(enums.DocumentTypeInvoice)
	err := service.SetTestCredentials("", "").Login("", "")
	if err != nil {
		t.Skip("Login failed")
	}

	inv, _ := models.NewInvoice("11111111111", 
		models.WithRecipientName("Mert", "Levent"),
		models.WithAddress("Papatya sk.", "Bursa", "Türkiye", "Nilüfer"),
	)
	
	item, _ := models.NewInvoiceItem("Test Item", 1, 100, 18)
	inv.AddItem(item)
	
	err = service.CreateDraft(inv)
	assert.Nil(t, err)
	
	// Get Document
	created, err := service.GetDocument(inv.GetUUID())
	assert.Nil(t, err)
	assert.Equal(t, "Mert", created["aliciAdi"])
	
	// Delete
	err = service.DeleteDraft([]string{inv.GetUUID()}, "Test Delete")
	assert.Nil(t, err)
}
