package persistence

import (
	"base/db"
	"ddd/adaptor"
	"ddd/po"
	"fmt"
	"github.com/pkg/errors"
)

func NewMySQLRepo(dns db.MySQLSourceName) (repo *mysqlRepo, err error)  {
	repo = &mysqlRepo{}
	if repo.MySQL, err = db.NewMySQL(dns); err != nil {
		return
	}
	var sql = fmt.Sprintf("CREATE DATABASE IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci; user `%s`", dns.Schemas, dns.Schemas)
	if err = repo.Exec(sql).Error; err != nil {
		return
	}
	tabs := []po.Model{&po.Event{}, &po.Command{}, &po.Aggregate{}}
	for _, tb := range tabs {
		if !repo.HasTable(tb) {
			if err = repo.CreateTable(tb).Error; err != nil {
				return nil, errors.WithMessagef(err, "建表")
			}
			tb.Index(repo.MySQL)
		}
	}
	return
}

type mysqlRepo struct {
	*db.MySQL
}

const sqlL = "select * from command where aggregate_id = %d and version > %d and version <= %d"

func (repo *mysqlRepo) Load(agg *po.DomainAggregate) error {
	var curr = agg.CurrVersion
	if agg.SnapshotVersion == 0 {
		var tb = &agg.Aggregate
		if err := repo.First(tb, agg.ID).Error; err != nil {
			return err
		}
	}
	if agg.SnapshotVersion >= curr {
		return nil
	}
	var sql = fmt.Sprintf(sqlL, agg.ID, agg.SnapshotVersion, agg.CurrVersion)
	rows, err := repo.Exec(sql).Rows()
	if err != nil {
		return errors.WithMessagef(err, "%v", *agg)
	}
	for rows.Next() {
		var result po.Command
		if err := repo.ScanRows(rows, &result); err != nil {
			return errors.WithMessagef(err, "%v", *agg)
		}
		agg.Events = append(agg.Events, &result)
	}
	return nil
}

const sqlS1 = "update aggregate set curr_version = %d where id = %d and curr_version = %d"

func (repo *mysqlRepo) Save(agg *po.DomainAggregate) error {
	var sql = fmt.Sprintf(sqlS1, agg.CurrVersion, agg.ID, agg.SnapshotVersion)
	var err error
	if agg.SnapshotVersion == 0 {
		var tb = &agg.Aggregate
		if err = repo.FirstOrCreate(tb, agg.ID).Error; err != nil {
			return err
		}
	}
	var commands []interface{}
	for _, e := range agg.Events {
		commands = append(commands, e)
	}
	var events []interface{}
	for _, e := range agg.Events {
		events = append(events, e)
	}

	return repo.Transaction(func(db *db.MySQL) error {
		var result = db.Exec(sql)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.Errorf("乐观锁版本冲突 %v", *agg)
		}
		if err = db.BatchInsert(commands); err != nil{
			return err
		}
		return db.BatchInsert(events)
	})
}
