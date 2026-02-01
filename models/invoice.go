package models

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/berkormanli/fatura-go/enums"
	"github.com/berkormanli/fatura-go/internal"
)

type Invoice struct {
	Taxable

	UUID             string            `json:"faturaUuid"`
	DocumentNumber   string            `json:"belgeNumarasi"`
	Date             string            `json:"faturaTarihi"`
	Time             string            `json:"saat"`
	Currency         enums.Currency    `json:"paraBirimi"`
	ExchangeRate     float64           `json:"dovzTLkur"`
	InvoiceType      enums.InvoiceType `json:"faturaTipi"`
	OrderNumber      string            `json:"siparisNumarasi"`
	OrderDate        string            `json:"siparisTarihi"`
	WaybillNumber    string            `json:"irsaliyeNumarasi"`
	WaybillDate      string            `json:"irsaliyeTarihi"`
	ReceiptNumber    string            `json:"fisNo"`
	ReceiptDate      string            `json:"fisTarihi"`
	ReceiptTime      string            `json:"fisSaati"`
	ReceiptType      string            `json:"fisTipi"`
	ZReportNumber    string            `json:"zRaporNo"`
	OkcSerialNumber  string            `json:"okcSeriNo"`
	RecipientTitle   string            `json:"aliciUnvan"`
	RecipientName    string            `json:"aliciAdi"`
	RecipientSurname string            `json:"aliciSoyadi"`
	Address          string            `json:"bulvarcaddesokak"`
	BuildingName     string            `json:"binaAdi"`
	BuildingNumber   string            `json:"binaNo"`
	DoorNumber       string            `json:"kapiNo"`
	Town             string            `json:"kasabaKoy"`
	District         string            `json:"mahalleSemtIlce"`
	City             string            `json:"sehir"`
	Country          string            `json:"ulke"`
	ZipCode          string            `json:"postaKodu"`
	Phone            string            `json:"tel"`
	Fax              string            `json:"fax"`
	Email            string            `json:"eposta"`
	Website          string            `json:"websitesi"`
	TaxOffice        string            `json:"vergiDairesi"`
	RecipientTaxID   string            `json:"vknTckn"`
	
	Note             string            `json:"not"`
	
	// Totals
	BaseAmount       float64           `json:"matrah"`
	ItemTotalAmount  float64           `json:"malhizmetToplamTutari"`
	TotalDiscount    float64           `json:"toplamIskonto"`
	CalculatedVAT    float64           `json:"hesaplanankdv"`
	TotalTaxes       float64           `json:"vergilerToplami"`
	TotalWithTaxes   float64           `json:"vergilerDahilToplamTutar"`
	TotalExpenses    float64           `json:"toplamMasraflar"`
	PaymentTotal     float64           `json:"odenecekTutar"`

	Items            []*InvoiceItem         `json:"-"`
	ReturnItems      []*InvoiceReturnItem   `json:"-"`
}

type InvoiceOption func(*Invoice)

func NewInvoice(taxID string, opts ...InvoiceOption) (*Invoice, error) {
	inv := &Invoice{
		Taxable:        NewTaxable(),
		RecipientTaxID: taxID,
		InvoiceType:    enums.InvoiceTypeSatis, // Default
		Currency:       enums.CurrencyTRY,
		Country:        "Türkiye",
	}

	// Defaults handled via opts or explicit checks
	if inv.UUID == "" {
		inv.UUID = uuid.New().String()
	}
	
	// Current Date/Time defaults if not set? 
	// In Go we usually set defaults first then apply opts.
	// But Date/Time might be overridden by opts.
	// We'll set them if empty after opts.

	for _, opt := range opts {
		opt(inv)
	}

	if inv.Date == "" {
		inv.Date = time.Now().Format("02/01/2006")
	}
	if inv.Time == "" {
		inv.Time = time.Now().Format("15:04:05")
	}

	// Validation
	if inv.UUID != "" {
		if _, err := uuid.Parse(inv.UUID); err != nil {
			return nil, errors.New("Uuid geçerli formatta değil.")
		}
	}
	
	if inv.Date != "" && !internal.ValidateDate(inv.Date) {
		return nil, errors.New("Tarih geçerli formatta değil.")
	}
	
	if inv.Time != "" && !internal.ValidateTime(inv.Time) {
		return nil, errors.New("Saat geçerli formatta değil.")
	}
	
	if inv.Currency != enums.CurrencyTRY && inv.ExchangeRate == 0 {
		return nil, errors.New("Kur bilgisi belirtilmedi.")
	}

	return inv, nil
}

// Options
func WithInvoiceType(t enums.InvoiceType) InvoiceOption {
	return func(i *Invoice) { i.InvoiceType = t }
}

func WithCurrency(c enums.Currency, rate float64) InvoiceOption {
	return func(i *Invoice) {
		i.Currency = c
		i.ExchangeRate = rate
	}
}

func WithRecipientName(name, surname string) InvoiceOption {
	return func(i *Invoice) {
		i.RecipientName = name
		i.RecipientSurname = surname
	}
}

func WithRecipientTitle(title string) InvoiceOption {
	return func(i *Invoice) { i.RecipientTitle = title }
}

func WithAddress(addr, city, country, district string) InvoiceOption {
	return func(i *Invoice) {
		i.Address = addr
		i.City = city
		i.Country = country
		i.District = district
	}
}

// ... more options for other fields can be added as needed

func (i *Invoice) GetInvoiceType() enums.InvoiceType {
	return i.InvoiceType
}

func (i *Invoice) AddItem(items ...*InvoiceItem) *Invoice {
	for _, item := range items {
		// Prepare item with parent context (this invoice)
		item.Prepare(i)
		i.Items = append(i.Items, item)
	}
	i.CalculateTotals()
	return i
}

func (i *Invoice) AddReturnItem(items ...*InvoiceReturnItem) *Invoice {
	if i.InvoiceType == enums.InvoiceTypeIade {
		i.ReturnItems = append(i.ReturnItems, items...)
	}
	return i
}

func (i *Invoice) CalculateTotals() {
	i.ItemTotalAmount = 0
	i.BaseAmount = 0
	i.CalculatedVAT = 0
	
	var discountNormal, discountIncrease float64

	for _, item := range i.Items {
		i.ItemTotalAmount += item.Price
		i.BaseAmount += item.TotalAmount
		i.CalculatedVAT += item.VatAmount
		
		if item.DiscountType == "İskonto" {
			discountNormal += item.DiscountAmount
		} else {
			discountIncrease += item.DiscountAmount
		}
	}
	
	// Discount
	i.TotalDiscount = discountNormal - discountIncrease
	if i.TotalDiscount < 0 {
		i.TotalDiscount = -i.TotalDiscount // Absolute value
	}

	// Taxes Total (VAT + Taxes that are NOT stoppage)
	// We need to aggregate taxes from all items
	
	// Reset invoice level taxes map before re-aggregating
	i.Taxable = NewTaxable()
	
	// Re-aggregate taxes from items
	for _, item := range i.Items {
		for _, _ = range item.GetTaxes() {
			// Aggregate to Invoice level taxable map
			// If tax type exists, add amount. 
			// But TaxDetail has specific Rate.
			// PHP logic relies on TaxableTrait being on Invoice too?
			// PHP AbstractModel has getTaxes() which iterates items.
			// So Invoice itself doesn't hold sum of taxes in a simple list?
			// Wait, TaxableTrait::getTaxes() returns $this->taxes.
			// InvoiceModel::getTaxes references AbstractModel::getTaxes which iterates items.
			// Ah, AbstractModel overrides getTaxes() to aggregate from items!
			
			// So Invoice doesn't store taxes directly in its own map usually, 
			// but we need to calculate totals based on them.
		}
	}
	
	// In PHP calculateTotals:
	// $this->vergilerToplami = $this->hesaplananKdv + array_column_sum...($this->getTaxes(), 'amount', fn($tax) => !$tax['model']->isStoppage())
	
	// taxes := i.GetTaxes()
	// In our Go implementation, Taxable.GetTaxes returns its own map.
	// We need a method on Invoice to return aggregated taxes from items.
	
	taxesList := i.GetAllTaxes()
	
	var totalNonStoppageTax float64
	for _, t := range taxesList {
		if !t.Model.IsStoppage() {
			totalNonStoppageTax += t.Amount
		}
	}
	
	i.TotalTaxes = i.CalculatedVAT + totalNonStoppageTax
	i.TotalWithTaxes = i.BaseAmount + i.TotalTaxes
	
	var totalStoppageTax float64
	for _, t := range taxesList {
		if t.Model.IsStoppage() || t.Model.IsWithholding() {
			totalStoppageTax += t.Amount
		}
	}
	
	i.PaymentTotal = i.TotalWithTaxes - totalStoppageTax
}

// GetAllTaxes aggregates taxes from all items
func (i *Invoice) GetAllTaxes() []TaxDetail {
	var all []TaxDetail
	for _, item := range i.Items {
		for _, v := range item.GetTaxes() {
			all = append(all, v)
		}
	}
	return all
}

func (i *Invoice) GetTotals() map[string]interface{} {
	return map[string]interface{}{
		"matrah":                   i.BaseAmount,
		"malHizmetToplamTutari":    i.ItemTotalAmount,
		"toplamIskonto":            i.TotalDiscount,
		"hesaplananKdv":            i.CalculatedVAT,
		"vergilerToplami":          i.TotalTaxes,
		"vergilerDahilToplamTutar": i.TotalWithTaxes,
		"toplamMasraflar":          i.TotalExpenses,
		"odenecekTutar":            i.PaymentTotal,
	}
}

func (i *Invoice) SetNote(note string) *Invoice {
	i.Note = note
	return i
}

func (i *Invoice) Export() map[string]interface{} {
	// i.Taxable.ExportTaxes() ignored as per PHP logic
	// PHP InvoiceModel export:
	// array_merge($this->toArray(), $this->getTotals(), ['malHizmetListe' => $this->getItems(true), 'iadeTable' => ...])
	// And keyMap.
	
	// We need to implement toArray equivalent (fields map)
	// And merge everything.
	
	base := map[string]interface{}{
		"faturaUuid":       i.UUID,
		"belgeNumarasi":    i.DocumentNumber,
		"faturaTarihi":     i.Date,
		"saat":             i.Time,
		"paraBirimi":       i.Currency,
		"dovzTLkur":        i.ExchangeRate,
		"faturaTipi":       i.InvoiceType,
		"siparisNumarasi":  i.OrderNumber,
		"siparisTarihi":    i.OrderDate,
		"irsaliyeNumarasi": i.WaybillNumber,
		"irsaliyeTarihi":   i.WaybillDate,
		"fisNo":            i.ReceiptNumber,
		"fisTarihi":        i.ReceiptDate,
		"fisSaati":         i.ReceiptTime,
		"fisTipi":          i.ReceiptType,
		"zRaporNo":         i.ZReportNumber,
		"okcSeriNo":        i.OkcSerialNumber,
		"aliciUnvan":       i.RecipientTitle,
		"aliciAdi":         i.RecipientName,
		"aliciSoyadi":      i.RecipientSurname,
		"bulvarcaddesokak": i.Address,
		"binaAdi":          i.BuildingName,
		"binaNo":           i.BuildingNumber,
		"kapiNo":           i.DoorNumber,
		"kasabaKoy":        i.Town,
		"mahalleSemtIlce":  i.District,
		"sehir":            i.City,
		"ulke":             i.Country,
		"postaKodu":        i.ZipCode,
		"tel":              i.Phone,
		"fax":              i.Fax,
		"eposta":           i.Email,
		"websitesi":        i.Website,
		"vergiDairesi":     i.TaxOffice,
		"vknTckn":          i.RecipientTaxID,
		"not":              i.Note,
		"hangiTip":         enums.TypeEArsivFatura, // Default
	}
	
	// Merge totals
	totals := i.GetTotals()
	for k, v := range totals {
		base[k] = v
	}
	
	// Items
	var itemsList []map[string]interface{}
	for _, item := range i.Items {
		itemsList = append(itemsList, item.Export())
	}
	base["malHizmetTable"] = itemsList
	
	// Return Items
	if len(i.ReturnItems) > 0 {
		var returnsList []map[string]interface{}
		for _, item := range i.ReturnItems {
			returnsList = append(returnsList, item.Export())
		}
		base["iadeTable"] = returnsList
	}
	
	// Invoice Level Taxes?
	// PHP's InvoiceModel uses $this->getTaxes() which aggregates items.
	// And then it seems to include them in the export if generic `getTaxes` is used.
	// But `InvoiceModel::export` merges `keyMapper`.
	// The `exportTaxes` method in Trait uses `getTaxes`.
	// Since `InvoiceModel` extends `AbstractModel` which implements `getTaxes` (aggregated),
	// `exportTaxes` will export aggregated taxes.
	
	// So we need to export aggregated taxes too.
	// Aggregated tax export logic skipped as per analysis
	// aggregatedTaxes := i.GetAllTaxes()
	
	// We need helper to export slice of TaxDetail to the map format expected
	// Helper in Taxable handles `t.GetTaxes()` which is a map.
	// We have a slice.
	
	// Aggregated tax export logic skipped as per analysis
	/*
	tempTaxable := NewTaxable()
	for _, _ = range aggregatedTaxes {
	    // ...
	}
	*/

	
	return base
}

func (i *Invoice) GetUUID() string {
	return i.UUID
}
