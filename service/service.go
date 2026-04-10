package service

import (
	"context"
	"sdk/dto"
	"sdk/repository"
)

type IService interface {
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
	GetPDF(ctx context.Context, key string) ([]byte, error)
}

type service struct {
	repo repository.IRepository
	s3   *S3Service
}

func NewService(repo repository.IRepository) IService {
	s3Service, _ := NewS3Service()

	return &service{
		repo: repo,
		s3:   s3Service,
	}
}

func (s *service) GetAccessToken(ctx context.Context) (any, error) {
	response, err := s.repo.GetAccessToken(ctx)
	if err != nil {
		return nil, err
	}

	rawResponse, ok := response.(dto.TaxBanditsAuthResponse)
	if !ok {
		return response, nil
	}

	return dto.TaxBanditsAuthAPIResponse{
		Status:      rawResponse.StatusName,
		Code:        rawResponse.StatusCode,
		Message:     rawResponse.StatusMessage,
		AccessToken: rawResponse.AccessToken,
	}, nil
}

func (s *service) CreateBusiness(ctx context.Context, req dto.CreateBusinessRequest) (any, error) {
	response, err := s.repo.CreateBusiness(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) UpdateBusiness(ctx context.Context, req dto.UpdateBusinessRequest) (any, error) {
	response, err := s.repo.UpdateBusiness(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) GetBusiness(ctx context.Context, req dto.GetBusinessRequest) (any, error) {
	response, err := s.repo.GetBusiness(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) ListBusinesses(ctx context.Context, req dto.ListBusinessesRequest) (any, error) {
	response, err := s.repo.ListBusinesses(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) ListWhCertificate(ctx context.Context, businessID string) (any, error) {
	response, err := s.repo.ListWhCertificate(ctx, businessID)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) RequestByURL(ctx context.Context, req dto.RequestByURLRequest) (any, error) {
	response, err := s.repo.RequestByURL(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) DeleteBusiness(ctx context.Context, req dto.DeleteBusinessRequest) (any, error) {
	response, err := s.repo.DeleteBusiness(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) RequestByEmail(ctx context.Context, req dto.WhCertificateRequestByEmailRequest) (any, error) {
	response, err := s.repo.RequestByEmail(ctx, req)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) ListRecipient(
	ctx context.Context,
	params dto.ListRecipientParams,
) (any, error) {

	response, err := s.repo.ListRecipient(ctx, params)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (s *service) GetPDF(ctx context.Context, key string) ([]byte, error) {
	return s.s3.GetFile(ctx, key)
}