package dto

type ErrorDetailInfo struct {
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
	Field   string `json:"field,omitempty"`
}

type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message"`
	Details []ErrorDetailInfo `json:"details,omitempty"`
}

type TaxBanditsAuthError struct {
	Id      string `json:"Id"`
	Message string `json:"Message"`
	Name    string `json:"Name"`
}

type TaxBanditsAuthResponse struct {
	StatusName    string `json:"StatusName"`
	StatusCode    int    `json:"StatusCode"`
	StatusMessage string `json:"StatusMessage"`
	AccessToken   string `json:"AccessToken"`
	// Errors        []TaxBanditsAuthError `json:"Errors"`
	// ExpiresIn     int                   `json:"ExpiresIn"`
	// TokenType     string                `json:"TokenType"`
}

type TaxBanditsAuthAPIResponse struct {
	Status      string `json:"status"`
	Code        int    `json:"code"`
	Message     string `json:"message"`
	AccessToken string `json:"access_token"`
}

type CreateBusinessResponse struct {
	StatusCode    int                   `json:"StatusCode"`
	StatusName    string                `json:"StatusName"`
	StatusMessage string                `json:"StatusMessage"`
	BusinessID    string                `json:"BusinessId"`
	PayerRef      string                `json:"PayerRef"`
	IsEIN         bool                  `json:"IsEIN"`
	EINorSSN      string                `json:"EINorSSN"`
	BusinessNm    string                `json:"BusinessNm"`
	FirstNm       string                `json:"FirstNm"`
	MiddleNm      string                `json:"MiddleNm"`
	LastNm        string                `json:"LastNm"`
	Suffix        string                `json:"Suffix"`
	Errors        []TaxBanditsAuthError `json:"Errors"`
}

type CreateBusinessAPIResponse struct {
	Status     string `json:"status"`
	Code       int    `json:"code"`
	Message    string `json:"message"`
	BusinessID string `json:"business_id,omitempty"`
}
