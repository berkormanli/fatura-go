package models

import (
	"github.com/berkormanli/fatura-go/enums"
	"github.com/berkormanli/fatura-go/internal"
)

type ProducerReceiptItem struct {
	Taxable

	Name            string     `json:"malHizmet"`
	Quantity        float64    `json:"miktar"`
	UnitPrice       float64    `json:"birimFiyat"`
	Unit            enums.Unit `json:"birim"`
	TotalAmount     float64    `json:"malHizmetTutari"`
	GvStoppageRate  int        `json:"gvStopajOrani"`

	imported bool
}

type ProducerItemOption func(*ProducerReceiptItem)

func NewProducerReceiptItem(name string, qty float64, unitPrice float64, gvStopajRate int, opts ...ProducerItemOption) *ProducerReceiptItem {
	item := &ProducerReceiptItem{
		Taxable:        NewTaxable(),
		Name:           name,
		Quantity:       qty,
		UnitPrice:      unitPrice,
		Unit:           enums.UnitAdet,
		GvStoppageRate: gvStopajRate,
	}

	for _, opt := range opts {
		opt(item)
	}

	if !item.imported {
		if item.TotalAmount == 0 {
			item.TotalAmount = item.Quantity * item.UnitPrice
		}
		item.AddTax(enums.TaxGVStopaj, item.GvStoppageRate, 0)
	}

	return item
}

// Options
func WithProducerUnit(u enums.Unit) ProducerItemOption {
	return func(i *ProducerReceiptItem) { i.Unit = u }
}

func (i *ProducerReceiptItem) AddTax(tax enums.Tax, rate int, amount float64) *ProducerReceiptItem {
	if amount == 0 {
		if tax == enums.TaxBorsaTescil {
			// BorsaTescil calculation: (Total - TaxesExceptBorsaTescil) * Rate
			// But current taxes map might not be full? 
			// PHP code: $this->malHizmetTutari - $this->totalTaxAmount(fn ($tax) => $tax['model'] != Tax::BorsaTescil)
			// This implies it subtracts ALREADY ADDED taxes?
			// Wait, TaxGVStopaj is added in constructor.
			// So if we add BorsaTescil after constructor, it subtracts GVStopaj?
			// Yes, Producer Receipts usually subtract Stoppage from base?
			// Or is Borsa Tescil calculated on net? 
			// Let's stick to PHP logic: subtract other taxes from base.
			
			otherTaxes := i.TotalTaxAmount(func(t enums.Tax) bool {
				return t != enums.TaxBorsaTescil
			})
			amount = internal.Percentage(i.TotalAmount-otherTaxes, float64(rate))
		} else {
			amount = internal.Percentage(i.TotalAmount, float64(rate))
		}
	}

	i.Taxable.AddTax(tax, rate, amount, 0)
	return i
}

func (i *ProducerReceiptItem) Prepare(parent Model) ItemModel {
	// Re-calculate taxes? PHP calculateTaxes uses array_map with callables.
	// We handle calculation in AddTax.
	return i
}

func (i *ProducerReceiptItem) GetTotals() map[string]interface{} {
	return map[string]interface{}{
		"birimFiyat":      i.UnitPrice,
		"malHizmetTutari": i.TotalAmount,
	}
}

func (i *ProducerReceiptItem) Export() map[string]interface{} {
	// PHP exports taxes with lowerFirst=true for Producer Receipts!
	// We need to implement lowerFirst logic in Taxable or here.
	
	taxes := i.Taxable.ExportTaxes()
	lowerTaxes := make(map[string]interface{})
	for k, v := range taxes {
		// Convert Key to lowercase first letter
		// k is like V0021Orani -> v0021Orani
		// Or using regex/string manipulation
		if len(k) > 0 {
			lowerKey := string(internal.ToLower(rune(k[0]))) + k[1:]
			lowerTaxes[lowerKey] = v
		}
	}

	m := map[string]interface{}{
		"malHizmet":       i.Name,
		"miktar":          i.Quantity,
		"birim":           i.Unit,
		"birimFiyat":      i.UnitPrice,
		"malHizmetTutari": i.TotalAmount,
		"gvStopajOrani":   i.GvStoppageRate,
	}
	
	for k, v := range lowerTaxes {
		m[k] = v
	}
	
	return m
}
