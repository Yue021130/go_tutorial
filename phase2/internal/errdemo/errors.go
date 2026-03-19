package errdemo

import (
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidUserID = errors.New("invalid user id")

type NotFoundError struct {
	Resource string
	ID       int
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("%s id=%d not found", e.Resource, e.ID)
}

func ValidateUserID(id int) error {
	if id <= 0 {
		return fmt.Errorf("validate user id %d: %w", id, ErrInvalidUserID)
	}
	return nil
}

func FindUserName(id int) (string, error) {
	if err := ValidateUserID(id); err != nil {
		return "", err
	}

	switch id {
	case 1:
		return "alice", nil
	case 2:
		return "bob", nil
	default:
		return "", fmt.Errorf("find user: %w", NotFoundError{
			Resource: "user",
			ID:       id,
		})
	}
}

func LoadUserDisplayName(id int) (string, error) {
	name, err := FindUserName(id)
	if err != nil {
		return "", fmt.Errorf("load user display name: %w", err)
	}
	return strings.ToUpper(name), nil
}

func IsInvalidID(err error) bool {
	return errors.Is(err, ErrInvalidUserID)
}

func IsNotFound(err error) bool {
	var nf NotFoundError
	return errors.As(err, &nf)
}
