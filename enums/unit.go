package enums

type Unit string

const (
	UnitGun   Unit = "DAY"
	UnitAy    Unit = "MON"
	UnitYil   Unit = "ANN"
	UnitSaat  Unit = "HUR"
	UnitDk    Unit = "D61"
	UnitSn    Unit = "D62"
	UnitAdet  Unit = "C62"
	UnitPk    Unit = "PA"
	UnitKutu  Unit = "BX"
	UnitMgm   Unit = "MGM"
	UnitGrm   Unit = "GRM"
	UnitKgm   Unit = "KGM"
	UnitLtr   Unit = "LTR"
	UnitTon   Unit = "TNE"
	UnitNt    Unit = "NT"
	UnitGt    Unit = "GT"
	UnitMmt   Unit = "MMT"
	UnitCmt   Unit = "CMT"
	UnitMtr   Unit = "MTR"
	UnitKtm   Unit = "KTM"
	UnitMlt   Unit = "MLT"
	UnitMm3   Unit = "MMQ"
	UnitCm2   Unit = "CMK"
	UnitCmq   Unit = "CMQ"
	UnitM2    Unit = "MTK"
	UnitM3    Unit = "MTQ"
	UnitKjo   Unit = "KJO"
	UnitClt   Unit = "CLT"
	UnitCt    Unit = "CT"
	UnitKwh   Unit = "KWH"
	UnitMwh   Unit = "MWH"
	UnitCct   Unit = "CCT"
	UnitGkj   Unit = "D30"
	UnitKlt   Unit = "D40"
	UnitLpa   Unit = "LPA"
	UnitKgm2  Unit = "B32"
	UnitNcl   Unit = "NCL"
	UnitPr    Unit = "PR"
	UnitKmt   Unit = "R9"
	UnitSet   Unit = "SET"
	UnitT3    Unit = "T3"
	UnitScm   Unit = "Q37"
	UnitNcm   Unit = "Q39"
	UnitMmbtu Unit = "J39"
	UnitCm3   Unit = "G52"
	UnitDzn   Unit = "DZN"
	UnitDm2   Unit = "DMK"
	UnitDmt   Unit = "DMT"
	UnitHar   Unit = "HAR"
	UnitLm    Unit = "LM"
)

func (u Unit) Alias() string {
	switch u {
	case UnitGun:
		return "Gün"
	case UnitAy:
		return "Ay"
	case UnitYil:
		return "Yıl"
	case UnitSaat:
		return "Saat"
	case UnitDk:
		return "Dakika"
	case UnitSn:
		return "Saniye"
	case UnitAdet:
		return "Adet"
	case UnitPk:
		return "Paket"
	case UnitKutu:
		return "Kutu"
	case UnitMgm:
		return "Mg"
	case UnitGrm:
		return "Gram"
	case UnitKgm:
		return "Kg"
	case UnitLtr:
		return "Lt"
	case UnitTon:
		return "Ton"
	case UnitNt:
		return "Net Ton"
	case UnitGt:
		return "Gross ton"
	case UnitMmt:
		return "Mm"
	case UnitCmt:
		return "Cm"
	case UnitMtr:
		return "M"
	case UnitKtm:
		return "Km"
	case UnitMlt:
		return "Ml"
	case UnitMm3:
		return "Mm3"
	case UnitCm2:
		return "Cm2"
	case UnitCmq:
		return "Cm3"
	case UnitM2:
		return "M2"
	case UnitM3:
		return "M3"
	case UnitKjo:
		return "Kj"
	case UnitClt:
		return "Cl"
	case UnitCt:
		return "Karat"
	case UnitKwh:
		return "Kwh"
	case UnitMwh:
		return "Mwh"
	case UnitCct:
		return "Ton Başına Taşıma Kapasitesi"
	case UnitGkj:
		return "Brüt Kalori"
	case UnitKlt:
		return "1000 Lt"
	case UnitLpa:
		return "Saf Alkol Lt"
	case UnitKgm2:
		return "Kg M2"
	case UnitNcl:
		return "Hücre Adet"
	case UnitPr:
		return "Çift"
	case UnitKmt:
		return "1000 M3"
	case UnitSet:
		return "Set"
	case UnitT3:
		return "1000 Adet"
	case UnitScm:
		return "Scm"
	case UnitNcm:
		return "Ncm"
	case UnitMmbtu:
		return "Mmbtu"
	case UnitCm3:
		return "Cm³"
	case UnitDzn:
		return "Düzine"
	case UnitDm2:
		return "Dm2"
	case UnitDmt:
		return "Dm"
	case UnitHar:
		return "Ha"
	case UnitLm:
		return "Metretül (LM)"
	default:
		return ""
	}
}
