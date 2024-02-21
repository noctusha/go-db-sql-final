package main

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return 0, err
	}
	defer db.Close()

	res, err := db.Exec("insert into parcel (client, status, address, created_at) values (:client, :status, :address, :created_at)",
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}

	l, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(l), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	p := Parcel{}

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return p, err
	}
	defer db.Close()

	row := db.QueryRow("select * from parcel where number = :number", sql.Named("number", number))

	err = row.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
	if err != nil {
		return p, err
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	var res []Parcel

	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return res, err
	}
	defer db.Close()

	rows, err := db.Query("select * from parcel where client = :client", sql.Named("client", client))
	if err != nil {
		return res, err
	}
	defer rows.Close()

	for rows.Next() {
		p := Parcel{}
		er := rows.Scan(&p.Number, &p.Client, &p.Status, &p.Address, &p.CreatedAt)
		if er != nil {
			return res, err
		}
		res = append(res, p)
	}
	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("update parcel set status = :status where number = :number", sql.Named("status", status), sql.Named("number", number))
	if err != nil {
		return err
	}
	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("update parcel set address = :address where (number = :number and status = :status)",
		sql.Named("address", address),
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) Delete(number int) error {
	db, err := sql.Open("sqlite", "tracker.db")
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("delete from parcel where (number = :number and status = :status)",
		sql.Named("number", number),
		sql.Named("status", ParcelStatusRegistered),
	)
	if err != nil {
		return err
	}

	return nil
}
