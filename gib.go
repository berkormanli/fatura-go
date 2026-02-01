package fatura

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/berkormanli/fatura-go/enums"
	"github.com/berkormanli/fatura-go/internal/client"
)

const (
	GatewayProd    = "https://earsivportal.efatura.gov.tr"
	GatewayTest    = "https://earsivportaltest.efatura.gov.tr"
	GatewayDigital = "https://dijital.gib.gov.tr"

	PathEsign    = "/earsiv-services/esign"
	PathLogin    = "/earsiv-services/assos-login"
	PathDispatch = "/earsiv-services/dispatch"
	PathDownload = "/earsiv-services/download"
	
	// New Digital Gateway Paths
	PathCaptcha      = "/apigateway/captcha/getnewcaptcha"
	PathDigitalLogin = "/apigateway/auth/tdvd/login"
)

type Gib struct {
	DocumentType enums.DocumentType
	TestMode     bool
	Username     string
	Password     string
	Token        string
	
	client       *client.Client
	
	// Filtering state
	filters    map[string]string
	column     []string
	limit      []int
	sortByDesc bool
	rowCount   int
	lastID     string
}

func NewGib(docType enums.DocumentType) *Gib {
	return &Gib{
		DocumentType: docType,
		client:       client.NewClient(30 * time.Second), // Default timeout
		filters:      make(map[string]string),
	}
}

func (g *Gib) SetTestMode(test bool) *Gib {
	g.TestMode = test
	return g
}

func (g *Gib) SetCredentials(username, password string) *Gib {
	g.Username = username
	g.Password = password
	return g
}

func (g *Gib) GetCredentials() (string, string) {
	return g.Username, g.Password
}

func (g *Gib) SetTestCredentials(username, password string) *Gib {
	if username != "" && password != "" {
		g.TestMode = true
		return g.SetCredentials(username, password)
	}
	
	// Fetch test credentials
	// PHP: new Client(esign, [assoscmd: kullaniciOner])
	g.TestMode = true
	gateway := g.GetGateway(PathEsign)
	
	resp, err := g.client.Request(gateway, map[string]string{
		"assoscmd": "kullaniciOner",
		"rtype":    "json",
	}, true)
	
	if err == nil && resp != nil {
		if userid, ok := resp["userid"].(string); ok && userid != "" {
			return g.SetCredentials(userid, "1")
		}
	}
	// Fallback or error? PHP throws exception if userid is empty.
	// We'll just return g, user should check credentials or we return error from this method?
	// Chained methods usually don't return error. 
	// But network op is involved here.
	// We'll stick to PHP behavior: set if successful, else it remains empty?
	return g
}

func (g *Gib) SetToken(token string) *Gib {
	g.Token = token
	return g
}

func (g *Gib) GetGateway(path string) string {
	base := GatewayProd
	if g.TestMode {
		base = GatewayTest
	}
	return base + path
}

func (g *Gib) Login(username, password string) error {
	if username != "" && password != "" {
		g.SetCredentials(username, password)
	}

	gateway := g.GetGateway(PathLogin)
	cmd := "anologin"
	if g.TestMode {
		cmd = "login"
	}

	params := map[string]string{
		"assoscmd": cmd,
		"userid":   g.Username,
		"sifre":    g.Password,
		"sifre2":   g.Password,
		"parola":   g.Password,
	}

	resp, err := g.client.Request(gateway, params, true)
	if err != nil {
		return err
	}

	if token, ok := resp["token"].(string); ok {
		g.SetToken(token)
		return nil
	}
	return errors.New("Login succeeded but token not found in response")
}

// CaptchaResponse holds the captcha data from GIB
type CaptchaResponse struct {
	ImageBase64 string `json:"captchaImgBase64"`
	CaptchaID   string `json:"cid"`
}

// GetCaptcha fetches a new captcha from GIB digital gateway
func (g *Gib) GetCaptcha() (*CaptchaResponse, error) {
	gateway := GatewayDigital + PathCaptcha
	
	resp, err := g.client.Request(gateway, nil, false) // GET request
	if err != nil {
		return nil, err
	}
	
	captcha := &CaptchaResponse{}
	if img, ok := resp["captchaImgBase64"].(string); ok {
		captcha.ImageBase64 = img
	}
	if cid, ok := resp["cid"].(string); ok {
		captcha.CaptchaID = cid
	}
	
	if captcha.ImageBase64 == "" || captcha.CaptchaID == "" {
		return nil, errors.New("invalid captcha response")
	}
	
	return captcha, nil
}

// LoginWithCaptcha performs login using the new digital gateway with captcha
func (g *Gib) LoginWithCaptcha(username, password, captchaSolution, captchaID string) error {
	g.SetCredentials(username, password)
	
	gateway := GatewayDigital + PathDigitalLogin
	
	payload := map[string]string{
		"dk":      captchaSolution,
		"userid":  username,
		"sifre":   password,
		"imageId": captchaID,
	}
	
	resp, err := g.client.RequestJSON(gateway, payload)
	if err != nil {
		return err
	}
	
	// Check result
	if result, ok := resp["result"].(bool); ok && result {
		if token, ok := resp["token"].(string); ok {
			g.SetToken(token)
			return nil
		}
	}
	
	return errors.New("login failed or invalid response")
}

func (g *Gib) Logout() error {
	gateway := g.GetGateway(PathLogin)
	params := map[string]string{
		"assoscmd": "logout",
		"token":    g.Token,
	}
	// Ignore response, just fire and forget?
	_, _ = g.client.Request(gateway, params, true)
	
	g.Token = ""
	g.Username = ""
	g.Password = ""
	return nil
}

// setParams helper
func (g *Gib) setParams(cmd, pageName string, payload interface{}) map[string]string {
	jpBytes, _ := json.Marshal(payload)
	if payload == nil {
		jpBytes = []byte("{}")
	} else if m, ok := payload.(map[string]interface{}); ok && len(m) == 0 {
		 jpBytes = []byte("{}")
	}
	
	return map[string]string{
		"callid":   uuid.New().String(),
		"token":    g.Token,
		"cmd":      cmd,
		"pageName": pageName,
		"jp":       string(jpBytes),
	}
}
