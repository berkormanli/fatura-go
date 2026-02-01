# fatura-go

GİB e-Arşiv Fatura Portalına (2.000 / 5.000 TL) entegrasyon sağlayan, `mlevent/fatura` kütüphanesinin (PHP) resmi olmayan Golang portudur. Bu kütüphane ile fatura ve makbuz oluşturabilir, imzalayabilir ve portal işlemlerini yönetebilirsiniz.

## Özellikler

- Fatura oluşturma, düzenleme ve silme
- Müstahsil Makbuzu ve Serbest Meslek Makbuzu desteği
- SMS ile doğrulama ve imzalama
- Belge sorgulama ve indirme (HTML/ZIP)
- İptal ve itiraz talepleri oluşturma
- Kullanıcı bilgileri güncelleme
- Test modu desteği

## Kurulum

```bash
go get github.com/berkormanli/fatura-go
```

## Kullanım

### Fatura Oluşturma

```go
package main

import (
    "fmt"
    "github.com/berkormanli/fatura-go"
    "github.com/berkormanli/fatura-go/enums"
    "github.com/berkormanli/fatura-go/models"
)

func main() {
    // Servis oluşturma
    service := fatura.NewGib(enums.DocumentTypeInvoice).SetTestCredentials("", "")
    
    // Giriş yapma
    if err := service.Login("", ""); err != nil {
        panic(err)
    }
    defer service.Logout()

    // Fatura detayları
    invoice, err := models.NewInvoice("11111111111", 
        models.WithRecipientName("Ad", "Soyad"),
        models.WithRecipientTitle("Unvan Ltd. Şti."),
        models.WithAddress("Papatya Sk.", "Bursa", "Türkiye", "Nilüfer"),
    )
    if err != nil {
         panic(err)
    }

    // Kalem ekleme
    item, _ := models.NewInvoiceItem("Danışmanlık Hizmeti", 1, 1000, 18)
    invoice.AddItem(item)

    // Not ekleme
    invoice.SetNote("İşbu fatura ...")

    // Taslak oluşturma
    if err := service.CreateDraft(invoice); err != nil {
        fmt.Println("Hata:", err)
    } else {
        fmt.Println("Fatura başarıyla oluşturuldu:", invoice.GetUUID())
    }
}
```

### Belge Sorgulama

```go
docs, err := service.GetAll("01/01/2023", "31/01/2023")
if err != nil {
    panic(err)
}

for _, doc := range docs {
    fmt.Printf("Belge No: %s, Tutar: %v\n", doc["belgeNumarasi"], doc["odenecekTutar"])
}
```

### Belge İndirme

```go
path, err := service.SaveToDisk("UUID-STRING", "./downloads", "fatura-1")
if err != nil {
    panic(err)
}
fmt.Println("Dosya kaydedildi:", path)
```

## Lisans

MIT
