package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/berkormanli/fatura-go/internal"
)

type ProducerReceipt struct {
	UUID             string                 `json:"faturaUuid"`
	DocumentNumber   string                 `json:"belgeNumarasi"`
	Date             string                 `json:"faturaTarihi"` // Note: JSON key faturaTarihi in PHP Base Model for 'tarih'
	Time             string                 `json:"saat"`
	RecipientTaxID   string                 `json:"vknTckn"`
	RecipientName    string                 `json:"aliciAdi"`
	RecipientSurname string                 `json:"aliciSoyadi"`
	City             string                 `json:"sehir"`
	Website          string                 `json:"websitesi"`
	Note             string                 `json:"not"`
	DeliveryDate     string                 `json:"teslimTarihi"`
	
	ItemTotalAmount  float64                `json:"malHizmetToplamTutari"`
	TotalWithTaxes   float64                `json:"vergilerDahilToplamTutar"`
	PaymentTotal     float64                `json:"odenecekTutar"`

	Items            []*ProducerReceiptItem `json:"-"`
}

type ProducerOption func(*ProducerReceipt)

func NewProducerReceipt(taxID, name, surname string, opts ...ProducerOption) (*ProducerReceipt, error) {
	pr := &ProducerReceipt{
		RecipientTaxID:   taxID,
		RecipientName:    name,
		RecipientSurname: surname,
	}

	if pr.UUID == "" {
		pr.UUID = uuid.New().String()
	}

	for _, opt := range opts {
		opt(pr)
	}

	if pr.Date == "" {
		pr.Date = time.Now().Format("02/01/2006")
	}
	if pr.Time == "" {
		pr.Time = time.Now().Format("15:04:05")
	}

	if pr.DeliveryDate != "" && !internal.ValidateDate(pr.DeliveryDate) {
		return nil, errors.New("Teslim tarihi geçerli formatta değil.")
	}

	return pr, nil
}

// Options
func WithProducerDate(date, time string) ProducerOption {
	return func(p *ProducerReceipt) {
		p.Date = date
		p.Time = time
	}
}

func WithDeliveryDate(date string) ProducerOption {
	return func(p *ProducerReceipt) { p.DeliveryDate = date }
}

func WithProducerCity(city string) ProducerOption {
	return func(p *ProducerReceipt) { p.City = city }
}

func (p *ProducerReceipt) AddItem(items ...*ProducerReceiptItem) *ProducerReceipt {
	for _, item := range items {
		item.Prepare(p)
		p.Items = append(p.Items, item)
	}
	p.CalculateTotals()
	return p
}

func (p *ProducerReceipt) CalculateTotals() {
	p.ItemTotalAmount = 0
	var totalTaxes float64

	for _, item := range p.Items {
		p.ItemTotalAmount += item.TotalAmount
		// In Producer Receipt:
		// Vergiler dahil toplam = malHizmetToplamTutari (Base Amount is Gross?)
		// Odenecek = Base - Taxes (Stoppage)
		
		// PHP: 
		// odenecek = vergilerDahil - taxes->amount.
		// items taxes are aggregated.
		
		itemTaxes := item.TotalTaxAmount(nil)
		totalTaxes += itemTaxes
	}
	
	p.TotalWithTaxes = p.ItemTotalAmount
	p.PaymentTotal = p.TotalWithTaxes - totalTaxes
}

func (p *ProducerReceipt) GetTotals() map[string]interface{} {
	return map[string]interface{}{
		"malHizmetToplamTutari":    p.ItemTotalAmount,
		"vergilerDahilToplamTutar": p.TotalWithTaxes,
		"odenecekTutar":            p.PaymentTotal,
	}
}

func (p *ProducerReceipt) SetNote(note string) *ProducerReceipt {
	p.Note = note
	return p
}

func (p *ProducerReceipt) Export() map[string]interface{} {
	base := map[string]interface{}{
		"faturaUuid":    p.UUID,
		"belgeNumarasi": p.DocumentNumber,
		"faturaTarihi":  p.Date, // PHP AbstractModel uses 'tarih' mapped to 'faturaTarihi' in Invoice?
		// Wait, ProducerReceiptModel DOES NOT implement keyMap for 'tarih'.
		// But AbstractModel defines 'tarih'.
		// InvoiceModel maps 'tarih' -> 'faturaTarihi'.
		// ProducerReceiptModel does NOT map it?
		// Let's check ProducerReceiptModel.php keyMap: 'teslimTarihi'->'teslimTarih', 'malHizmetListe'->'mustahsilTable'.
		// It does NOT map 'tarih'. So it stays 'tarih'?
		// Actually AbstractModel does NOT hold 'faturaTarihi'.
		// The API likely expects 'tarih'?
		// Let's check step 128 (ProducerReceiptModel.php).
		// It does not map 'tarih'.
		// BUT `InvoiceModel` DOES map 'tarih' -> 'faturaTarihi'.
		// So Producer might use 'tarih'.
		// However, I used 'faturaTarihi' in struct tag.
		// I should check strict API requirements from PHP usage or assumed knowledge.
		// Usually GIB uses `faturaTarihi` for Invoices and `belgeTarihi` or `tarih` for others?
		// Note 639 in README output shows: `[belgeTarihi] => 09-10-2022`.
		// But that's response.
		// Request: PHP `AbstractModel` has public `$tarih`.
		// `InvoiceModel` keyMap: `'tarih' => 'faturaTarihi'`.
		// `ProducerReceiptModel` keyMap: NO 'tarih' mapping.
		// So it sends "tarih".
	}
	
	// Fix keys based on analysis
	base["tarih"] = p.Date
	base["saat"] = p.Time
	base["vknTckn"] = p.RecipientTaxID
	base["aliciAdi"] = p.RecipientName
	base["aliciSoyadi"] = p.RecipientSurname
	base["sehir"] = p.City
	base["websitesi"] = p.Website
	base["not"] = p.Note
	base["teslimTarih"] = p.DeliveryDate
	
	totals := p.GetTotals()
	for k, v := range totals {
		base[k] = v
	}
	
	var itemsList []map[string]interface{}
	for _, item := range p.Items {
		itemsList = append(itemsList, item.Export())
	}
	base["mustahsilTable"] = itemsList
	
	return base
}

func (p *ProducerReceipt) GetUUID() string {
	return p.UUID
}
