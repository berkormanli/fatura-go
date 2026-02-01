package enums

type Tax string

const (
	TaxBankaMuameleleri  Tax = "0021"
	TaxKKDFKesintisi     Tax = "0061"
	TaxOTV1Liste         Tax = "0071"
	TaxOTV2Liste         Tax = "9077"
	TaxOTV3Liste         Tax = "0073"
	TaxOTV4Liste         Tax = "0074"
	TaxOTV3AListe        Tax = "0075"
	TaxOTV3BListe        Tax = "0076"
	TaxOTV3CListe        Tax = "0077"
	TaxDamga             Tax = "1047"
	TaxDamga5035         Tax = "1048"
	TaxOzelIletisim      Tax = "4080"
	TaxOzelIletisim5035  Tax = "4081"
	TaxKDVTevkifat       Tax = "9015" // Tevkifat
	TaxBSMV4961          Tax = "9021"
	TaxBorsaTescil       Tax = "8001"
	TaxEnerjiFonu        Tax = "8002"
	TaxElkHavagazTuketim Tax = "4071"
	TaxTRTPayi           Tax = "8004"
	TaxElkTuketim        Tax = "8005"
	TaxTKKullanim        Tax = "8006"
	TaxTKRuhsat          Tax = "8007"
	TaxCevreTemizlik     Tax = "8008"
	TaxGVStopaj          Tax = "0003" // Stopaj
	TaxKVStopaj          Tax = "0011" // Stopaj
	TaxMeraFonu          Tax = "9040" // Stopaj
	TaxOTV1ListeTevkifat Tax = "4171" // Tevkifat
	TaxBelOdHalRusum     Tax = "9944"
	TaxKonaklama         Tax = "0059"
	TaxSGKPrim           Tax = "SGK_PRIM" // Müstahsil
)

func (t Tax) Alias() string {
	switch t {
	case TaxBankaMuameleleri:
		return "Banka Muameleleri Vergisi"
	case TaxKKDFKesintisi:
		return "KKDF Kesintisi"
	case TaxOTV1Liste:
		return "ÖTV 1. Liste"
	case TaxOTV2Liste:
		return "ÖTV 2. Liste"
	case TaxOTV3Liste:
		return "ÖTV 3. Liste"
	case TaxOTV4Liste:
		return "ÖTV 4. Liste"
	case TaxOTV3AListe:
		return "ÖTV 3A Liste"
	case TaxOTV3BListe:
		return "ÖTV 3B Liste"
	case TaxOTV3CListe:
		return "ÖTV 3C Liste"
	case TaxDamga:
		return "Damga Vergisi"
	case TaxDamga5035:
		return "5035 Sayılı Kanuna Göre Damga Vergisi"
	case TaxOzelIletisim:
		return "Özel İletişim Vergisi"
	case TaxOzelIletisim5035:
		return "5035 Sayılı Kanuna Göre Özel İletişim Vergisi"
	case TaxKDVTevkifat:
		return "KDV Tevkifat"
	case TaxBSMV4961:
		return "Banka ve Sigorta Muameleleri Vergisi"
	case TaxBorsaTescil:
		return "Borsa Tescil Ücreti"
	case TaxEnerjiFonu:
		return "Enerji Fonu"
	case TaxElkHavagazTuketim:
		return "Elektrik Havagaz Tüketim Vergisi"
	case TaxTRTPayi:
		return "TRT Payı"
	case TaxElkTuketim:
		return "Elektrik Tüketim Vergisi"
	case TaxTKKullanim:
		return "TK Kullanım"
	case TaxTKRuhsat:
		return "TK Ruhsat"
	case TaxCevreTemizlik:
		return "Çevre Temizlik Vergisi"
	case TaxGVStopaj:
		return "Gelir Vergisi Stopajı"
	case TaxKVStopaj:
		return "Kurumlar Vergisi Stopajı"
	case TaxMeraFonu:
		return "Mera Fonu"
	case TaxOTV1ListeTevkifat:
		return "ÖTV 1. Liste Tevkifat"
	case TaxBelOdHalRusum:
		return "Belediyelere Ödenen Hal Rüsumu"
	case TaxKonaklama:
		return "Konaklama Vergisi"
	case TaxSGKPrim:
		return "SGK Prim Kesintisi"
	default:
		return ""
	}
}

func (t Tax) HasVat() bool {
	switch t {
	case TaxKKDFKesintisi, TaxOTV1Liste, TaxOTV2Liste, TaxOTV3Liste, TaxOTV4Liste,
		TaxOTV3AListe, TaxOTV3BListe, TaxOTV3CListe, TaxEnerjiFonu, TaxElkHavagazTuketim,
		TaxTRTPayi, TaxElkTuketim, TaxOTV1ListeTevkifat, TaxBelOdHalRusum:
		return true
	default:
		return false
	}
}

func (t Tax) IsStoppage() bool {
	switch t {
	case TaxKDVTevkifat, TaxGVStopaj, TaxKVStopaj, TaxMeraFonu, TaxSGKPrim:
		return true
	default:
		return false
	}
}

func (t Tax) IsWithholding() bool {
	switch t {
	case TaxKDVTevkifat, TaxOTV1ListeTevkifat:
		return true
	default:
		return false
	}
}

func (t Tax) HasDefaultRate() bool {
	switch t {
	case TaxOTV1Liste, TaxOTV1ListeTevkifat, TaxKonaklama:
		return true
	default:
		return false
	}
}

func (t Tax) DefaultRate() int {
	switch t {
	case TaxOTV1Liste:
		return 0
	case TaxOTV1ListeTevkifat:
		return 100
	case TaxKonaklama:
		return 2
	default:
		return 0
	}
}

type TaxCodeInfo struct {
	Rate int
	Name string
}

func (t Tax) Codes() map[int]TaxCodeInfo {
	if t == TaxKDVTevkifat {
		return map[int]TaxCodeInfo{
			601:      {40, "Yapım İşleri ile Bu İşlerle Birlikte İfa Edilen Mühendislik-Mimarlık ve Etüt-Proje Hizmetleri [KDVGUT-(I/C-2.1.3.2.1)]"},
			602:      {90, "Etüt, plan-proje, danışmanlık, denetim vb"},
			603:      {70, "Makine, Teçhizat, Demirbaş ve Taşıtlara Ait Tadil, Bakım ve Onarım Hizmetleri [KDVGUT- (I/C-2.1.3.2.3)]"},
			604:      {50, "Yemek servis hizmeti"},
			605:      {50, "Organizasyon hizmeti"},
			606:      {90, "İşgücü temin hizmetleri"},
			607:      {90, "Özel güvenlik hizmeti"},
			608:      {90, "Yapı denetim hizmetleri"},
			609:      {70, "Fason Olarak Yaptırılan Tekstil ve Konfeksiyon İşleri, Çanta ve Ayakkabı Dikim İşleri ve Bu İşlere Aracılık Hizmetleri [KDVGUT-(I/C-2.1.3.2.7)]"},
			610:      {90, "Turistik mağazalara verilen müşteri bulma/ götürme hizmetleri"},
			611:      {90, "Spor kulüplerinin yayın, reklam ve isim hakkı gelirlerine konu işlemleri"},
			612:      {90, "Temizlik Hizmeti [KDVGUT-(I/C-2.1.3.2.10)]"},
			613:      {90, "Çevre, Bahçe ve Bakım Hizmetleri [KDVGUT-(I/C-2.1.3.2.11)]"},
			614:      {50, "Servis taşımacıliğı"},
			615:      {70, "Her Türlü Baskı ve Basım Hizmetleri [KDVGUT-(I/C-2.1.3.2.12)]"},
			616:      {50, "Diğer Hizmetler [KDVGUT-(I/C-2.1.3.2.13)]"},
			617:      {70, "Hurda metalden elde edilen külçe teslimleri"},
			618:      {70, "Hurda Metalden Elde Edilenler Dışındaki Bakır, Çinko, Demir Çelik, Alüminyum ve Kurşun Külçe Teslimi [KDVGUT-(I/C-2.1.3.3.1)]"},
			619:      {70, "Bakir, çinko ve alüminyum ürünlerinin teslimi"},
			620:      {70, "istisnadan vazgeçenlerin hurda ve atık teslimi"},
			621:      {90, "Metal, plastik, lastik, kauçuk, kâğit ve cam hurda ve atıklardan elde edilen hammadde teslimi"},
			622:      {90, "Pamuk, tiftik, yün ve yapaği ile ham post ve deri teslimleri"},
			623:      {50, "Ağaç ve orman ürünleri teslimi"},
			624:      {20, "Yük Taşımacılığı Hizmeti [KDVGUT-(I/C-2.1.3.2.11)]"},
			625:      {30, "Ticari Reklam Hizmetleri [KDVGUT-(I/C-2.1.3.2.15)]"},
			626:      {20, "Diğer Teslimler [KDVGUT-(I/C-2.1.3.3.7.)]"},
			627:      {50, "Demir-Çelik Ürünlerinin Teslimi [KDVGUT-(I/C-2.1.3.3.8)]"},
			// '627-Ex': {40, "Demir-Çelik Ürünlerinin Teslimi [KDVGUT-(I/C-2.1.3.3.8)] (01/11/2022 tarihi öncesi)"}, // Map keys must be same type
			801: {100, "[Tam Tevkifat] Yapım İşleri ile Bu İşlerle Birlikte İfa Edilen Mühendislik-Mimarlık ve Etüt-Proje Hizmetleri[KDVGUT-(I/C-2.1.3.2.1)]"},
			802: {100, "[Tam Tevkifat] Etüt, Plan-Proje, Danışmanlık, Denetim ve Benzeri Hizmetler[KDVGUT-(I/C-2.1.3.2.2)]"},
			803: {100, "[Tam Tevkifat] Makine, Teçhizat, Demirbaş ve Taşıtlara Ait Tadil, Bakım ve Onarım Hizmetleri[KDVGUT- (I/C-2.1.3.2.3)]"},
			804: {100, "[Tam Tevkifat] Yemek Servis Hizmeti[KDVGUT-(I/C-2.1.3.2.4)]"},
			805: {100, "[Tam Tevkifat] Organizasyon Hizmeti[KDVGUT-(I/C-2.1.3.2.4)]"},
			806: {100, "[Tam Tevkifat] İşgücü Temin Hizmetleri[KDVGUT-(I/C-2.1.3.2.5)]"},
			807: {100, "[Tam Tevkifat] Özel Güvenlik Hizmeti[KDVGUT-(I/C-2.1.3.2.5)]"},
			808: {100, "[Tam Tevkifat] Yapı Denetim Hizmetleri[KDVGUT-(I/C-2.1.3.2.6)]"},
			809: {100, "[Tam Tevkifat] Fason Olarak Yaptırılan Tekstil ve Konfeksiyon İşleri, Çanta ve Ayakkabı Dikim İşleri ve Bu İşlere Aracılık Hizmetleri[KDVGUT-(I/C-2.1.3.2.7)]"},
			810: {100, "[Tam Tevkifat] Turistik Mağazalara Verilen Müşteri Bulma/ Götürme Hizmetleri[KDVGUT-(I/C-2.1.3.2.8)]"},
			811: {100, "[Tam Tevkifat] Spor Kulüplerinin Yayın, Reklâm ve İsim Hakkı Gelirlerine Konu İşlemleri[KDVGUT-(I/C-2.1.3.2.9)]"},
			812: {100, "[Tam Tevkifat] Temizlik Hizmeti[KDVGUT-(I/C-2.1.3.2.10)]"},
			813: {100, "[Tam Tevkifat] Çevreve Bahçe Bakım Hizmetleri[KDVGUT-(I/C-2.1.3.2.10)]"},
			814: {100, "[Tam Tevkifat] Servis Taşımacılığı Hizmeti[KDVGUT-(I/C-2.1.3.2.11)]"},
			815: {100, "[Tam Tevkifat] Her Türlü Baskı ve Basım Hizmetleri[KDVGUT-(I/C-2.1.3.2.12)]"},
			816: {100, "[Tam Tevkifat] Hurda Metalden Elde Edilen Külçe Teslimleri[KDVGUT-(I/C-2.1.3.3.1)]"},
			817: {100, "[Tam Tevkifat] Hurda Metalden Elde Edilenler Dışındaki Bakır, Çinko, Demir Çelik, Alüminyum ve Kurşun Külçe Teslimi [KDVGUT-(I/C-2.1.3.3.1)]"},
			818: {100, "[Tam Tevkifat] Bakır, Çinko, Alüminyum ve Kurşun Ürünlerinin Teslimi[KDVGUT-(I/C-2.1.3.3.2)]"},
			819: {100, "[Tam Tevkifat] İstisnadan Vazgeçenlerin Hurda ve Atık Teslimi[KDVGUT-(I/C-2.1.3.3.3)]"},
			820: {100, "[Tam Tevkifat] Metal, Plastik, Lastik, kauçuk, Kâğıt ve Cam Hurda ve Atıklardan Elde Edilen Hammadde Teslimi[KDVGUT-(I/C-2.1.3.3.4)]"},
			821: {100, "[Tam Tevkifat] Pamuk, Tiftik, Yün ve Yapağı İle Ham Post ve Deri Teslimleri[KDVGUT-(I/C-2.1.3.3.5)]"},
			822: {100, "[Tam Tevkifat] Ağaç ve Orman Ürünleri Teslimi[KDVGUT-(I/C-2.1.3.3.6)]"},
			823: {100, "[Tam Tevkifat] Yük Taşımacılığı Hizmeti [KDVGUT-(I/C-2.1.3.2.11)]"},
			824: {100, "[Tam Tevkifat] Ticari Reklam Hizmetleri [KDVGUT-(I/C-2.1.3.2.15)]"},
			825: {100, "[Tam Tevkifat] Demir-Çelik Ürünlerinin Teslimi [KDVGUT-(I/C-2.1.3.3.8)]"},
		}
	}
	return map[int]TaxCodeInfo{}
}

func (t Tax) GetRate(code int) (int, bool) {
	if codes := t.Codes(); len(codes) > 0 {
		if info, ok := codes[code]; ok {
			return info.Rate, true
		}
	}
	return 0, false
}
