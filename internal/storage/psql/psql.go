package psql

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

type People struct {
	Id         int    `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Surname    string `json:"surname"`
	Patronymic string `json:"patronymic,omitempty"`
	Country    string `json:"country,omitempty"`
	Gender     string `json:"gender,omitempty"`
	Age        int    `json:"age,omitempty"`
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.psql.New"

	connStr := storagePath
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SavePeople(name string, surname string, patronymic string, country string, gender string, age int) error {
	const op = "internal.storage.psql.SavePeople"

	stmt, err := s.db.Prepare("INSERT INTO peoples(name, surname, patronymic, country, gender, age) VALUES ($1, $2, $3, $4, $5, $6)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(name, surname, patronymic, country, gender, age)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = res

	return nil
}

func (s *Storage) DelPeople(id int) error {
	const op = "internal.storage.psql.DelPeople"

	stmt, err := s.db.Prepare("DELETE FROM peoples WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = res

	return nil
}

func (s *Storage) SearchId(id int) (People, error) {
	const op = "internal.storage.psql.SearchId"

	row := s.db.QueryRow("SELECT * FROM peoples WHERE id = $1", id)
	p := People{}
	err := row.Scan(&p.Id, &p.Name, &p.Surname, &p.Patronymic, &p.Country, &p.Gender, &p.Age)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	return p, nil
}

func (s *Storage) SearchName(name string, surname string) (People, error) {
	const op = "internal.storage.psql.SearchName"

	row := s.db.QueryRow("SELECT * FROM peoples WHERE (name = $1 AND surname = $2)", name, surname)
	p := People{}
	err := row.Scan(&p.Id, &p.Name, &p.Surname, &p.Patronymic, &p.Country, &p.Gender, &p.Age)
	if err != nil {
		fmt.Errorf("%s: %w", op, err)
	}

	return p, nil
}

func (s *Storage) Edit(id int, name string, surname string, patronymic string, country string, gender string, age int) error {
	const op = "internal.storage.psql.Edit"

	stmt, err := s.db.Prepare("UPDATE peoples SET name = $2, surname = $3, patronymic = $4, country = $5, gender = $6, age = $7 WHERE id = $1")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(id, name, surname, patronymic, country, gender, age)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_ = res

	return nil
}
