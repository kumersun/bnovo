package repository

import "fmt"

type CountryRepository struct{}

func NewCountryRepository() *CountryRepository {
	return &CountryRepository{}
}

func (r *CountryRepository) getCountryByPhoneNumber(phone string) (string, error) {
	return fmt.Sprintf("{country_code} by phone number {%v}", phone), nil
}
