package models

import (
	"errors"
	"github.com/berkormanli/fatura-go/enums"
	"github.com/berkormanli/fatura-go/internal"
)

type InvoiceItem struct {
	Taxable
	
	Name             string         `json:"malHizmet"`
	Quantity         float64        `json:"miktar"`
	UnitPrice        float64        `json:"birimFiyat"`
	VatRate          float64        `json:"kdvOrani"`
	Unit             enums.Unit     `json:"birim"`
	Price            float64        `json:"fiyat"`
	DiscountType     string         `json:"iskontoTipi"`
	DiscountRate     float64        `json:"iskontoOrani"`
	DiscountAmount   float64        `json:"iskontoTutari"`
	DiscountReason   string         `json:"iskontoNedeni"`
	TotalAmount      float64        `json:"malHizmetTutari"`
	VatAmount        float64        `json:"kdvTutari"`
	WithholdingCode  int            `json:"tevkifatKodu"`
	SpecialBaseReason int           `json:"ozelMatrahNedeni"`
	SpecialBaseAmount float64       `json:"ozelMatrahTutari"`
	Gtip             string         `json:"gtip"`

	imported bool
}

type InvoiceItemOption func(*InvoiceItem)

func NewInvoiceItem(name string, qty float64, unitPrice float64, vatRate float64, opts ...InvoiceItemOption) (*InvoiceItem, error) {
	item := &InvoiceItem{
		Taxable:      NewTaxable(),
		Name:         name,
		Quantity:     qty,
		UnitPrice:    unitPrice,
		VatRate:      vatRate,
		Unit:         enums.UnitAdet,
		DiscountType: "İskonto",
	}

	for _, opt := range opts {
		opt(item)
	}

	// Validation
	validVat := false
	for _, r := range []float64{0, 1, 8, 10, 18, 20} {
		if item.VatRate == r {
			validVat = true
			break
		}
	}
	if !validVat {
		return nil, errors.New("Invalid VAT rate")
	}

	if item.DiscountType != "İskonto" && item.DiscountType != "Arttırım" {
		return nil, errors.New("Invalid discount type")
	}

	if !item.imported {
		if item.Price == 0 {
			item.Price = item.Quantity * item.UnitPrice
		}

		if item.DiscountRate > 0 && item.DiscountAmount == 0 {
			item.DiscountAmount = internal.Percentage(item.Price, item.DiscountRate)
		}

		if item.TotalAmount == 0 {
			if item.DiscountAmount == 0 {
				item.TotalAmount = item.Price
			} else {
				if item.DiscountType == "İskonto" {
					item.TotalAmount = item.Price - item.DiscountAmount
				} else {
					item.TotalAmount = item.Price + item.DiscountAmount
				}
			}
		}

		if item.VatAmount == 0 {
			item.VatAmount = internal.Percentage(item.TotalAmount, item.VatRate)
		}
	}

	return item, nil
}

// Options
func WithUnit(u enums.Unit) InvoiceItemOption {
	return func(i *InvoiceItem) { i.Unit = u }
}

func WithDiscount(typ string, rate float64, reason string) InvoiceItemOption {
	return func(i *InvoiceItem) {
		i.DiscountType = typ
		i.DiscountRate = rate
		i.DiscountReason = reason
	}
}

func WithGtip(gtip string) InvoiceItemOption {
	return func(i *InvoiceItem) { i.Gtip = gtip }
}

func WithWithholding(code int) InvoiceItemOption {
	return func(i *InvoiceItem) { i.WithholdingCode = code }
}

func WithSpecialBase(reason int, amount float64) InvoiceItemOption {
	return func(i *InvoiceItem) { 
		i.SpecialBaseReason = reason
		i.SpecialBaseAmount = amount
	}
}

func (i *InvoiceItem) AddTax(tax enums.Tax, rate int, amount float64, vat float64) *InvoiceItem {
	if amount == 0 {
		if tax == enums.TaxKDVTevkifat {
			amount = internal.Percentage(i.VatAmount, float64(rate))
		} else {
			amount = internal.Percentage(i.TotalAmount, float64(rate))
		}
	}

	if tax == enums.TaxOTV1ListeTevkifat {
		amount *= i.Quantity
	}

	if vat == 0 && tax.HasVat() {
		vat = internal.Percentage(amount, i.VatRate)
	}

	i.Taxable.AddTax(tax, rate, amount, vat)
	return i
}

func (i *InvoiceItem) Prepare(parent Model) ItemModel {
	// Need to access parent invoice type
	// Assuming parent is *Invoice or exposes it.
	// We'll define an interface or type assertion.
	
	// Since we haven't defined Invoice yet, we can't assert *Invoice.
	// But we can check if parent implements a method GetInvoiceType
	type InvoiceTyper interface {
		GetInvoiceType() enums.InvoiceType
	}
	
	var invType enums.InvoiceType = enums.InvoiceTypeSatis
	if p, ok := parent.(InvoiceTyper); ok {
		invType = p.GetInvoiceType()
	}

	i.VatAmount += i.TotalTaxVat(nil)

	if invType == enums.InvoiceTypeTevkifat && i.WithholdingCode > 0 {
		// Tevkifat logic
		// Check code validity?
		if info, ok := enums.TaxKDVTevkifat.GetRate(i.WithholdingCode); ok {
			i.AddTax(enums.TaxKDVTevkifat, info, 0, 0)
		}
	}

	if invType == enums.InvoiceTypeOzelMatrah && i.SpecialBaseReason > 0 {
		// Ozel Matrah Logic
		// Check reason?
		i.VatAmount = internal.Percentage(i.SpecialBaseAmount, i.VatRate) + i.TotalTaxVat(nil)
	}
	
	// Calculate total taxes internally update the map
	// The AddTax method handles addition to the map.
	
	return i
}

func (i *InvoiceItem) GetTotals() map[string]interface{} {
	return map[string]interface{}{
		"birimFiyat":       i.UnitPrice,
		"fiyat":            i.Price,
		"iskontoTutari":    i.DiscountAmount,
		"malHizmetTutari":  i.TotalAmount,
		"kdvTutari":        i.VatAmount,
		"ozelMatrahTutari": i.SpecialBaseAmount,
	}
}

func (i *InvoiceItem) Export() map[string]interface{} {
	m := map[string]interface{}{
		"malHizmet":        i.Name,
		"miktar":           i.Quantity,
		"birim":            i.Unit,
		"birimFiyat":       i.UnitPrice,
		"kdvOrani":         i.VatRate,
		"iskontoTipi":      i.DiscountType,
		"iskontoOrani":     i.DiscountRate,
		"iskontoTutari":    i.DiscountAmount,
		"iskontoNedeni":    i.DiscountReason,
		"malHizmetTutari":  i.TotalAmount,
		"kdvTutari":        i.VatAmount,
		"tevkifatKodu":     i.WithholdingCode,
		"ozelMatrahNedeni": i.SpecialBaseReason,
		"ozelMatrahTutari": i.SpecialBaseAmount,
		"gtip":             i.Gtip,
		
		// Mapped keys
		"iskontoArttm":     i.DiscountType,
		
		// Totals
		"fiyat":            i.Price,
	}
	
	// Merge taxes
	for k, v := range i.ExportTaxes() {
		m[k] = v
	}
	
	return m
}
