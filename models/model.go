package models

type Model interface {
	GetUUID() string
	Export() map[string]interface{}
}

type ItemModel interface {
	GetTotals() map[string]interface{} // Changed to return map to match PHP usage
	Export() map[string]interface{}
	Prepare(parent Model) ItemModel
}
