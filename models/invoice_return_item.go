package models

import (
	"errors"
	"github.com/berkormanli/fatura-go/internal"
)

type InvoiceReturnItem struct {
	InvoiceNumber string `json:"faturaNo"`
	IssueDate     string `json:"duzenlenmeTarihi"`
}

func NewInvoiceReturnItem(invoiceNo string, date string) (*InvoiceReturnItem, error) {
	if !internal.ValidateInvoiceNumber(invoiceNo) {
		return nil, errors.New("Fatura numarası geçerli formatta değil.")
	}
	if !internal.ValidateDate(date) {
		return nil, errors.New("Tarih geçerli formatta değil.")
	}

	return &InvoiceReturnItem{
		InvoiceNumber: invoiceNo,
		IssueDate:     date,
	}, nil
}

func (i *InvoiceReturnItem) Export() map[string]interface{} {
	return map[string]interface{}{
		"faturaNo":         i.InvoiceNumber,
		"duzenlenmeTarihi": i.IssueDate,
	}
}
