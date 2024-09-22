package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kumersun/bnovo/entity"
)

var countryRepo = NewCountryRepository()

type GuestRepository struct {
	dbpool *pgxpool.Pool
}

func NewGuestRepository(dbpool *pgxpool.Pool) *GuestRepository {
	return &GuestRepository{dbpool}
}

func prepareGuest(guest *entity.Guest) {
	if guest.Country == "" {
		country, err := countryRepo.getCountryByPhoneNumber(guest.Phone)
		if err != nil {
			guest.Country = err.Error()
		} else {
			guest.Country = country
		}
	}
}

func (gr *GuestRepository) CreateGuest(guest *entity.Guest) error {
	query := "INSERT INTO guest (name, surname, phone, email, country) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	return gr.dbpool.QueryRow(context.Background(), query, guest.Name, guest.Surname, guest.Phone, guest.Email, guest.Country).
		Scan(&guest.ID)
}

func (gr *GuestRepository) GetGuests(ctx context.Context) ([]entity.Guest, error) {
	rows, err := gr.dbpool.Query(
		ctx,
		"SELECT id, name, surname, phone, email, country FROM guest",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var guests []entity.Guest
	for rows.Next() {
		var guest entity.Guest
		if err := rows.Scan(&guest.ID, &guest.Name, &guest.Surname, &guest.Phone, &guest.Email, &guest.Country); err != nil {
			return nil, err
		}
		prepareGuest(&guest)
		guests = append(guests, guest)
	}

	return guests, nil
}

func (gr *GuestRepository) GetGuest(id int) (*entity.Guest, error) {
	var guest entity.Guest
	query := "SELECT id, name, surname, phone, email, country FROM guest WHERE id = $1"
	err := gr.dbpool.QueryRow(context.Background(), query, id).
		Scan(&guest.ID, &guest.Name, &guest.Surname, &guest.Phone, &guest.Email, &guest.Country)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	prepareGuest(&guest)
	return &guest, nil
}

func (gr *GuestRepository) UpdateGuest(guest *entity.Guest) error {
	query := "UPDATE guest SET name = $1, surname = $2, phone = $3, email = $4, country = $5 WHERE id = $6"
	_, err := gr.dbpool.Exec(
		context.Background(),
		query,
		guest.Name,
		guest.Surname,
		guest.Phone,
		guest.Email,
		guest.Country,
		guest.ID,
	)
	return err
}

func (gr *GuestRepository) DeleteGuest(id int) error {
	query := "DELETE FROM guest WHERE id = $1"
	_, err := gr.dbpool.Exec(context.Background(), query, id)
	return err
}
