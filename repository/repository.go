package repository

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sdk/dto"
	"strconv"
	"strings"
	"time"
)

const (
	taxBanditsAuthURL            = "https://V1-tbs-oauth.stssprint.com/v2/tbsauth"
	createBusinessURL            = "https://v1-tbs-api.stssprint.com/v1.7.3/Business/Create"
	updateBusinessURL            = "https://v1-tbs-api.stssprint.com/v1.7.3/Business/update"
	getBusinessURL               = "https://v1-tbs-api.stssprint.com/v1.7.3/Business/get"
	listBusinessesURL            = "https://v1-tbs-api.stssprint.com/v1.7.3/Business/List"
	deleteBusinessURL            = "https://v1-tbs-api.stssprint.com/v1.7.3/Business/delete"
	requestByBusinessURL         = "https://v1-tbs-api.stssprint.com/v1.7.3/Business/requestbyurl"
	whCertificateGetURL          = "https://v1-tbs-api.stssprint.com/v1.7.3/WhCertificate/List"
	whCertificateRequestEmailURL = "https://v1-tbs-api.stssprint.com/v1.7.3/WhCertificate/RequestByEmail"
	recipientListURL             = "https://v1-tbs-api.stssprint.com/v1.7.3/Recipient/List"
)

type IRepository interface {
	GetAccessToken(context.Context) (any, error)
	CreateBusiness(context.Context, dto.CreateBusinessRequest) (any, error)
	UpdateBusiness(context.Context, dto.UpdateBusinessRequest) (any, error)
	GetBusiness(context.Context, dto.GetBusinessRequest) (any, error)
	ListBusinesses(context.Context, dto.ListBusinessesRequest) (any, error)
	ListWhCertificate(context.Context, string) (any, error)
	RequestByURL(context.Context, dto.RequestByURLRequest) (any, error)
	DeleteBusiness(context.Context, dto.DeleteBusinessRequest) (any, error)
	RequestByEmail(context.Context, dto.WhCertificateRequestByEmailRequest) (any, error)
	ListRecipient(context.Context, dto.ListRecipientParams) (any, error)
}

type repository struct {
	client    *http.Client
	clientID  string
	secret    string
	userToken string
}

func NewRepository() IRepository {
	return &repository{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		clientID:  strings.TrimSpace(os.Getenv("TBS_CLIENT_ID")),
		secret:    strings.TrimSpace(os.Getenv("TBS_SECRET")),
		userToken: strings.TrimSpace(os.Getenv("TBS_USER_TOKEN")),
	}
}

func (r *repository) GetAccessToken(ctx context.Context) (any, error) {
	rawResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	return rawResponse, nil
}

func (r *repository) fetchAccessToken(ctx context.Context) (dto.TaxBanditsAuthResponse, error) {
	if r.clientID == "" || r.secret == "" || r.userToken == "" {
		return dto.TaxBanditsAuthResponse{}, fmt.Errorf("missing TaxBandits credentials: set TBS_CLIENT_ID, TBS_SECRET, and TBS_USER_TOKEN")
	}

	jwtToken, err := r.buildJWT()
	if err != nil {
		return dto.TaxBanditsAuthResponse{}, fmt.Errorf("build jwt: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, taxBanditsAuthURL, nil)
	if err != nil {
		return dto.TaxBanditsAuthResponse{}, fmt.Errorf("create auth request: %w", err)
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authentication", jwtToken)

	statusCode, parsedBody, err := doJSONRequest[dto.TaxBanditsAuthResponse](r.client, req)
	if err != nil {
		return dto.TaxBanditsAuthResponse{}, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return parsedBody, fmt.Errorf("taxbandits auth returned status %d: %s", statusCode, stringifyBody(parsedBody))
	}

	if strings.TrimSpace(parsedBody.AccessToken) == "" {
		return parsedBody, fmt.Errorf("auth response did not include AccessToken")
	}

	return parsedBody, nil
}

func (r *repository) buildJWT() (string, error) {
	headerBytes, err := json.Marshal(dto.TaxBanditsJWTHeader{
		Alg: "HS256",
		Typ: "JWT",
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt header: %w", err)
	}

	payloadBytes, err := json.Marshal(dto.TaxBanditsJWTPayload{
		Iss: r.clientID,
		Sub: r.clientID,
		Aud: r.userToken,
		Iat: time.Now().Unix(),
	})
	if err != nil {
		return "", fmt.Errorf("marshal jwt payload: %w", err)
	}

	encodedHeader := base64.RawURLEncoding.EncodeToString(headerBytes)
	encodedPayload := base64.RawURLEncoding.EncodeToString(payloadBytes)
	signingInput := encodedHeader + "." + encodedPayload

	mac := hmac.New(sha256.New, []byte(r.secret))
	if _, err := mac.Write([]byte(signingInput)); err != nil {
		return "", fmt.Errorf("sign jwt: %w", err)
	}

	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return signingInput + "." + signature, nil
}

func doJSONRequest[T any](client *http.Client, req *http.Request) (int, T, error) {
	resp, err := client.Do(req)
	if err != nil {
		var zero T
		return 0, zero, fmt.Errorf("call upstream api: %w", err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		var zero T
		return 0, zero, fmt.Errorf("read upstream response: %w", err)
	}

	var parsedBody T
	if len(responseBody) == 0 {
		var zero T
		return resp.StatusCode, zero, fmt.Errorf("upstream API returned an empty response")
	}

	if err := json.Unmarshal(responseBody, &parsedBody); err != nil {
		var zero T
		return resp.StatusCode, zero, fmt.Errorf("unmarshal upstream response: %w", err)
	}

	return resp.StatusCode, parsedBody, nil
}

func stringifyBody(body any) string {
	jsonBody, err := json.Marshal(body)
	if err == nil {
		return string(jsonBody)
	}

	return fmt.Sprintf("%v", body)
}

func (r *repository) CreateBusiness(ctx context.Context, req dto.CreateBusinessRequest) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal create business payload: %w", err)
	}

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, createBusinessURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create business request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[dto.CreateBusinessResponse](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits create business returned status %d: %s", statusCode, stringifyBody(response))
	}

	if response.BusinessID == "" && strings.EqualFold(response.StatusName, "Success") {
		return response, fmt.Errorf("create business response did not include BusinessId")
	}

	if len(response.Errors) > 0 {
		return response, fmt.Errorf("taxbandits create business returned validation errors: %s", stringifyBody(response.Errors))
	}

	return response, nil
}

func (r *repository) GetBusiness(ctx context.Context, req dto.GetBusinessRequest) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	queryURL, err := url.Parse(getBusinessURL)
	if err != nil {
		return nil, fmt.Errorf("parse get business url: %w", err)
	}

	q := queryURL.Query()
	if strings.TrimSpace(req.BusinessID) != "" {
		q.Set("BusinessId", req.BusinessID)
	}
	if strings.TrimSpace(req.TIN) != "" {
		q.Set("TIN", req.TIN)
	}
	if strings.TrimSpace(req.PayerRef) != "" {
		q.Set("PayerRef", req.PayerRef)
	}
	queryURL.RawQuery = q.Encode()

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create get business request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits get business returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) ListWhCertificate(ctx context.Context, businessID string) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	queryURL, err := url.Parse(whCertificateGetURL)
	if err != nil {
		return nil, fmt.Errorf("parse whcertificate url: %w", err)
	}

	q := queryURL.Query()

	if strings.TrimSpace(businessID) != "" {
		q.Set("BusinessId", businessID) // ✅ THIS is the key line
	}

	queryURL.RawQuery = q.Encode()

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create whcertificate get request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits whcertificate get returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) ListBusinesses(ctx context.Context, req dto.ListBusinessesRequest) (any, error) {
	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	queryURL, err := url.Parse(listBusinessesURL)
	if err != nil {
		return nil, fmt.Errorf("parse list businesses url: %w", err)
	}

	q := queryURL.Query()
	q.Set("Page", strconv.Itoa(req.Page))
	q.Set("PageSize", strconv.Itoa(req.PageSize))
	if strings.TrimSpace(req.FromDate) != "" {
		q.Set("FromDate", req.FromDate)
	}
	if strings.TrimSpace(req.ToDate) != "" {
		q.Set("ToDate", req.ToDate)
	}
	queryURL.RawQuery = q.Encode()

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create list businesses request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits list businesses returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) UpdateBusiness(ctx context.Context, req dto.UpdateBusinessRequest) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal update business payload: %w", err)
	}

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPut, updateBusinessURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create update business request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits update business returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) RequestByURL(ctx context.Context, req dto.RequestByURLRequest) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	queryURL, err := url.Parse(requestByBusinessURL)
	if err != nil {
		return nil, fmt.Errorf("parse request by url: %w", err)
	}

	q := queryURL.Query()
	q.Set("payerref", req.PayerRef)
	q.Set("formtype", req.FormType)
	q.Set("returnurl", req.ReturnURL)
	q.Set("cancelurl", req.CancelURL)
	queryURL.RawQuery = q.Encode()

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodGet, queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create request by url request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits request by url returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) DeleteBusiness(ctx context.Context, req dto.DeleteBusinessRequest) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	queryURL, err := url.Parse(deleteBusinessURL)
	if err != nil {
		return nil, fmt.Errorf("parse delete business url: %w", err)
	}

	q := queryURL.Query()
	if strings.TrimSpace(req.BusinessID) != "" {
		q.Set("BusinessId", req.BusinessID)
	}
	if strings.TrimSpace(req.EinOrSSN) != "" {
		q.Set("EinOrSSN", req.EinOrSSN)
	}
	if strings.TrimSpace(req.PayerRef) != "" {
		q.Set("PayerRef", req.PayerRef)
	}
	queryURL.RawQuery = q.Encode()

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodDelete, queryURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("create delete business request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits delete business returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) RequestByEmail(ctx context.Context, req dto.WhCertificateRequestByEmailRequest) (any, error) {
	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request by email payload: %w", err)
	}

	upstreamReq, err := http.NewRequestWithContext(ctx, http.MethodPost, whCertificateRequestEmailURL, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("create whcertificate request by email request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Content-Type", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, err
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf("taxbandits whcertificate request by email returned status %d: %s", statusCode, stringifyBody(response))
	}

	return response, nil
}

func (r *repository) ListRecipient(
	ctx context.Context,
	params dto.ListRecipientParams,
) (any, error) {

	authResponse, err := r.fetchAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("fetch access token: %w", err)
	}

	queryURL, err := url.Parse(recipientListURL)
	if err != nil {
		return nil, fmt.Errorf("parse recipient list url: %w", err)
	}

	q := queryURL.Query()

	if strings.TrimSpace(params.BusinessID) != "" {
		q.Set("BusinessId", params.BusinessID)
	}

	if strings.TrimSpace(params.RecipientID) != "" {
		q.Set("RecipientId", params.RecipientID)
	}

	page := params.Page
	if page <= 0 {
		page = 1
	}
	q.Set("Page", strconv.Itoa(page))

	pageSize := params.PageSize
	if pageSize <= 0 {
		pageSize = 50
	}
	if pageSize > 100 {
		pageSize = 100
	}
	q.Set("PageSize", strconv.Itoa(pageSize))

	queryURL.RawQuery = q.Encode()

	upstreamReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		queryURL.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("create recipient list request: %w", err)
	}

	upstreamReq.Header.Set("Accept", "application/json")
	upstreamReq.Header.Set("Authorization", "Bearer "+authResponse.AccessToken)

	statusCode, response, err := doJSONRequest[map[string]any](r.client, upstreamReq)
	if err != nil {
		return nil, fmt.Errorf("recipient list request failed: %w", err)
	}

	if statusCode < http.StatusOK || statusCode >= http.StatusMultipleChoices {
		return response, fmt.Errorf(
			"taxbandits recipient list returned status %d: %s",
			statusCode,
			stringifyBody(response),
		)
	}

	return response, nil
}
