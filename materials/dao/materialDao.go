package dao

import (
	"fmt"
	"genMaterials/db"
	"genMaterials/log"
	"genMaterials/model"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

var appLog *log.AppLogger

type materialDao struct{}

type MaterialDaoIF interface {
	AddMaterial(newMaterial []*model.Material) ([]model.Material, error)
	GetAllMaterials() ([]model.Material, error)
	GetMaterialsOnGroup(groupCode string) ([]model.Material, error)
	GetMaterialGroups() ([]model.Material, error)
	GetNestedMaterials(pid *model.Pids) ([]model.MaterialsObj, error)
	GetAllMaterialWithTrades() ([]*model.MaterialBom, error)
	GetMaterialWithTrades(tradeList model.TradeList) ([]*model.MaterialBom, error)
	GetTopMaterialWithTrades(tradeList model.TradeList) ([]*model.MaterialBom, error)
	GetTradesCategories(tradeName string) ([]model.TradeCategoriesList, error)
	GetAllCategories() ([]model.TradeCategoriesList, error)
	SearchMaterialExist(searchObj model.SearchMaterial) (int, error)
	SearchAndGetMaterialExist(searchObj []string) ([]model.Material, error)
}

func MaterialDao() MaterialDaoIF {
	return &materialDao{}
}

func (self *materialDao) AddMaterial(materialsList []*model.Material) ([]model.Material, error) {

	materials := []model.Material{}

	// db, errs := db.DBConnect()
	db, ConnectionErrs := db.SqlxConnect()
	if ConnectionErrs != nil {
		return nil, ConnectionErrs
	}
	var err error

	sqlStatement := `INSERT INTO materials (pid,material_code,material_group,discription,material_name,material_unit,base_price,gst,created_at,created_by,last_updated_by,last_updated_at,sort_seq) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13) RETURNING id`
	fmt.Println(sqlStatement)
	for _, material := range materialsList {
		var materialId int

		errSql := db.Get(&materialId, sqlStatement, material.Pid, material.MaterialCode, material.MaterialGroupName, material.Discription,
			material.MaterialName, material.MaterialUnit, material.BasePrice, material.Gst, material.CreatedAt, material.CreatedBy, material.LastUpdatedBy, material.LastUpdatedAt, material.SortSeq)
		if errSql != nil {
			fmt.Println(errSql.Error())
			return materials, errSql
		}
		if materialId > 0 {
			var tradeId int
			sqlStatement1 := `INSERT INTO material_trade(
	trade_id, material_id, created_by, created_at, last_updated_by, last_updated_at)
	VALUES ($1,$2,$3,$4,$5,$6) RETURNING material_trade_id`
			errSql := db.Get(&tradeId, sqlStatement1, material.TreadId, materialId,
				material.CreatedBy, material.CreatedAt, material.LastUpdatedBy, material.LastUpdatedAt)
			if errSql != nil {
				fmt.Println("Error in adding rfq supplier")
				fmt.Println(errSql)
				return materials, errSql
			}
			appLog.Info("Trade Material Mapping ID:" + strconv.Itoa(tradeId))
		}
	}
	defer db.Close()
	return materials, err
}

func (self *materialDao) GetAllMaterials() ([]model.Material, error) {
	appLog.Info("### dao.GetAllMaterials() called. ###")

	materials := []model.Material{}
	// db, errs := db.DBConnect()
	db, errs := db.SqlxConnect()
	if errs != nil {
		fmt.Println(errs)
		return nil, errs
	}
	sqlStatement := `SELECT
					*	
				     FROM
				     materials
				`
	err := db.Select(&materials, sqlStatement)

	if err != nil {
		fmt.Printf("The error is: ", err)
		return nil, err
	}
	defer db.Close()
	return materials, err
}

func (self *materialDao) GetMaterialsOnGroup(groupCode string) ([]model.Material, error) {
	appLog.Info("GetMaterialsOnGroup is Called")

	materials := []model.Material{}
	// db, errs := db.DBConnect()
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return nil, connectionErrors
	}
	sqlStatement := `SELECT
					*	
				     FROM
				     materials where pid=$1
				`
	selectQueryError := db.Select(&materials, sqlStatement, groupCode)

	if selectQueryError != nil {
		fmt.Printf("The error is: ", selectQueryError)
		return nil, selectQueryError
	}
	defer db.Close()
	return materials, selectQueryError
}

func (self *materialDao) GetMaterialGroups() ([]model.Material, error) {
	appLog.Info("GetMaterialGroups is Called")

	materials := []model.Material{}
	// db, errs := db.DBConnect()
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return nil, connectionErrors
	}
	sqlStatement := `select DISTINCT pid,material_group from materials order by pid`
	selectQueryError := db.Select(&materials, sqlStatement)

	if selectQueryError != nil {
		fmt.Printf("The error is: ", selectQueryError)
		return nil, selectQueryError
	}
	defer db.Close()
	return materials, selectQueryError
}

func (self *materialDao) GetNestedMaterials(pids *model.Pids) ([]model.MaterialsObj, error) {
	appLog.Info("GetNestedMaterials is Called")
	var materialObjList = []model.MaterialsObj{}

	// db, errs := db.DBConnect()
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return materialObjList, connectionErrors
	}
	for _, groupName := range pids.Pid {
		var materialObj = model.MaterialsObj{}
		materials := []*model.Material{}
		materialObj.MaterialGroupName = groupName
		var pidList []string
		sqlStatement1 := `select pid from materials where material_group=$1`
		selectQueryError1 := db.Select(&pidList, sqlStatement1, groupName)
		materialObj.Pid = pidList[0]
		if selectQueryError1 != nil {
			fmt.Printf("The error is: ", selectQueryError1)
			return materialObjList, selectQueryError1
		}
		sqlStatement2 := `select * from materials where material_group=$1`
		selectQueryError2 := db.Select(&materials, sqlStatement2, groupName)
		if selectQueryError2 != nil {
			fmt.Printf("The error is: ", selectQueryError2)
			return materialObjList, selectQueryError2
		}
		materialObj.Child = materials
		materialObjList = append(materialObjList, materialObj)
	}
	defer db.Close()
	return materialObjList, connectionErrors
}

func (self *materialDao) GetAllMaterialWithTrades() ([]*model.MaterialBom, error) {
	materials := []*model.MaterialBom{}
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return materials, connectionErrors
	}

	sqlStatement := `select m.pid,m.material_code,m.material_group,m.discription,m.material_name,m.material_unit,m.base_price,m.gst,
		t.trade_name 
		from materials m join material_trade mt on 
		m.id=mt.material_id
		join trade t on
		mt.trade_id=t.trade_id;`
	selectQueryError := db.Select(&materials, sqlStatement)
	if selectQueryError != nil {
		fmt.Println(selectQueryError.Error())
	}

	defer db.Close()
	return materials, connectionErrors
}

func (self *materialDao) GetMaterialWithTrades(tradeList model.TradeList) ([]*model.MaterialBom, error) {
	materials := []*model.MaterialBom{}
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return materials, connectionErrors
	}

	sqlStatement := `select m.id,m.pid,m.material_code,m.material_group,
	    m.material_name,m.material_unit,m.sort_seq,
		t.trade_name 
		from materials m join material_trade mt on 
		m.id=mt.material_id
		join trade t on
		mt.trade_id=t.trade_id  where t.trade_name in(?)`
	query, args, err := sqlx.In(sqlStatement, tradeList.TradeNames)
	// fmt.Println(query)
	if err != nil {
		fmt.Printf("The error is", err)
		return materials, err
	}

	query = db.Rebind(query)
	queryErr := db.Select(&materials, query, args...)
	if queryErr != nil {
		fmt.Printf(queryErr.Error())
		return materials, queryErr
	}

	// fmt.Printf("projects are", address)

	defer db.Close()
	return materials, connectionErrors
}

func (self *materialDao) GetTopMaterialWithTrades(tradeList model.TradeList) ([]*model.MaterialBom, error) {
	materials := []*model.MaterialBom{}
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return materials, connectionErrors
	}

	sqlStatement := `select m.id,m.pid,m.material_code,m.material_group,m.material_name,
	m.material_unit,m.sort_seq,
		t.trade_name 
		from materials m join material_trade mt on 
		m.id=mt.material_id
		join trade t on
		mt.trade_id=t.trade_id  where t.trade_name in(?)
		and m.sort_seq <=25 order by m.sort_seq`
	query, args, err := sqlx.In(sqlStatement, tradeList.TradeNames)
	fmt.Println(query)

	if err != nil {
		fmt.Printf("The error is", err)
		return materials, err
	}

	query = db.Rebind(query)
	queryErr := db.Select(&materials, query, args...)
	if queryErr != nil {
		fmt.Printf(queryErr.Error())
		return materials, queryErr
	}

	// fmt.Printf("projects are", address)

	defer db.Close()
	return materials, connectionErrors
}

func (self *materialDao) GetTradesCategories(tradeName string) ([]model.TradeCategoriesList, error) {
	categories := []model.TradeCategoriesList{}
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return categories, connectionErrors
	}

	sqlStatement := `select DISTINCT m.material_group,pid
	from materials m join material_trade mt on 
	m.id=mt.material_id
	join trade t on
	mt.trade_id=t.trade_id  where t.trade_name = $1`
	query, args, err := sqlx.In(sqlStatement, tradeName)

	if err != nil {
		fmt.Println("Error:", err)
		return categories, err
	}
	query = db.Rebind(query)
	queryErr := db.Select(&categories, query, args...)
	if queryErr != nil {
		fmt.Printf(queryErr.Error())
		return categories, queryErr
	}
	defer db.Close()
	return categories, connectionErrors
}

func (self *materialDao) SearchMaterialExist(searchObj model.SearchMaterial) (int, error) {
	appLog.Info("### dao.GetRfqDetails() called. ###")

	// BeforeIndentMaterialSave(indentMat)

	count := []int{}
	db, connectionErrs := db.SqlxConnect()
	if connectionErrs != nil {

		fmt.Println(connectionErrs)
		return -1, connectionErrs
	}

	var queryErr error
	if (searchObj.TradeName == "" && len(searchObj.TradeName) == 0) && ((searchObj.CategoryName == "") && (len(searchObj.CategoryName) == 0)) {
		sqlStatement := `select count(*) from materials where material_name = $1`

		queryErr = db.Select(&count, sqlStatement, searchObj.MaterialName)
	} else if (searchObj.TradeName != "" && len(searchObj.TradeName) > 0) && ((searchObj.CategoryName == "") && (len(searchObj.CategoryName) == 0)) {
		sqlStatement := `select count(*)
		from materials m join material_trade mt on 
		m.id=mt.material_id
		join trade t on
		mt.trade_id=t.trade_id  where t.trade_name = $1 
		and m.material_name = $2`
		queryErr = db.Select(&count, sqlStatement, searchObj.TradeName, searchObj.MaterialName)
	} else {
		sqlStatement := `select count(*)
		from materials m join material_trade mt on 
		m.id=mt.material_id
		join trade t on
		mt.trade_id=t.trade_id  where t.trade_name = $1 
		and m.material_name = $2 and m.material_group=$3`
		queryErr = db.Select(&count, sqlStatement, searchObj.TradeName, searchObj.MaterialName, searchObj.CategoryName)
	}

	if queryErr != nil {
		fmt.Println(queryErr)
		return 1, queryErr
	}
	defer db.Close()
	if len(count) > 0 {
		return count[0], queryErr
	}
	return 1, queryErr
}
func ReplaceSQL(stmt, pattern string, len int) string {
	pattern += ","
	stmt = fmt.Sprintf(stmt, strings.Repeat(pattern, len))
	n := 0
	for strings.IndexByte(stmt, '?') != -1 {
		n++
		param := "$" + strconv.Itoa(n)
		stmt = strings.Replace(stmt, "?", param, 1)
	}
	return strings.TrimSuffix(stmt, ",")
}

func (self *materialDao) GetAllCategories() ([]model.TradeCategoriesList, error) {
	categories := []model.TradeCategoriesList{}
	db, connectionErrors := db.SqlxConnect()
	if connectionErrors != nil {
		fmt.Println(connectionErrors)
		return categories, connectionErrors
	}

	sqlStatement := `select DISTINCT m.material_group,pid
	from materials m join material_trade mt on 
	m.id=mt.material_id
	join trade t on
	mt.trade_id=t.trade_id`
	query, args, err := sqlx.In(sqlStatement)

	if err != nil {
		fmt.Println("Error:", err)
		return categories, err
	}
	query = db.Rebind(query)
	queryErr := db.Select(&categories, query, args...)
	if queryErr != nil {
		fmt.Printf(queryErr.Error())
		return categories, queryErr
	}
	defer db.Close()
	return categories, connectionErrors
}

func (self *materialDao) SearchAndGetMaterialExist(searchObj []string) ([]model.Material, error) {
	appLog.Info("### dao.GetRfqDetails() called. ###")

	// BeforeIndentMaterialSave(indentMat)

	materials := []model.Material{}
	db, connectionErrs := db.SqlxConnect()
	if connectionErrs != nil {

		fmt.Println(connectionErrs)
		return materials, connectionErrs
	}

	sqlStatement := `select * from materials where material_name in(?)`

	query, args, err := sqlx.In(sqlStatement, searchObj)
	fmt.Println(query)

	if err != nil {
		fmt.Println(err)
		return materials, err
	}

	query = db.Rebind(query)
	queryErr := db.Select(&materials, query, args...)
	if queryErr != nil {
		fmt.Printf(queryErr.Error())
		return materials, queryErr
	}
	defer db.Close()

	return materials, queryErr
}
