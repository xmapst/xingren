package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"reflect"
	"time"
	"xingren/utils"
)

var DB *gorm.DB

type Model struct {
	CreatedOn  int64 `gorm:"column:created_on" json:"created_on"`
	ModifiedOn int64 `gorm:"column:modified_on" json:"modified_on"`
}

const (
	PREFIX_COL_STR = "tb_"
)

func Setup() {
	var (
		err error
	)

	// /gorm.db
	DB, err = gorm.Open("sqlite3", "xingren.db")
	if err != nil {
		panic(err)
	}
	//defer DB.Close()
	DB.LogMode(false)
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return PREFIX_COL_STR + defaultTableName
	}

	DB.SingularTable(true)
	DB.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
	DB.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
	DB.DB().SetMaxIdleConns(10)
	DB.DB().SetMaxOpenConns(100)
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok {
			if createTimeField.IsBlank {
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifiedOn` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok {
		scope.SetColumn("ModifiedOn", time.Now().Unix())
	}
}

// addExtraSpaceIfExist adds a separator
func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}

// SelectData 查询数据
// listPtr struct的slice的指针
// condit map[string]interface{} 条件 id=? 1
// page int 页码，从0开始累加，获取所有传-1
// count int 总数
// isTotal bool 是否返回总条数
// groupBy string 分组
// orderBy string 排序
func SelectData(listPtr interface{}, where interface{}, page, count int64, groupBy, orderBy string) (err error) {
	db := DB.Model(reflect.New(reflect.ValueOf(listPtr).Elem().Type().Elem().Elem()).Interface())
	db, err = utils.BuildWhere(db, where)
	if err != nil {
		return err
	}

	if page > -1 { // 分页
		db = db.Limit(count).Offset(page * count)
	}

	if groupBy != "" { // 分组
		db = db.Group(groupBy)
	}

	if orderBy == "" { // 排序
		orderBy = "id desc"
	}

	db = db.Order(orderBy)

	if err = db.Find(listPtr).Error; err != nil {
		return err
	}

	return
}

// SelectOne 查询单条数据
// where map[string]interface{} 条件 id=? 1
// objPtr struct的指针
func SelectOne(objPtr interface{}, where interface{}) (err error) {
	db := DB.Model(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface())
	db, err = utils.BuildWhere(db, where)
	if err != nil {
		return err
	}
	return db.First(objPtr).Error
}

// UpsetData 存在则更新，否则插入
// where map[string]interface{} 条件 id=? 1
// updateField 更新字段名列表
// isNil 去除0, "", false字段
func UpSetData(objPtr interface{}, where interface{}, updateField []string, isNil bool) (err error) {
	return DBBegin(func(gormDB *gorm.DB) error {
		db := gormDB.Model(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface())
		db, err = utils.BuildWhere(db, where)
		if err != nil {
			return err
		}
		if err := db.First(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface()).Error; gorm.IsRecordNotFoundError(err) || where == nil { // 不存在或者条件为空时直接插入
			// 不存在
			return db.Create(objPtr).Error
		} else if err != nil {
			return err
		}

		// 存在则更新
		updataDataMap, err := utils.StructToMapFilterFields(objPtr, updateField)
		if err != nil {
			return err
		}
		if isNil {
			if err := utils.MapFullStruct(updataDataMap, reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface()); err != nil {
				return err
			}
			return db.Updates(objPtr).Error
		} else {
			return db.Updates(updataDataMap).Error
		}
	})
}

// UpdateData 更新数据
func UpdateData(objPtr interface{}, where interface{}, updateMap map[string]interface{}) (err error) {
	return TxBegin(DB, func(gormDB *gorm.DB) (err error) {

		db := gormDB.Model(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface())
		db, err = utils.BuildWhere(db, where)
		if err != nil {
			return err
		}

		if err := db.First(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface()).Error; gorm.IsRecordNotFoundError(err) {
			// 不存在
			return err
		}

		return db.Updates(updateMap).Error
	})
}

func DBBegin(transaction func(gormDB *gorm.DB) error) (err error) {
	db := DB.Begin()
	defer func() {
		if err != nil {
			db.Rollback()
		}
	}()
	if err = transaction(db); err != nil {
		return err
	}
	return db.Commit().Error
}

func TxBegin(db *gorm.DB, txFun func(gormDB *gorm.DB) error) (err error) {
	if nil == db {
		db = DB
	}
	if tx := db.Begin(); nil == tx.Error {
		db = tx
		defer func() {
			if nil != err {
				tx.Rollback()
			} else {
				err = tx.Commit().Error
			}
		}()
	}
	if err = txFun(db); err != nil {
		return err
	}
	return
}

// Exit 是否存在
func Exit(objPtr interface{}, where interface{}) (bool, error) {
	var err error
	db := DB.Model(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface())
	db, err = utils.BuildWhere(db, where)
	if err != nil {
		return false, err
	}

	if err := db.First(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface()).Error; gorm.IsRecordNotFoundError(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func Total(objPtr interface{}, where interface{}) (int64, error) {
	var err error
	db := DB.Model(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface())
	db, err = utils.BuildWhere(db, where)
	if err != nil {
		return 0, err
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func DelData(objPtr interface{}, where interface{}) (err error) {
	db := DB.Model(reflect.New(reflect.ValueOf(objPtr).Elem().Type()).Interface())
	db, err = utils.BuildWhere(db, where)
	if err != nil {
		return err
	}
	return db.Delete(objPtr).Error
}
