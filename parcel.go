package main

import (
	"database/sql"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p

	query := "INSERT INTO parcel (client, status, address, created_at) VALUES (:client, :status, :address, :created_at)"
	result, err := s.db.Exec(query,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt))

	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()

	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	p := Parcel{}

	query := "SELECT number, client, status, address, created_at FROM parcel WHERE number = :number"
	result := s.db.QueryRow(query, sql.Named("number", number))

	// заполните объект Parcel данными из таблицы
	err := result.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

	return p, err
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	query := "SELECT number, client, status, address, created_at FROM parcel WHERE client = :client"

	rows, err := s.db.Query(query, sql.Named("client", client))

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// заполните срез Parcel данными из таблицы
	var res []Parcel
	for rows.Next() {
		p := Parcel{}

		err = rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)

		if err != nil {
			return nil, err
		}

		res = append(res, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := "UPDATE parcel SET status = :status WHERE number = :id"

	_, err := s.db.Exec(query,
		sql.Named("status", status),
		sql.Named("id", number))

	return err
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	query := "UPDATE parcel SET address = :address WHERE number = :id AND status = :status"

	_, err := s.db.Exec(query,
		sql.Named("address", address),
		sql.Named("id", number),
		sql.Named("status", ParcelStatusRegistered))

	return err
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	query := "DELETE FROM parcel WHERE number = :id AND status = :status"

	_, err := s.db.Exec(query,
		sql.Named("id", number),
		sql.Named("status", ParcelStatusRegistered))

	return err
}
