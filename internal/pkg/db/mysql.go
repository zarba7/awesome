package db

import (
	"bytes"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"base/log"
	"time"
)

type MySQL struct {
	*gorm.DB
}

type MySQLSourceName struct {
	Addr     string
	Schemas  string
	User     string
	Password string
}

func (dsn *MySQLSourceName) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s)?charset=utf8&parseTime=True&loc=Local", dsn.User, dsn.Password, dsn.Addr)
}

func NewMySQL(dsn MySQLSourceName) (*MySQL, error) {
	schema, err := gorm.Open("mysql", dsn.String())
	schema.SingularTable(true)
	schema.DB().SetConnMaxLifetime(10 * time.Minute)
	schema.DB().SetMaxIdleConns(8)
	schema.DB().SetMaxOpenConns(30)
	log.Infof("connect to mysql by dsn | %s", dsn)
	return &MySQL{schema}, err
}

func (the *MySQL) PessimisticLock(row interface{}, primaryKey interface{}, exec func(db *MySQL) error) (err error) {
	return the.Transaction(func(db *MySQL) error {
		if err := db.Set("gorm:query_option", "FOR UPDATE").First(row, primaryKey).Error; err != nil {
			return err
		}
		return exec(db)
	})
}

func (the *MySQL) Transaction(fun func(db *MySQL) error) (err error) {
	// 开始事务
	tx := the.Begin()
	if tx.Error != nil {
		return tx.Error
	}
	defer func() {
		e := recover()
		if e != nil {
			log.Errorf("Tran panic %s", err)
			err = fmt.Errorf("%v", e)
			tx.Rollback()
			return
		}
		if !errors.Is(err, nil) {
			// 发生错误时回滚事务
			log.Error("Roll back", err)
			tx.Rollback()
			return
		}
		// 或提交事务
		err = tx.Commit().Error
	}()
	return fun(&MySQL{tx})
}

func (the *MySQL) NewRecord(rec interface{}) *MySQL {
	if !the.DB.NewRecord(rec) {
		the.Error = fmt.Errorf("%s pk exist", rec)
	}
	return the
}

func (the *MySQL) BatchInsert(rows []interface{}) error {
	if len(rows) == 0 {
		return nil
	}
	sql := the.genInsertSql(rows)
	return the.Exec(sql).Error
}

//清空所有表
func (the *MySQL) TruncateAllTable() error {
	var tabs []string
	err := the.DB.Raw("select database()").Pluck("database()", &tabs).Error
	if err != nil {
		return err
	}
	name := fmt.Sprintf("Tables_in_%s", tabs[0])
	err = the.DB.Raw("show tables").Pluck(name, &tabs).Error
	if err != nil {
		return err
	}

	return the.Transaction(func(db *MySQL) error {
		for _, ss := range tabs {
			if err := db.Exec(fmt.Sprintf("truncate %s", ss)).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (the *MySQL) genInsertSql(source []interface{}) string {
	scope := the.NewScope(source[0])
	var buffer bytes.Buffer
	sql := fmt.Sprintf("INSERT INTO %s VALUES ", scope.TableName())
	buffer.WriteString(sql)
	for i, e := range source {
		fields := the.NewScope(e).Fields()
		val := ""
		for j, field := range fields {
			if j == 0 {
				val += "("
			}
			v := field.Field.Interface()
			switch vv := v.(type) {
			case time.Time:
				val += fmt.Sprintf("'%s'", vv.Format("2006-01-02 15:04:05"))
			case float64, float32, string:
				val += fmt.Sprintf("'%v'", vv)
			default:
				val += fmt.Sprintf("'%d'", vv)
			}
			if j == len(fields)-1 {
				val += ")"
			} else {
				val += ","
			}
		}
		if i == len(source)-1 {
			val += ";"
		} else {
			val += ","
		}
		buffer.WriteString(val)
	}
	return buffer.String()
}
