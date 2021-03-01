package transaction

import (
	"bwa-startup/campaign"
	"bwa-startup/payment"
	"errors"
)


type Service interface {
	GetTransactionByCampaignID(input GetTransactionByCampaignIdInput) ([]Transaction, error)
	GetTransactionByUserID(userID int) ([]Transaction, error)
	CreateTransaction(input CreateTransactionInput) (Transaction, error)
}

type service struct {
	repository       	Repository
	campaignRepository 	campaign.Repository
	paymentService		payment.Service
}

func NewService(repository Repository, campaignRepository campaign.Repository, 	paymentService payment.Service) *service {
	return &service{repository, campaignRepository, paymentService}
}

func (s *service) GetTransactionByCampaignID(input GetTransactionByCampaignIdInput) ([]Transaction, error) {
	campaign, err := s.campaignRepository.FindById(input.ID)
	if err != nil {
		return []Transaction{}, err
	}
	
	if campaign.UserID != input.User.ID {
		return []Transaction{}, errors.New("Not an owner of the campaign")
	}
	
	
	transactions, err := s.repository.GetByCampaignID(input.ID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) GetTransactionByUserID(userID int) ([]Transaction, error) {
	transactions, err := s.repository.GetByUserID(userID)
	if err != nil {
		return transactions, err
	}

	return transactions, nil
}

func (s *service) CreateTransaction(input CreateTransactionInput) (Transaction, error) {
	transaction := Transaction{}
	transaction.Amount = input.Amount
	transaction.CampaignID = input.CampaignID
	transaction.UserID = input.User.ID
	transaction.Status = "pending"

	newTransaction, err := s.repository.Save(transaction)
	if err != nil {
		return newTransaction, err
	}
	paymentTransaction := payment.Transaction {
		ID: newTransaction.ID,
		Amount: newTransaction.Amount,
	}
	
	paymentURL, err := s.paymentService.GetPaymentURL(paymentTransaction, input.User)
	if err != nil {
		return newTransaction, err
	}
	
	newTransaction.PaymentURL = paymentURL

	newTransaction, err = s.repository.Update(transaction)
	if err != nil {
		return newTransaction, err
	}

	return newTransaction, nil
}
