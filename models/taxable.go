package models

import (
	"fmt"
	"unicode"

	"github.com/berkormanli/fatura-go/enums"
)

type TaxDetail struct {
	Model  enums.Tax
	Rate   int
	Amount float64
	Vat    float64
}

type Taxable struct {
	Taxes map[enums.Tax]TaxDetail
}

func NewTaxable() Taxable {
	return Taxable{
		Taxes: make(map[enums.Tax]TaxDetail),
	}
}

func (t *Taxable) AddTax(tax enums.Tax, rate int, amount float64, vat float64) *Taxable {
	if t.Taxes == nil {
		t.Taxes = make(map[enums.Tax]TaxDetail)
	}
	t.Taxes[tax] = TaxDetail{
		Model:  tax,
		Rate:   rate,
		Amount: amount,
		Vat:    vat,
	}
	return t
}

func (t *Taxable) GetTaxes() map[enums.Tax]TaxDetail {
	if t.Taxes == nil {
		return map[enums.Tax]TaxDetail{}
	}
	return t.Taxes
}

func (t *Taxable) TotalTaxAmount(filterFn func(enums.Tax) bool) float64 {
	var total float64
	for _, tax := range t.GetTaxes() {
		if filterFn == nil || filterFn(tax.Model) {
			total += tax.Amount
		}
	}
	return total
}

func (t *Taxable) TotalTaxVat(filterFn func(enums.Tax) bool) float64 {
	var total float64
	for _, tax := range t.GetTaxes() {
		if filterFn == nil || filterFn(tax.Model) {
			total += tax.Vat
		}
	}
	return total
}

func (t *Taxable) ExportTaxes() map[string]interface{} {
	taxes := make(map[string]interface{})
	for _, taxDetail := range t.GetTaxes() {
		// PHP: V0021Orani, V0021Tutari, V0021KdvTutari
		// Tax Value usually numbers e.g. "0021"
		
		keyRate := fmt.Sprintf("V%sOrani", taxDetail.Model)
		keyAmount := fmt.Sprintf("V%sTutari", taxDetail.Model)
		
		taxes[keyRate] = taxDetail.Rate
		taxes[keyAmount] = taxDetail.Amount
		
		if taxDetail.Model.HasVat() {
			keyVat := fmt.Sprintf("V%sKdvTutari", taxDetail.Model)
			taxes[keyVat] = taxDetail.Vat
		}
	}
	return taxes
}

// Helper to lower first char (though PHP code used it optionally, mostly not needed for standard export)
func lcfirst(s string) string {
	if len(s) == 0 {
		return s
	}
	r := []rune(s)
	r[0] = unicode.ToLower(r[0])
	return string(r)
}
