// Package repository
package repository

const (
	// ErrorIdNotFound error message
	ErrorIdNotFound = `id not found`

	// ErrorValueAlreadyExist error message
	ErrorValueAlreadyExist = `value already exist`
)

// Repository interface
type Repository interface {
	Add(id string, value string) error
	Get(id string) (value string, error error)
}
