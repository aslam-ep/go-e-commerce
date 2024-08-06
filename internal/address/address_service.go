package address

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/aslam-ep/go-e-commerce/config"
	"github.com/aslam-ep/go-e-commerce/utils"
)

type AddressService interface {
	CreateAddress(c context.Context, req *AddressReq) (*Address, error)
	GetAllAddress(c context.Context, userID int) (*AddressRes, error)
	GetAddressByID(c context.Context, id int, userID int) (*Address, error)
	UpdateAddress(c context.Context, req *AddressReq) (*Address, error)
	SetDefaultAddress(c context.Context, id int, userID int) (*utils.MessageRes, error)
	DeleteAddress(c context.Context, id int, userID int) (*utils.MessageRes, error)
}

type addressService struct {
	repository   AddressRepository
	timeout      time.Duration
	addressLimit int
}

func NewAddressService(addressRepo AddressRepository) AddressService {
	return &addressService{
		repository:   addressRepo,
		timeout:      time.Duration(config.AppConfig.DBTimeout) * time.Second,
		addressLimit: 10,
	}
}

func (s *addressService) CreateAddress(c context.Context, req *AddressReq) (*Address, error) {
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

func (s *addressService) GetAllAddress(c context.Context, userID int) (*AddressRes, error) {
	ctx, cancel := context.WithTimeout(c, s.timeout)
	defer cancel()

	addresses, err := s.repository.GetAll(ctx, userID)
	if err != nil {
		return nil, err
	}

	res := &AddressRes{
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

func (s *addressService) UpdateAddress(c context.Context, req *AddressReq) (*Address, error) {
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
