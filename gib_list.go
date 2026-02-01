package fatura

import (
	"strings"

	"github.com/berkormanli/fatura-go/enums"
	"github.com/berkormanli/fatura-go/errors"
	"github.com/berkormanli/fatura-go/internal"
)

// GetAll retrieves all documents in range
func (g *Gib) GetAll(startDate, endDate string) ([]map[string]interface{}, error) {
    if !internal.ValidateDate(startDate) || !internal.ValidateDate(endDate) {
        return nil, errors.NewInvalidFormatError("Tarih formatı geçersiz", nil)
    }

    hangiTip := enums.TypeEArsivFatura
    if g.TestMode {
        hangiTip = enums.TypeEArsivDiger
    }
    
    payload := map[string]interface{}{
        "baslangic": startDate,
        "bitis":     endDate,
        "hangiTip":  hangiTip,
    }
    
    gateway := g.GetGateway(PathDispatch)
    params := g.setParams("EARSIV_PORTAL_TASLAKLARI_GETIR", "RG_TASLAKLAR", payload)
    
    resp, err := g.client.Request(gateway, params, true)
    if err != nil {
        return nil, err
    }
    
    return g.filterDocuments(resp["data"])
}

// GetAllIssuedToMe retrieves documents issued to the user
func (g *Gib) GetAllIssuedToMe(startDate, endDate string) ([]map[string]interface{}, error) {
    if !internal.ValidateDate(startDate) || !internal.ValidateDate(endDate) {
        return nil, errors.NewInvalidFormatError("Tarih formatı geçersiz", nil)
    }
    
    payload := map[string]interface{}{
        "baslangic":            startDate,
        "bitis":                endDate,
        "hourlySearchInterval": "NONE",
    }
    
    gateway := g.GetGateway(PathDispatch)
    params := g.setParams("EARSIV_PORTAL_ADIMA_KESILEN_BELGELERI_GETIR", "RG_ALICI_TASLAKLAR", payload)
    
    resp, err := g.client.Request(gateway, params, true)
    if err != nil {
        return nil, err
    }
    
    return g.filterDocuments(resp["data"])
}

func (g *Gib) filterDocuments(data interface{}) ([]map[string]interface{}, error) {
    // data might be []interface{}, need to cast to []map[string]interface{}
    
    var docs []map[string]interface{}
    if list, ok := data.([]interface{}); ok {
        for _, item := range list {
            if m, ok := item.(map[string]interface{}); ok {
                docs = append(docs, m)
            }
        }
    } else if list, ok := data.([]map[string]interface{}); ok {
        docs = list
    } else {
        return []map[string]interface{}{}, nil
    }
    
    // Apply filters
    if len(g.filters) > 0 {
        var filtered []map[string]interface{}
        for _, doc := range docs {
            match := true
            for k, v := range g.filters {
                if val, ok := doc[k].(string); ok {
                    // Exact or contains (case insensitive)
                    if val != v && !strings.Contains(strings.ToLower(val), strings.ToLower(v)) {
                         match = false
                         break
                    }
                } else {
                    match = false
                    break
                }
            }
            if match {
                filtered = append(filtered, doc)
            }
        }
        docs = filtered
    }
    
    g.rowCount = len(docs)
    g.filters = make(map[string]string) // Reset filters
    
    // Helper to reverse if needed (Skip for now or simple implementation)
    // Helper to slice if limit set
    if len(g.limit) == 2 {
        offset, count := g.limit[0], g.limit[1]
        if offset < len(docs) {
             end := offset + count
             if end > len(docs) {
                 end = len(docs)
             }
             docs = docs[offset:end]
        } else {
            docs = []map[string]interface{}{}
        }
        g.limit = []int{}
    }

    return docs, nil
}

// Filter setters ...
func (g *Gib) FindEttn(ettn string) *Gib {
    g.filters["ettn"] = ettn
    return g
}

func (g *Gib) SetLimit(limit, offset int) *Gib {
    g.limit = []int{offset, limit}
    return g
}

// Internal error helper?
type invalidFormatError struct{ msg string }
func (e invalidFormatError) Error() string { return e.msg }
func (e invalidFormatError) GetRequest() interface{} { return nil }
