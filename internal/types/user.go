package types

import (
	"time"

	"github.com/google/uuid"
)

// region repo types

type User struct {
	ID                     uuid.UUID `db:"id"`
	Name                   string    `db:"name"`
	NationalIdentityNumber string    `db:"national_identity_number"`
	PhoneNumber            string    `db:"phone_number"`
	CreatedAt              time.Time `db:"created_at"`
}

// endregion repo types
