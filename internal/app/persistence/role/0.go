package role

import (
	"pbRole"
	"src/persistence/dao"
)

var Events = dao.Table{
	pbRole.EventOpt_addItem:&addItem{},
	pbRole.EventOpt_setName:&setName{},
}