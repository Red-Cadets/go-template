package models

import (
	"fmt"

	"gorm.io/gorm"
)

type Pet struct {
	gorm.Model
	Name string `gorm:"unique"`
	Type PetType
	Age  int
}

func (p Pet) String() string {
	return fmt.Sprintf("name: %s; type: %s; age: %d", p.Name, p.Type.String(), p.Age)
}

type PetType int

const (
	Dog PetType = iota
	Cat
	Parrot
)

func (t PetType) String() string {
	return [...]string{"dog", "cat", "parrot"}[t]
}
