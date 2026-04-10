package dto

type TaxBanditsJWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type TaxBanditsJWTPayload struct {
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Iat int64  `json:"iat"`
}

type CreateBusinessSigningAuthority struct {
	Name               string `json:"Name,omitempty"`
	Phone              string `json:"Phone,omitempty"`
	BusinessMemberType string `json:"BusinessMemberType,omitempty"`
}

type CreateBusinessForm1042SDetails struct {
	WHAgtCh3Cd string `json:"WHAgtCh3Cd,omitempty"`
	WHAgtCh4Cd string `json:"WHAgtCh4Cd,omitempty"`
	WHAgtGIIN  string `json:"WHAgtGIIN,omitempty"`
	FTIN       string `json:"FTIN,omitempty"`
	Country    string `json:"Country,omitempty"`
}

type CreateBusinessUSAddress struct {
	Address1 string `json:"Address1,omitempty"`
	Address2 string `json:"Address2,omitempty"`
	City     string `json:"City,omitempty"`
	State    string `json:"State,omitempty"`
	ZipCd    string `json:"ZipCd,omitempty"`
}

type CreateBusinessForeignAddress struct {
	Address1          string `json:"Address1,omitempty"`
	Address2          string `json:"Address2,omitempty"`
	City              string `json:"City,omitempty"`
	ProvinceOrStateNm string `json:"ProvinceOrStateNm,omitempty"`
	Country           string `json:"Country,omitempty"`
	PostalCd          string `json:"PostalCd,omitempty"`
}

type CreateBusinessACADetails struct {
	FirstName          string `json:"FirstName,omitempty"`
	MiddleName         string `json:"MiddleName,omitempty"`
	LastName           string `json:"LastName,omitempty"`
	Suffix             string `json:"Suffix,omitempty"`
	Phone              string `json:"Phone,omitempty"`
	IsGovernmentalUnit bool   `json:"IsGovernmentalUnit,omitempty"`
}

type CreateBusinessRequest struct {
	BusinessNm           string                          `json:"BusinessNm,omitempty"`
	FirstNm              string                          `json:"FirstNm,omitempty"`
	MiddleNm             string                          `json:"MiddleNm,omitempty"`
	LastNm               string                          `json:"LastNm,omitempty"`
	Suffix               string                          `json:"Suffix,omitempty"`
	PayerRef             string                          `json:"PayerRef,omitempty"`
	TradeNm              string                          `json:"TradeNm,omitempty"`
	IsEIN                bool                            `json:"IsEIN"`
	EINorSSN             string                          `json:"EINorSSN" binding:"required"`
	IsDefaultBusiness    bool                            `json:"IsDefaultBusiness"`
	Email                string                          `json:"Email,omitempty"`
	ContactNm            string                          `json:"ContactNm,omitempty"`
	Phone                string                          `json:"Phone,omitempty"`
	PhoneExtn            string                          `json:"PhoneExtn,omitempty"`
	Fax                  string                          `json:"Fax,omitempty"`
	BusinessType         string                          `json:"BusinessType,omitempty"`
	SigningAuthority     *CreateBusinessSigningAuthority `json:"SigningAuthority,omitempty"`
	KindOfEmployer       string                          `json:"KindOfEmployer,omitempty"`
	KindOfPayer          string                          `json:"KindOfPayer,omitempty"`
	IsBusinessTerminated bool                            `json:"IsBusinessTerminated"`
	Form1042SDetails     *CreateBusinessForm1042SDetails `json:"Form1042SDetails,omitempty"`
	IsForeign            bool                            `json:"IsForeign"`
	USAddress            *CreateBusinessUSAddress        `json:"USAddress,omitempty"`
	ForeignAddress       *CreateBusinessForeignAddress   `json:"ForeignAddress,omitempty"`
	ACADetails           *CreateBusinessACADetails       `json:"ACADetails,omitempty"`
}

type UpdateBusinessRequest struct {
	BusinessID           string                          `json:"BusinessId" binding:"required"`
	BusinessNm           string                          `json:"BusinessNm" binding:"required"`
	FirstNm              string                          `json:"FirstNm,omitempty"`
	MiddleNm             string                          `json:"MiddleNm,omitempty"`
	LastNm               string                          `json:"LastNm,omitempty"`
	Suffix               string                          `json:"Suffix,omitempty"`
	PayerRef             string                          `json:"PayerRef,omitempty"`
	TradeNm              string                          `json:"TradeNm,omitempty"`
	IsEIN                bool                            `json:"IsEIN"`
	EINorSSN             string                          `json:"EINorSSN" binding:"required"`
	IsDefaultBusiness    bool                            `json:"IsDefaultBusiness"`
	Email                string                          `json:"Email,omitempty"`
	ContactNm            string                          `json:"ContactNm,omitempty"`
	Phone                string                          `json:"Phone,omitempty"`
	PhoneExtn            string                          `json:"PhoneExtn,omitempty"`
	Fax                  string                          `json:"Fax,omitempty"`
	BusinessType         string                          `json:"BusinessType,omitempty"`
	SigningAuthority     *CreateBusinessSigningAuthority `json:"SigningAuthority,omitempty"`
	KindOfEmployer       string                          `json:"KindOfEmployer,omitempty"`
	KindOfPayer          string                          `json:"KindOfPayer,omitempty"`
	IsBusinessTerminated bool                            `json:"IsBusinessTerminated"`
	Form1042SDetails     *CreateBusinessForm1042SDetails `json:"Form1042SDetails,omitempty"`
	IsForeign            bool                            `json:"IsForeign"`
	USAddress            *CreateBusinessUSAddress        `json:"USAddress,omitempty"`
	ForeignAddress       *CreateBusinessForeignAddress   `json:"ForeignAddress,omitempty"`
	ACADetails           *CreateBusinessACADetails       `json:"ACADetails,omitempty"`
}

type ListBusinessesRequest struct {
	Page     int    `form:"Page"`
	PageSize int    `form:"PageSize"`
	FromDate string `form:"FromDate"`
	ToDate   string `form:"ToDate"`
}

type GetBusinessRequest struct {
	BusinessID string `form:"BusinessId"`
	TIN        string `form:"TIN"`
	PayerRef   string `form:"PayerRef"`
}

type WhCertificateListRequest struct {
	PayeeRef string `form:"PayeeRef" binding:"required"`
}

type RequestByURLRequest struct {
	PayerRef  string `form:"payerref" binding:"required"`
	FormType  string `form:"formtype" binding:"required"`
	ReturnURL string `form:"returnurl" binding:"required"`
	CancelURL string `form:"cancelurl" binding:"required"`
}

type DeleteBusinessRequest struct {
	BusinessID string `form:"BusinessId"`
	EinOrSSN   string `form:"EinOrSSN"`
	PayerRef   string `form:"PayerRef"`
}

type WhCertificateRequestByEmailSubmissionManifest struct {
	IsTINMatching bool `json:"IsTINMatching"`
}

type WhCertificateRequestByEmailRequester struct {
	BusinessID *string `json:"BusinessId"`
	PayerRef   *string `json:"PayerRef"`
	TIN        *string `json:"TIN"`
}

type WhCertificateRequestByEmailRecipient struct {
	PayeeRef string `json:"PayeeRef" binding:"required"`
	Name     string `json:"Name" binding:"required"`
	Email    string `json:"Email" binding:"required,email"`
}

type WhCertificateRequestByEmailRequest struct {
	SubmissionManifest *WhCertificateRequestByEmailSubmissionManifest `json:"SubmissionManifest" binding:"required"`
	Requester          *WhCertificateRequestByEmailRequester          `json:"Requester"`
	Recipients         []WhCertificateRequestByEmailRecipient         `json:"Recipients" binding:"required,min=1,dive"`
}

type ListRecipientParams struct {
	BusinessID  string
	RecipientID string
	Page        int
	PageSize    int
}