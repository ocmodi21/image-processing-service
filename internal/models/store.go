package models

type Store struct {
	ID       string `db:"id" json:"id"`
	Name     string `db:"store_name" json:"store_name"`
	AreaCode string `db:"area_code" json:"area_code"`
}
