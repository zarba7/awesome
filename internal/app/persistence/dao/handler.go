package dao

import (
	store "src/_infrastructure/store/mysql"
)



type Handler interface {
	Write(db *store.MySQLClient) error
}

type Table struct {
	hash map[int32]Handler
}

func (t *Table)Register()  {

}

func (t *Table)Find(id int32)Handler  {

}

