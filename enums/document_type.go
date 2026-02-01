package enums

type DocumentType string

const (
	DocumentTypeInvoice             DocumentType = "FATURA"
	DocumentTypeProducerReceipt     DocumentType = "MÜSTAHSİL MAKBUZU"
	DocumentTypeSelfEmployedReceipt DocumentType = "SERBEST MESLEK MAKBUZU"
)
