package fatura

import (
	"errors"
	
	"github.com/berkormanli/fatura-go/models"
)

// GetRecipientData retrieves recipient data by tax ID
func (g *Gib) GetRecipientData(taxOrTrID string) (map[string]interface{}, error) {
	gateway := g.GetGateway(PathDispatch)
	payload := map[string]string{"vknTcknn": taxOrTrID}
	
	params := g.setParams("SICIL_VEYA_MERNISTEN_BILGILERI_GETIR", "RG_BASITFATURA", payload)
	
	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return nil, err
	}
	
	if data, ok := resp["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return nil, nil // Or error?
}

// GetUserData retrieves current user data
func (g *Gib) GetUserData() (map[string]interface{}, error) {
	gateway := g.GetGateway(PathDispatch)
	
	params := g.setParams("EARSIV_PORTAL_KULLANICI_BILGILERI_GETIR", "RG_KULLANICI", nil)
	
	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return nil, err
	}
	
	if data, ok := resp["data"].(map[string]interface{}); ok {
		return data, nil
	}
	return nil, nil
}

// UpdateUserData updates user data
func (g *Gib) UpdateUserData(userData models.UserData) error {
	// userData.Export() or similar handling
	// PHP: $userData->export()
	exported := userData.Export()
	
	gateway := g.GetGateway(PathDispatch)
	params := g.setParams("EARSIV_PORTAL_KULLANICI_BILGILERI_KAYDET", "RG_KULLANICI", exported)
	
	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return err
	}
	
	// Check response data? PHP: returns response->get('data') ? true : false
	if resp["data"] != nil {
		return nil
	}
	return errors.New("Update failed")
}

// GetPhoneNumber retrieves registered phone number
func (g *Gib) GetPhoneNumber() (string, error) {
	gateway := g.GetGateway(PathDispatch)
	params := g.setParams("EARSIV_PORTAL_TELEFONNO_SORGULA", "RG_BASITTASLAKLAR", nil)
	
	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return "", err
	}
	
	// PHP: object('data')->telefon ?? null
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if val, ok := data["telefon"].(string); ok {
			return val, nil
		}
	}
	return "", nil
}

// StartSmsVerification starts SMS verification process
func (g *Gib) StartSmsVerification() (string, error) {
	phone, err := g.GetPhoneNumber()
	if err != nil || phone == "" {
		return "", errors.New("Phone number not found or error")
	}
	
	payload := map[string]interface{}{
		"CEPTEL":  phone,
		"KCEPTEL": false,
		"TIP":     "",
	}
	
	gateway := g.GetGateway(PathDispatch)
	params := g.setParams("EARSIV_PORTAL_SMSSIFRE_GONDER", "RG_SMSONAY", payload)
	
	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return "", err
	}
	
	// Return oid
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if oid, ok := data["oid"].(string); ok {
			return oid, nil
		}
	}
	return "", errors.New("OID not found in response")
}

// CompleteSmsVerification completes SMS verification
func (g *Gib) CompleteSmsVerification(code, oid string, documents []string) error {
	// Prepare items to sign
	// PHP: array_map -> ['belgeTuru' => ..., 'ettn' => uuid]
	
	setToSign := make([]map[string]interface{}, len(documents))
	for i, id := range documents {
		setToSign[i] = map[string]interface{}{
			"belgeTuru": string(g.DocumentType),
			"ettn":      id,
		}
	}
	
	payload := map[string]interface{}{
		"DATA":  setToSign,
		"SIFRE": code,
		"OID":   oid,
		"OPR":   1,
	}
	
	gateway := g.GetGateway(PathDispatch)
	params := g.setParams("0lhozfib5410mp", "RG_SMSONAY", payload)
	
	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return err
	}
	
	// Check result
	// PHP: object('data')->sonuc === '1'
	if data, ok := resp["data"].(map[string]interface{}); ok {
		if res, ok := data["sonuc"].(string); ok && res == "1" {
			g.rowCount = len(documents)
			return nil
		}
		// If float64 or int
		if res, ok := data["sonuc"].(float64); ok && res == 1 {
			g.rowCount = len(documents)
			return nil
		}
	}
	
	return errors.New("Verification failed or Result != 1")
}
