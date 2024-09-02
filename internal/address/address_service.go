package address

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aslam-ep/go-e-commerce/config"
	"github.com/aslam-ep/go-e-commerce/utils"
)

// Service interface defines the methods required for address services.
type Service interface {
	// CreateAdddress Creates a new address based on the provided request and returns the created address details.
	CreateAddress(c context.Context, req *CreateUpdateAddressRequest) (*Address, error)

	// GetAllAddress Get all address based on the given userID and returns the AdderessRes
	GetAllAddress(c context.Context, userID int) (*ListAddressRes, error)

	// GetAddressByID Get a address by the address and user ID and return it
	GetAddressByID(c context.Context, id int, userID int) (*Address, error)

	// UpdateAddress Update the address based on user request and returns the updated address
	UpdateAddress(c context.Context, req *CreateUpdateAddressRequest) (*Address, error)

	// SetDefaultAddress Make the given address ID and return status
	SetDefaultAddress(c context.Context, id int, userID int) (*utils.MessageRes, error)

	// DeleteAddress Delete address based on given ID
	DeleteAddress(c context.Context, id int, userID int) (*utils.MessageRes, error)
}

type addressService struct {
	repository   Repository
	timeout      time.Duration
	addressLimit int
}

// NewService creates a new instance of the address service.
func NewService(addressRepo Repository) Service {
	return &addressService{
		repository:   addressRepo,
		timeout:      time.Duration(config.AppConfig.DBTimeout) * time.Second,
		addressLimit: 10,
	}
}

func (s *addressService) CreateAddress(c context.Context, req *CreateUpdateAddressRequest) (*Address, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	count, err := s.repository.GetCountByUserID(ctx, int(req.UserID))
	if err != nil {
		return nil, err
	}

	if count >= s.addressLimit {
		return nil, errors.New("user can't have more than 10 addresses")
	}

	a := &Address{
		UserID:       req.UserID,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		PostalCode:   req.PostalCode,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
	}

	address, err := s.repository.Create(ctx, a)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressService) GetAllAddress(c context.Context, userID int) (*ListAddressRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	addresses, err := s.repository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &ListAddressRes{
		Count:     len(*addresses),
		Addresses: addresses,
	}

	return res, nil
}

func (s *addressService) GetAddressByID(c context.Context, id int, userID int) (*Address, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	address, err := s.repository.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	return address, nil
}

func (s *addressService) UpdateAddress(c context.Context, req *CreateUpdateAddressRequest) (*Address, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	_, err := s.repository.GetByID(ctx, int(req.ID), int(req.UserID))
	if err != nil {
		return nil, err
	}

	a := &Address{
		ID:           req.ID,
		UserID:       req.UserID,
		AddressLine1: req.AddressLine1,
		AddressLine2: req.AddressLine2,
		PostalCode:   req.PostalCode,
		City:         req.City,
		State:        req.State,
		Country:      req.Country,
	}

	updatedAddress, err := s.repository.Update(ctx, a)
	if err != nil {
		return nil, err
	}

	return updatedAddress, nil
}

func (s *addressService) SetDefaultAddress(c context.Context, id int, userID int) (*utils.MessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	_, err := s.repository.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	err = s.repository.SetDefault(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	res := &utils.MessageRes{
		Success: true,
		Message: fmt.Sprintf("Address(%d) set as the default address.", id),
	}

	return res, nil
}

func (s *addressService) DeleteAddress(c context.Context, id int, userID int) (*utils.MessageRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	_, err := s.repository.GetByID(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	err = s.repository.Delete(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	res := &utils.MessageRes{
		Success: true,
		Message: fmt.Sprintf("Address(%d) deleted.", id),
	}

	return res, nil
}
