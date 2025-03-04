package service

import (
	"be-assessment-test/internal/repository"
	"be-assessment-test/internal/types"
	"be-assessment-test/internal/util/dbutil"
	"context"
	"net/http"
	"time"

	"github.com/go-errors/errors"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Registration interface {
	Create(ctx context.Context, req types.RegistrationCreateReq) (types.RegistrationCreateRes, error)
}

type registrationImpl struct {
	db              *sqlx.DB
	userRepo        repository.User
	bankAccountRepo repository.BankAccount
}

func NewRegistration(db *sqlx.DB, userRepo repository.User, bankAccountRepo repository.BankAccount) Registration {
	return &registrationImpl{
		db,
		userRepo,
		bankAccountRepo,
	}
}

func (s *registrationImpl) Create(ctx context.Context, req types.RegistrationCreateReq) (types.RegistrationCreateRes, error) {
	res := types.RegistrationCreateRes{}

	err := req.Validate()
	if err != nil {
		return res, err
	}

	userExists := true
	_, err = s.userRepo.FindByNationalIDNumberOrPhoneNumber(ctx, req.NationalIdentityNumber, req.PhoneNumber)
	if errors.Is(err, types.ErrNoData) {
		userExists = false
	} else if err != nil {
		return res, err
	}

	if userExists {
		return res, errors.New(types.AppErr{StatusCode: http.StatusBadRequest, Message: "national identity number or phone number already used"})
	}

	userID, err := uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}

	bankAccountID, err := uuid.NewV7()
	if err != nil {
		return res, errors.New(err)
	}

	timeNow := time.Now()

	user := types.User{
		ID:                     userID,
		Name:                   req.Name,
		NationalIdentityNumber: req.NationalIdentityNumber,
		PhoneNumber:            req.PhoneNumber,
		CreatedAt:              timeNow,
	}

	bankAccount := types.BankAccount{
		ID:        bankAccountID,
		UserID:    user.ID,
		CreatedAt: timeNow,
	}

	tx, err := dbutil.NewSqlxTx(ctx, s.db, nil)
	if err != nil {
		return res, errors.New(err)
	}

	defer tx.Rollback()

	err = s.userRepo.CreateTx(ctx, tx, user)
	if err != nil {
		return res, err
	}

	err = s.bankAccountRepo.CreateTx(ctx, tx, bankAccount)
	if err != nil {
		return res, err
	}

	res.BankAccountID = bankAccountID

	err = tx.Commit()
	if err != nil {
		return res, errors.New(err)
	}

	return res, nil
}
