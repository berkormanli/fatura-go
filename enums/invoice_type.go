package enums

type InvoiceType string

const (
	InvoiceTypeSatis            InvoiceType = "SATIS"
	InvoiceTypeIade             InvoiceType = "IADE"
	InvoiceTypeTevkifat         InvoiceType = "TEVKIFAT"
	InvoiceTypeIstisna          InvoiceType = "ISTISNA"
	InvoiceTypeOzelMatrah       InvoiceType = "OZELMATRAH"
	InvoiceTypeIhracKayitli     InvoiceType = "IHRACKAYITLI"
	InvoiceTypeKonaklamaVergisi InvoiceType = "KONAKLAMAVERGISI"
)

func (i InvoiceType) Alias() string {
	switch i {
	case InvoiceTypeSatis:
		return "Satış"
	case InvoiceTypeIade:
		return "İade"
	case InvoiceTypeTevkifat:
		return "Tevkifat"
	case InvoiceTypeIstisna:
		return "İstisna"
	case InvoiceTypeOzelMatrah:
		return "Özel Matrah"
	case InvoiceTypeIhracKayitli:
		return "İhraç Kayıtlı"
	case InvoiceTypeKonaklamaVergisi:
		return "Konaklama Vergisi"
	default:
		return ""
	}
}

func (i InvoiceType) Reasons() map[int]string {
	if i == InvoiceTypeOzelMatrah {
		return map[int]string{
			801: "Milli Piyango, Spor Toto vb. Oyunlar",
			802: "At yarışları ve diğer müşterek bahis ve talih oyunları",
			803: "Profesyonel Sanatçıların Yer Aldığı Gösteriler, Konserler, Profesyonel Sporcuların Katıldığı Sportif Faaliyetler, Maçlar, Yarışlar ve Yarışmalar",
			804: "Gümrük Depolarında ve Müzayede Mahallerinde Yapılan Satışla",
			805: "Altından Mamül veya Altın İçeren Ziynet Eşyaları İle Sikke Altınların Teslimi",
			806: "Tütün Mamülleri",
			807: "Muzır Neşriyat Kapsamındaki  Gazete, Dergi vb. Periyodik Yayınlar",
			808: "Gümüşten Mamul veya Gümüş İçeren Ziynet Eşyaları ile Sikke Gümüşlerin Teslimi",
			809: "Belediyeler taraf. yap. şehiriçi yolcu taşımacılığında kullanılan biletlerin ve kartların bayiler tarafından satışı",
			810: "Ön Ödemeli Elektronik Haberleşme Hizmetleri",
			811: "TŞOF Tarafından Araç Plakaları ile Sürücü Kurslarında Kullanılan Bir Kısım Evrakın Teslimi",
			812: "KDV Uygulanmadan Alınan İkinci El Motorlu Kara Taşıtı veya Taşınmaz Teslimi",
		}
	}
	return map[int]string{}
}
