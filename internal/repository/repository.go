// Package repository
package repository

// Repository interface
type Repository interface {
	Add(id string, value string) error
	Get(id string) (value string, error error)
}
