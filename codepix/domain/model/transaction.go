package model

import (
	"errors"
	"time"

	"github.com/asaskevich/govalidator"
	uuid "github.com/satori/go.uuid"
)

const (
	TransactionPending   string = "pending"
	TransactionCompleted string = "completed"
	TransactionError     string = "error"
	TransactionConfirmed string = "confirmed"
)

type TransactionRepositoryInterface interface {
	Register(transaction *Transaction) error
	Save(transaction *Transaction) error
	Find(id string) (*Transaction, error)
}

type Transactions struct {
	Transaction []Transaction
}

type Transaction struct {
	Base              `valid:"required"`
	AccountFrom       *Account `valid:"-"`
	Amount            float64  `json:"amount"  valid:"notnull"`
	PixKeyTo          *PixKey  `valid:"-"`
	Status            string   `json:"status" valid:"notnull"`
	Description       string   `json:"description" valid:"notnull"`
	CancelDescription string   `json:"cancel_description" valid:"notnull"`
}

func (tr *Transaction) isValid() error {
	_, err := govalidator.ValidateStruct(tr)

	if tr.Amount <= 0 {
		return errors.New("The amount must be grater than 0")
	}

	if tr.Status != TransactionPending && tr.Status != TransactionCompleted && tr.Status != TransactionError {
		return errors.New("invalid status from this transaction")
	}

	if tr.PixKeyTo.AccountID == tr.AccountFrom.ID {
		return errors.New("the source and destination accoun not be the same")
	}

	if err != nil {
		return err
	}

	return nil
}

func NewTransaction(accountFrom *Account, amount float64, pixKeyTo *PixKey, description string) (*Transaction, error) {
	tr := Transaction{
		AccountFrom: accountFrom,
		Amount:      amount,
		PixKeyTo:    pixKeyTo,
		Status:      TransactionPending,
		Description: description,
	}

	tr.ID = uuid.NewV4().String()
	tr.CreatedAt = time.Now()

	err := tr.isValid()

	if err != nil {
		return nil, err
	}

	return &tr, nil
}

func (tr *Transaction) Complete() error {
	tr.Status = TransactionCompleted
	tr.UpdatedAt = time.Now()
	err := tr.isValid()

	return err
}

func (tr *Transaction) Confirm() error {
	tr.Status = TransactionConfirmed
	tr.UpdatedAt = time.Now()
	err := tr.isValid()

	return err
}

func (tr *Transaction) Cancel(description string) error {
	tr.Status = TransactionError
	tr.Description = description
	tr.UpdatedAt = time.Now()
	err := tr.isValid()

	return err
}
