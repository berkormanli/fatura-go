package models

import (
	"github.com/berkormanli/fatura-go/internal"
)

type SelfEmployedReceiptItem struct {
	Reason              string  `json:"neIcinAlindigi"`
	GrossWage           float64 `json:"brutUcret"`
	VatRate             float64 `json:"kdvOrani"`
	GvStoppageRate      int     `json:"gvStopajOrani"`
	NetWage             float64 `json:"netUcret"`
	KdvWithholdingRate  int     `json:"kdvTevkifatOrani"`
	NetReceived         float64 `json:"netAlinan"`

	GvStoppageAmount    float64 `json:"gvStopajTutari"`
	VatAmount           float64 `json:"kdvTutari"`
	KdvWithholdingAmount float64 `json:"kdvTevkifatTutari"`
	
	imported bool
}

type SelfEmployedItemOption func(*SelfEmployedReceiptItem)

func NewSelfEmployedReceiptItem(reason string, grossWage float64, vatRate float64, opts ...SelfEmployedItemOption) (*SelfEmployedReceiptItem, error) {
	item := &SelfEmployedReceiptItem{
		Reason:    reason,
		GrossWage: grossWage,
		VatRate:   vatRate,
	}

	for _, opt := range opts {
		opt(item)
	}

	if !item.imported {
		if item.GvStoppageAmount == 0 && item.GvStoppageRate > 0 {
			item.GvStoppageAmount = internal.Percentage(item.GrossWage, float64(item.GvStoppageRate))
		}

		if item.NetWage == 0 {
			item.NetWage = item.GrossWage - item.GvStoppageAmount
		}

		if item.VatAmount == 0 {
			item.VatAmount = internal.Percentage(item.GrossWage, item.VatRate)
		}

		if item.KdvWithholdingAmount == 0 && item.KdvWithholdingRate > 0 {
			item.KdvWithholdingAmount = internal.Percentage(item.VatAmount, float64(item.KdvWithholdingRate))
		}

		if item.NetReceived == 0 {
			item.NetReceived = item.NetWage + item.VatAmount - item.KdvWithholdingAmount
		}
	}

	return item, nil
}

// Options
func WithGvStoppage(rate int) SelfEmployedItemOption {
	return func(i *SelfEmployedReceiptItem) { i.GvStoppageRate = rate }
}

func WithKdvWithholding(rate int) SelfEmployedItemOption {
	return func(i *SelfEmployedReceiptItem) { i.KdvWithholdingRate = rate }
}

func (i *SelfEmployedReceiptItem) Prepare(parent Model) ItemModel {
	return i
}

func (i *SelfEmployedReceiptItem) GetTotals() map[string]interface{} {
	return map[string]interface{}{
		"brutUcret":         i.GrossWage,
		"netUcret":          i.NetWage,
		"gvStopajTutari":    i.GvStoppageAmount,
		"kdvTutari":         i.VatAmount,
		"kdvTevkifatTutari": i.KdvWithholdingAmount,
		"netAlinan":         i.NetReceived,
	}
}

func (i *SelfEmployedReceiptItem) Export() map[string]interface{} {
	m := map[string]interface{}{
		"neIcinAlindigi":   i.Reason,
		"brutUcret":        i.GrossWage,
		"kdvOrani":         i.VatRate,
		"gvStopajOrani":    i.GvStoppageRate,
		"netUcret":         i.NetWage,
		"kdvTevkifatOrani": i.KdvWithholdingRate,
		"netAlinan":        i.NetReceived,
		
		"gvStopajTutari":    i.GvStoppageAmount,
		"kdvTutari":         i.VatAmount,
		"kdvTevkifatTutari": i.KdvWithholdingAmount,
		
		// KeyMap
		"stopaj": i.GvStoppageRate,
		"kdv":    i.VatRate,
	}
	return m
}
