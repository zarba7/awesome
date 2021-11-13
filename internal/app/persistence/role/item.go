package role

import (
	"pbRole"
	store "src/_infrastructure/store/mysql"
)

type addItem struct {
	pbRole.EventAddItem
}

func (a *addItem) Write(db *store.MySQLClient) error {
	panic("implement me")
}
