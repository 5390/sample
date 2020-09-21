package services

import (
	"encoding/csv"
	"fmt"
	"genMaterials/common"
	"genMaterials/dao"
	"genMaterials/model"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

func Uploadcsvfile(fileNameWithPath string) {

	// b, err := ioutil.ReadFile("/home/rohit/GenMaterials_Documents/db.csv") // just pass the file name
	b, err := ioutil.ReadFile(fileNameWithPath) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	r := csv.NewReader(strings.NewReader(string(b)))

	counter := 0
	// materialList := []*model.Material{}
	materialList := []*model.Material{}
	for {
		// material := models.Material
		material := new(model.Material)
		counter++
		record, err := r.Read()
		if counter == 1 {
			continue
		}

		// Processing
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		// material.
		// materialsData = append(materialsData, record)
		// fmt.Println(record)
		material.Pid = record[0]
		material.MaterialCode = record[1]
		material.MaterialGroupName = record[2]
		material.MaterialName = record[3]
		material.Alias = record[4]
		material.Discription = record[5]
		material.MaterialUnit = record[6]
		material.Gst, _ = strconv.Atoi(record[7])
		material.BasePrice, _ = strconv.ParseFloat(record[8], 64)
		material.TreadId, _ = strconv.Atoi(record[9])
		SortSeq, seqError := strconv.Atoi(record[10])
		if seqError == nil {
			material.SortSeq = SortSeq
		}
		material.CreatedAt = time.Now()
		material.LastUpdatedAt = time.Now()
		material.CreatedBy = "Admin 1"
		material.LastUpdatedBy = "Admin 1"

		materialList = append(materialList, material)
		// fmt.Println(material)
	}
	_, _ = dao.MaterialDao().AddMaterial(materialList)
	//fmt.Println(materialList)
}

func GetAllMaterials() ([]model.Material, error) {

	getAllMaterials, getMaterialErr := dao.MaterialDao().GetAllMaterials()
	//fmt.Println(materialList)\
	return getAllMaterials, getMaterialErr
}

func GetMaterialsOnGroup(groupCode string) ([]model.Material, error) {

	getAllMaterials, getMaterialErr := dao.MaterialDao().GetMaterialsOnGroup(groupCode)
	//fmt.Println(materialList)\
	return getAllMaterials, getMaterialErr
}

func GetMaterialGroups() ([]model.Material, error) {

	getAllMaterials, getMaterialErr := dao.MaterialDao().GetMaterialGroups()
	//fmt.Println(materialList)\
	return getAllMaterials, getMaterialErr
}

func GetNestedMaterials(groupName *model.Pids) ([]model.MaterialsObj, error) {

	getAllMaterials, getMaterialErr := dao.MaterialDao().GetNestedMaterials(groupName)
	// fmt.Println(getAllMaterials)
	return getAllMaterials, getMaterialErr
}

func GetAllMaterialWithTrades() ([]model.MaterialsTrades, error) {
	materialMap := make(map[string][]string)
	getAllMaterialsWithTradeList, getMaterialErr := dao.MaterialDao().GetAllMaterialWithTrades()
	materialTrades := make(map[string][]string)

	materialByMap := make(map[string]model.MaterialBom)
	materialTradeList := []model.MaterialsTrades{}
	for _, material := range getAllMaterialsWithTradeList {
		materialMap[material.MaterialGroupName] = append(materialMap[material.MaterialGroupName], material.MaterialCode)
		materialTrades[material.MaterialCode] = append(materialTrades[material.MaterialCode], material.TradeName)
		if _, ok := materialByMap[material.MaterialCode]; !ok {
			materialByMap[material.MaterialCode] = *material
		}
	}
	// var exists = struct{}{}
	for materialGroupName, materialsByGrp := range materialMap {
		materialsTradesObj := new(model.MaterialsTrades)
		materialsTradesObj.GroupName = materialGroupName
		for _, materialCode := range materialsByGrp {
			if material, ok := materialByMap[materialCode]; ok {
				materialsTradesObj.TradeName = material.TradeName
				trades := materialTrades[materialCode]
				material.TradeList = trades
				material.TradeList = unique(trades)
				materialsTradesObj.MaterialList = append(materialsTradesObj.MaterialList, material)
			}
			delete(materialByMap, materialCode)
		}
		materialTradeList = append(materialTradeList, *materialsTradesObj)
	}
	return materialTradeList, getMaterialErr
}

// func GetMaterialWithTrades(tradeList model.TradeList) ([]model.MaterialsTrades, error) {
// 	materialMap := make(map[string][]model.Material)
// 	getAllMaterialsWithTradeList, getMaterialErr := dao.MaterialDao().GetMaterialWithTrades(tradeList)
// 	materialTrades := make(map[string][]string)
// 	materialTradeList := []model.MaterialsTrades{}
// 	for _, material := range getAllMaterialsWithTradeList {
// 		material.TradeList = append(material.TradeList, material.TradeName)

// 		materialMap[material.MaterialGroupName] = append(materialMap[material.MaterialGroupName], *material)
// 		materialTrades[material.MaterialCode] = append(materialTrades[material.MaterialCode], material.TradeName)
// 	}
// 	// var exists = struct{}{}
// 	for materialGroupName, materialsByGrp := range materialMap {
// 		materialsTradesObj := new(model.MaterialsTrades)
// 		materialsTradesObj.GroupName = materialGroupName
// 		materialsTradesObj.MaterialList = materialsByGrp
// 		for _, material := range materialsTradesObj.MaterialList {
// 			for _, materialTrade := range materialTrades[material.MaterialCode] {
// 				material.TradeList = append(material.TradeList, materialTrade)
// 			}
// 		}
// 		materialTradeList = append(materialTradeList, *materialsTradesObj)
// 	}
// 	return materialTradeList, getMaterialErr
// }

func GetMaterialWithTrades(tradeList model.TradeList) ([]model.MaterialsTrades, error) {
	materialMap := make(map[string][]string)
	getAllMaterialsWithTradeList, getMaterialErr := dao.MaterialDao().GetMaterialWithTrades(tradeList)
	for _, mat := range getAllMaterialsWithTradeList {
		if mat.TradeName == "" || len(mat.TradeName) == 0 {

			fmt.Printf("%v", mat)
		}
	}
	materialTrades := make(map[string][]string)
	materialByMap := make(map[string]model.MaterialBom)
	materialTradeList := []model.MaterialsTrades{}
	for _, material := range getAllMaterialsWithTradeList {
		materialMap[material.MaterialGroupName] = append(materialMap[material.MaterialGroupName], material.MaterialCode)
		materialTrades[material.MaterialCode] = append(materialTrades[material.MaterialCode], material.TradeName)
		if _, ok := materialByMap[material.MaterialCode]; !ok {
			materialByMap[material.MaterialCode] = *material
		}
	}
	// var exists = struct{}{}
	for materialGroupName, materialsByGrp := range materialMap {
		materialsTradesObj := new(model.MaterialsTrades)
		materialsTradesObj.GroupName = materialGroupName
		for _, materialCode := range materialsByGrp {
			if material, ok := materialByMap[materialCode]; ok {
				materialsTradesObj.TradeName = material.TradeName
				trades := materialTrades[materialCode]
				// sort.Strings(trades)
				material.TradeList = trades
				material.TradeList = unique(trades)
				materialsTradesObj.MaterialList = append(materialsTradesObj.MaterialList, material)
				delete(materialByMap, materialCode)
			}
		}
		if (materialsTradesObj.TradeName != "") && len(materialsTradesObj.MaterialList) > 0 {
			materialTradeList = append(materialTradeList, *materialsTradesObj)
		}
	}
	sort.Sort(MaterialSorter(materialTradeList))
	fmt.Println(len(materialTradeList))
	return materialTradeList, getMaterialErr
}

func GetTopMaterialWithTrades(tradeList model.TradeList) ([]model.MaterialsTrades, error) {
	materialMap := make(map[string][]string)
	getAllMaterialsWithTradeList, getMaterialErr := dao.MaterialDao().GetTopMaterialWithTrades(tradeList)
	materialTrades := make(map[string][]string)
	materialByMap := make(map[string]model.MaterialBom)
	materialTradeList := []model.MaterialsTrades{}
	for _, material := range getAllMaterialsWithTradeList {
		materialMap[material.MaterialGroupName] = append(materialMap[material.MaterialGroupName], material.MaterialCode)
		materialTrades[material.MaterialCode] = append(materialTrades[material.MaterialCode], material.TradeName)
		if _, ok := materialByMap[material.MaterialCode]; !ok {
			materialByMap[material.MaterialCode] = *material
		}
	}
	// var exists = struct{}{}
	for materialGroupName, materialsByGrp := range materialMap {
		materialsTradesObj := new(model.MaterialsTrades)
		materialsTradesObj.GroupName = materialGroupName
		for _, materialCode := range materialsByGrp {
			if material, ok := materialByMap[materialCode]; ok {
				materialsTradesObj.TradeName = material.TradeName
				trades := materialTrades[materialCode]
				// sort.Strings(trades)
				material.TradeList = trades
				material.TradeList = unique(trades)
				materialsTradesObj.MaterialList = append(materialsTradesObj.MaterialList, material)
				delete(materialByMap, materialCode)
			}
		}
		materialTradeList = append(materialTradeList, *materialsTradesObj)
	}
	sort.Sort(MaterialSorter(materialTradeList))
	return materialTradeList, getMaterialErr
}

func GetTradesCategories(tradeName string) (string, error) {
	tradeCategoriesList, tradeCategoriesErr := dao.MaterialDao().GetTradesCategories(tradeName)
	_, response := common.FormatResult(tradeCategoriesList, true, tradeCategoriesErr)
	return common.MarshalJson(response), tradeCategoriesErr
}

func GetAllCategories() (string, error) {
	tradeCategoriesList, tradeCategoriesErr := dao.MaterialDao().GetAllCategories()
	_, response := common.FormatResult(tradeCategoriesList, true, tradeCategoriesErr)
	return common.MarshalJson(response), tradeCategoriesErr
}

func SearchMaterialExist(searchObj model.SearchMaterial) (string, error) {
	tradeCategoriesList, tradeCategoriesErr := dao.MaterialDao().SearchMaterialExist(searchObj)
	if tradeCategoriesErr != nil {
		_, response := common.FormatResult(nil, false, tradeCategoriesErr)
		return common.MarshalJson(response), tradeCategoriesErr
	}
	if tradeCategoriesList == 0 {
		_, response := common.FormatResult(true, true, tradeCategoriesErr)
		return common.MarshalJson(response), tradeCategoriesErr
	}
	_, response := common.FormatResult(false, true, tradeCategoriesErr)
	return common.MarshalJson(response), tradeCategoriesErr
}

func unique(intSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func DeleteFile(path string) {
	// delete file
	var err = os.Remove(path)
	if err != nil {
		return
	}

	fmt.Println("==> done deleting file")
}

type MaterialSorter []model.MaterialsTrades

func (a MaterialSorter) Len() int           { return len(a) }
func (a MaterialSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a MaterialSorter) Less(i, j int) bool { return a[i].GroupName < a[j].GroupName }

//search and get materials if name matched
func SearchAndGetMaterialExist(searchObj []string) (string, error) {
	tradeCategoriesList, tradeCategoriesErr := dao.MaterialDao().SearchAndGetMaterialExist(searchObj)
	if tradeCategoriesErr != nil {
		_, response := common.FormatResult(tradeCategoriesList, false, tradeCategoriesErr)
		return common.MarshalJson(response), tradeCategoriesErr
	}

	_, response := common.FormatResult(tradeCategoriesList, true, tradeCategoriesErr)
	return common.MarshalJson(response), tradeCategoriesErr
}
