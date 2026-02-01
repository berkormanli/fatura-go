package models

type UserData struct {
	VknTckn         string `json:"vknTckn"`
	Title           string `json:"unvan"`
	Name            string `json:"ad"`
	Surname         string `json:"soyad"`
	Street          string `json:"cadde"`
	BuildingName    string `json:"apartmanAdi"`
	BuildingNumber  string `json:"apartmanNo"`
	DoorNumber      string `json:"kapiNo"`
	Town            string `json:"kasaba"`
	District        string `json:"ilce"`
	City            string `json:"il"`
	ZipCode         string `json:"postaKodu"`
	Country         string `json:"ulke"`
	Phone           string `json:"telNo"`
	Fax             string `json:"faksNo"`
	Email           string `json:"ePostaAdresi"`
	Website         string `json:"webSitesiAdresi"`
	TaxOffice       string `json:"vergiDairesi"`
	RegistryNumber  string `json:"sicilNo"`
	BusinessCenter  string `json:"isMerkezi"`
	MersisNumber    string `json:"mersisNo"`
}

func (u *UserData) Export() map[string]interface{} {
	return map[string]interface{}{
		"vknTckn":         u.VknTckn,
		"unvan":           u.Title,
		"ad":              u.Name,
		"soyad":           u.Surname,
		"cadde":           u.Street,
		"apartmanAdi":     u.BuildingName,
		"apartmanNo":      u.BuildingNumber,
		"kapiNo":          u.DoorNumber,
		"kasaba":          u.Town,
		"ilce":            u.District,
		"il":              u.City,
		"postaKodu":       u.ZipCode,
		"ulke":            u.Country,
		"telNo":           u.Phone,
		"faksNo":          u.Fax,
		"ePostaAdresi":    u.Email,
		"webSitesiAdresi": u.Website,
		"vergiDairesi":    u.TaxOffice,
		"sicilNo":         u.RegistryNumber,
		"isMerkezi":       u.BusinessCenter,
		"mersisNo":        u.MersisNumber,
	}
}
