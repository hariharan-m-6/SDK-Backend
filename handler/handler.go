package handler

import (
	"net/http"
	"sdk/dto"
	"sdk/service"
	"strconv"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	GetAccessToken(*gin.Context)
	CreateBusiness(*gin.Context)
	UpdateBusiness(*gin.Context)
	GetBusiness(*gin.Context)
	ListBusinesses(*gin.Context)
	ListWhCertificate(*gin.Context)
	RequestByURL(*gin.Context)
	DeleteBusiness(*gin.Context)
	RequestByEmail(*gin.Context)
	ListRecipient(*gin.Context)
	GetPDF(*gin.Context)
}

type handler struct {
	service service.IService
}

func NewHandler(service service.IService) IHandler {
	return &handler{
		service: service,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func handleServiceError(c *gin.Context, err error, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   message,
		Details: err.Error(),
	})
}

func respondWithValidationError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, ErrorResponse{
		Error:   message,
		Details: "validation failed",
	})
}

func (h *handler) GetAccessToken(c *gin.Context) {
	response, err := h.service.GetAccessToken(c.Request.Context())
	if err != nil {
		handleServiceError(c, err, "Failed to get access token")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) CreateBusiness(c *gin.Context) {
	var req dto.CreateBusinessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithValidationError(c, "Invalid request body")
		return
	}

	response, err := h.service.CreateBusiness(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to create business")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateBusiness(c *gin.Context) {
	var req dto.UpdateBusinessRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithValidationError(c, "Invalid request body")
		return
	}

	response, err := h.service.UpdateBusiness(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to update business")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) GetBusiness(c *gin.Context) {
	req := dto.GetBusinessRequest{
		BusinessID: c.Query("BusinessId"),
		TIN:        c.Query("TIN"),
		PayerRef:   c.Query("PayerRef"),
	}

	if req.BusinessID == "" && req.TIN == "" && req.PayerRef == "" {
		respondWithValidationError(c, "At least one of BusinessId, TIN, or PayerRef is required")
		return
	}

	response, err := h.service.GetBusiness(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to fetch business")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) ListBusinesses(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("Page", "1"))
	if err != nil || page <= 0 {
		respondWithValidationError(c, "Invalid page parameter")
		return
	}

	pageSize, err := strconv.Atoi(c.DefaultQuery("PageSize", "10"))
	if err != nil || pageSize <= 0 {
		respondWithValidationError(c, "Invalid page size parameter")
		return
	}

	req := dto.ListBusinessesRequest{
		Page:     page,
		PageSize: pageSize,
		FromDate: c.Query("FromDate"),
		ToDate:   c.Query("ToDate"),
	}

	response, err := h.service.ListBusinesses(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to list businesses")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) ListWhCertificate(c *gin.Context) {
	businessID := c.Query("BusinessId")

	response, err := h.service.ListWhCertificate(c.Request.Context(), businessID)
	if err != nil {
		handleServiceError(c, err, "Failed to list certificates")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) RequestByURL(c *gin.Context) {
	req := dto.RequestByURLRequest{
		PayerRef:  c.Query("payerref"),
		FormType:  c.Query("formtype"),
		ReturnURL: c.Query("returnurl"),
		CancelURL: c.Query("cancelurl"),
	}

	if req.PayerRef == "" || req.FormType == "" || req.ReturnURL == "" || req.CancelURL == "" {
		respondWithValidationError(c, "Missing required query parameters")
		return
	}

	response, err := h.service.RequestByURL(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to generate request URL")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) DeleteBusiness(c *gin.Context) {
	req := dto.DeleteBusinessRequest{
		BusinessID: c.Query("BusinessId"),
		EinOrSSN:   c.Query("EinOrSSN"),
		PayerRef:   c.Query("PayerRef"),
	}

	if req.BusinessID == "" && req.EinOrSSN == "" && req.PayerRef == "" {
		respondWithValidationError(c, "At least one of BusinessId, EinOrSSN, or PayerRef is required")
		return
	}

	response, err := h.service.DeleteBusiness(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to delete business")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) RequestByEmail(c *gin.Context) {
	var req dto.WhCertificateRequestByEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithValidationError(c, "Invalid request body")
		return
	}

	response, err := h.service.RequestByEmail(c.Request.Context(), req)
	if err != nil {
		handleServiceError(c, err, "Failed to send email request")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) ListRecipient(c *gin.Context) {
	businessID := c.Query("BusinessId")
	recipientID := c.Query("RecipientId")

	page, _ := strconv.Atoi(c.DefaultQuery("Page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("PageSize", "50"))

	if pageSize > 100 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "PageSize cannot exceed 100",
		})
		return
	}

	params := dto.ListRecipientParams{
		BusinessID:  businessID,
		RecipientID: recipientID,
		Page:        page,
		PageSize:    pageSize,
	}

	response, err := h.service.ListRecipient(c.Request.Context(), params)
	if err != nil {
		handleServiceError(c, err, "Failed to list recipients")
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *handler) GetPDF(c *gin.Context) {
	key := c.Query("key")
	if key == "" {
		c.JSON(400, gin.H{"error": "key is required"})
		return
	}

	data, err := h.service.GetPDF(c.Request.Context(), key)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Type", "application/pdf")
	c.Header("Content-Disposition", "inline")
	c.Data(200, "application/pdf", data)
}