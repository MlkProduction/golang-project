package domain

import (
	"fmt"
	"regexp"

	core_error "github.com/MlkProduction/golang-project/internal/core/errors"
)

type User struct {
	ID      int
	Version int

	FullName    string
	PhoneNumber *string
}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}

}

func NewUserUnitialized(fullName string, phoneNumber *string) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)

}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName))

	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("Invalid `FullName` len: %d: %w", fullNameLength, core_error.ErrInvalidArgument)

	}

	if u.PhoneNumber != nil {
		phoneNumberLen := len([]rune(*u.PhoneNumber))
		if phoneNumberLen < 10 || phoneNumberLen > 15 {
			return fmt.Errorf("Invalid `Phone Number` len: %d: %w", phoneNumberLen, core_error.ErrInvalidArgument)
		}

		re := regexp.MustCompile(`^\+[0-9]+$`)

		if !re.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("Invalid `Phone Number` format: %w", core_error.ErrInvalidArgument)
		}
	}

	return nil

}

type UserPatch struct {
	FullName    Nullable[string]
	PhoneNumber Nullable[string]
}

func (p *UserPatch) Validate() error {

	if p.FullName.Set && p.FullName.Value == nil {
		return fmt.Errorf("Full Name not can be patched to null in the validate:%w", core_error.ErrInvalidArgument)
	}

	return nil
}

func (u *User) ApplyPatch(patch UserPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("Full Name not can be patched to null to apply patch: %w", err)
	}

	tmp := *u

	if patch.FullName.Set {
		tmp.FullName = *patch.FullName.Value
	}

	if patch.PhoneNumber.Set {
		tmp.PhoneNumber = patch.PhoneNumber.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate patched user: %w", err)
	}

	*u = tmp
	return nil

}
