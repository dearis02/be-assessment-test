package types

import (
	pkg "be-assessment-test/pkg/custom_validator"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

// region service types

type RegistrationCreateReq struct {
	Name                   string `json:"name"`
	NationalIdentityNumber string `json:"national_identity_number"`
	PhoneNumber            string `json:"phone_number"`
}

func (r RegistrationCreateReq) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Name, validation.Required),
		validation.Field(&r.NationalIdentityNumber, validation.Required, validation.Length(16, 16)),
		validation.Field(&r.PhoneNumber, validation.Required, validation.By(pkg.NumericRule)),
	)
}

type RegistrationCreateRes struct {
	BankAccountNumber string `json:"bank_account_number"`
}

// endregion service types
