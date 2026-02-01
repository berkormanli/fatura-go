package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/berkormanli/fatura-go/enums"

)

type SelfEmployedReceipt struct {
	UUID             string                     `json:"ettn"` // KeyMapped to ettn
	DocumentNumber   string                     `json:"belgeNumarasi"`
	Date             string                     `json:"tarih"`
	Time             string                     `json:"saat"`
	Currency         enums.Currency             `json:"paraBirimi"`
	ExchangeRate     float64                    `json:"dovizKuru"` // KeyMapped to kur
	RecipientTaxID   string                     `json:"vknTckn"`
	RecipientName    string                     `json:"aliciAdi"` // KeyMapped to adi
	RecipientSurname string                     `json:"aliciSoyadi"` // KeyMapped to soyadi
	RecipientTitle   string                     `json:"aliciUnvan"` // KeyMapped to unvan
	Address          string                     `json:"adres"` // KeyMapped to bulvarCaddeSokak
	BuildingName     string                     `json:"binaAdi"`
	BuildingNumber   string                     `json:"binaNo"`
	DoorNumber       string                     `json:"kapiNo"`
	Town             string                     `json:"kasabaKoy"`
	District         string                     `json:"mahalleSemtIlce"`
	City             string                     `json:"sehir"`
	Country          string                     `json:"ulke"`
	ZipCode          string                     `json:"postaKodu"`
	TaxOffice        string                     `json:"vergiDairesi"`
	Description      string                     `json:"aciklama"`
	KdvAccrual       bool                       `json:"kdvTahakkukIcin"`
	
	// Totals
	GrossWageTotal         float64 `json:"brutUcret"` // KeyMapped to brtUcret
	GvStoppageTotal        float64 `json:"gvStopajTutari"` // KeyMapped to gvStpjTtari
	NetWageTotal           float64 `json:"netUcretTutari"` // KeyMapped to netUcretTtr
	VatTotal               float64 `json:"kdvTutari"` // KeyMapped to kdvTtri
	KdvWithholdingTotal    float64 `json:"kdvTevkifatTutari"` // KeyMapped to kdvTvkftTtri
	CollectedVatTotal      float64 `json:"tahsilEdilenKdv"` // KeyMapped to thsilEdilenKdv
	NetReceivedTotal       float64 `json:"netAlinanToplam"`

	Items            []*SelfEmployedReceiptItem `json:"-"`
}

type SelfEmployedOption func(*SelfEmployedReceipt)

func NewSelfEmployedReceipt(taxID, name, surname string, opts ...SelfEmployedOption) (*SelfEmployedReceipt, error) {
	ser := &SelfEmployedReceipt{
		RecipientTaxID:   taxID,
		RecipientName:    name,
		RecipientSurname: surname,
		Country:          "Türkiye",
		Currency:         enums.CurrencyTRY,
	}

	if ser.UUID == "" {
		ser.UUID = uuid.New().String()
	}

	for _, opt := range opts {
		opt(ser)
	}

	if ser.Date == "" {
		ser.Date = time.Now().Format("02/01/2006")
	}
	if ser.Time == "" {
		ser.Time = time.Now().Format("15:04:05")
	}

	if ser.Currency != enums.CurrencyTRY && ser.ExchangeRate == 0 {
		return nil, errors.New("Kur bilgisi belirtilmedi.")
	}

	return ser, nil
}

// Options
func WithSelfEmployedAddress(addr, city, country, district string) SelfEmployedOption {
	return func(s *SelfEmployedReceipt) {
		s.Address = addr
		s.City = city
		s.Country = country
		s.District = district
	}
}

func WithSelfEmployedCurrency(c enums.Currency, rate float64) SelfEmployedOption {
	return func(s *SelfEmployedReceipt) {
		s.Currency = c
		s.ExchangeRate = rate
	}
}

func (s *SelfEmployedReceipt) AddItem(items ...*SelfEmployedReceiptItem) *SelfEmployedReceipt {
	for _, item := range items {
		item.Prepare(s)
		s.Items = append(s.Items, item)
	}
	s.CalculateTotals()
	return s
}

func (s *SelfEmployedReceipt) CalculateTotals() {
	s.GrossWageTotal = 0
	s.GvStoppageTotal = 0
	s.NetWageTotal = 0
	s.VatTotal = 0
	s.KdvWithholdingTotal = 0
	s.NetReceivedTotal = 0
	
	for _, item := range s.Items {
		s.GrossWageTotal += item.GrossWage
		s.GvStoppageTotal += item.GvStoppageAmount
		s.NetWageTotal += item.NetWage
		s.VatTotal += item.VatAmount
		s.KdvWithholdingTotal += item.KdvWithholdingAmount
		s.NetReceivedTotal += item.NetReceived
	}
	
	s.CollectedVatTotal = s.VatTotal - s.KdvWithholdingTotal
}

func (s *SelfEmployedReceipt) GetTotals() map[string]interface{} {
	return map[string]interface{}{
		"brutUcret":         s.GrossWageTotal,
		"gvStopajTutari":    s.GvStoppageTotal,
		"netUcretTutari":    s.NetWageTotal,
		"kdvTutari":         s.VatTotal,
		"kdvTevkifatTutari": s.KdvWithholdingTotal,
		"tahsilEdilenKdv":   s.CollectedVatTotal,
		"netAlinanToplam":   s.NetReceivedTotal,
	}
}

func (s *SelfEmployedReceipt) Export() map[string]interface{} {
	base := map[string]interface{}{
		"ettn":              s.UUID,
		"belgeNumarasi":     s.DocumentNumber,
		"tarih":             s.Date,
		"saat":              s.Time,
		"paraBirimi":        s.Currency,
		"kur":               s.ExchangeRate,
		"vknTckn":           s.RecipientTaxID,
		"adi":               s.RecipientName,
		"soyadi":            s.RecipientSurname,
		"unvan":             s.RecipientTitle,
		"bulvarCaddeSokak":  s.Address,
		"binaAdi":           s.BuildingName,
		"binaNo":            s.BuildingNumber,
		"kapiNo":            s.DoorNumber,
		"kasabaKoy":         s.Town,
		"mahalleSemtIlce":   s.District,
		"sehir":             s.City,
		"ulke":              s.Country,
		"postaKodu":         s.ZipCode,
		"vergiDairesi":      s.TaxOffice,
		"aciklama":          s.Description,
		"kdvTahakkukIcin":   s.KdvAccrual,
		
		// Map totals
		"brtUcret":          s.GrossWageTotal,
		"gvStpjTtari":       s.GvStoppageTotal,
		"netUcretTtr":       s.NetWageTotal,
		"kdvTtri":           s.VatTotal,
		"kdvTvkftTtri":      s.KdvWithholdingTotal,
		"thsilEdilenKdv":    s.CollectedVatTotal,
		"netAlinanToplam":   s.NetReceivedTotal,
	}
	
	var itemsList []map[string]interface{}
	for _, item := range s.Items {
		itemsList = append(itemsList, item.Export())
	}
	base["serbestTable"] = itemsList
	
	return base
}

func (s *SelfEmployedReceipt) GetUUID() string {
	return s.UUID
}
