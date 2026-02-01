package enums

type Type string

const (
	TypeEArsivFatura Type = "5000/30000"
	TypeEArsivDiger  Type = "Buyuk"
)

func (t Type) Alias() string {
	switch t {
	case TypeEArsivFatura:
		return "E-Arşiv Fatura"
	case TypeEArsivDiger:
		return "E-Arşiv Diğer"
	default:
		return ""
	}
}
