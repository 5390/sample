package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"genMaterials/log"
	"genMaterials/model"
	"genMaterials/services"
	"net/http"

	"github.com/labstack/echo"
)

var appLog *log.AppLogger

func AddMaterial(c echo.Context) error {
	// Read file
	// Source
	file, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()
	// Destination
	dst, err := os.Create(file.Filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	// Copy
	if _, err = io.Copy(dst, src); err != nil {
		return err
	}
	// fmt.Printf("file Name:" + file.Filename)
	appLog.Info("File is recived sucessfully" + file.Filename)
	fileWithAbsolutePath, fileErr := filepath.Abs(file.Filename)
	appLog.Info("file is Strored at:" + fileWithAbsolutePath)
	if fileErr != nil {
		return fileErr
	}
	services.Uploadcsvfile(fileWithAbsolutePath)
	// services.DeleteFile(fileWithAbsolutePath)
	return c.String(http.StatusOK, "response")
}

func GetMaterial(c echo.Context) error {
	materialsList, ServiceError := services.GetAllMaterials()
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}

func GetAllMaterialWithTrades(c echo.Context) error {
	materialsList, ServiceError := services.GetAllMaterialWithTrades()
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}
func GetMaterialWithTrades(c echo.Context) error {
	tradeList := new(model.TradeList)
	bindErr := c.Bind(tradeList)
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, bindErr.Error())
	}
	materialsList, ServiceError := services.GetMaterialWithTrades(*tradeList)
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}

func GetTopMaterialWithTrades(c echo.Context) error {
	tradeList := new(model.TradeList)
	bindErr := c.Bind(tradeList)
	if bindErr != nil {
		return c.JSON(http.StatusInternalServerError, bindErr.Error())
	}
	materialsList, ServiceError := services.GetTopMaterialWithTrades(*tradeList)
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}

func GetTradesCategories(c echo.Context) error {
	tradeName := c.Param("tradename")

	tradeCategoriesList, ServiceError := services.GetTradesCategories(tradeName)
	if ServiceError != nil {
		return c.String(http.StatusInternalServerError, tradeCategoriesList)
	}
	return c.String(http.StatusOK, tradeCategoriesList)
}
func SearchMaterialExist(c echo.Context) error {
	searchMaterial := new(model.SearchMaterial)
	bindErr := c.Bind(searchMaterial)
	if bindErr != nil {
		fmt.Println(bindErr)
		return c.String(http.StatusInternalServerError, bindErr.Error())
	}

	tradeCategoriesList, ServiceError := services.SearchMaterialExist(*searchMaterial)
	if ServiceError != nil {
		return c.String(http.StatusInternalServerError, tradeCategoriesList)
	}
	return c.String(http.StatusOK, tradeCategoriesList)
}

//Function to check server health
func Health(c echo.Context) error {
	return c.String(http.StatusOK, "All ok")
}

func GetMaterialOnGroup(c echo.Context) error {
	groupCode := c.Param("groupCode")
	materialsList, ServiceError := services.GetMaterialsOnGroup(groupCode)
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}

func GetMaterialGroups(c echo.Context) error {
	// groupCode := c.Param("groupCode")
	materialsList, ServiceError := services.GetMaterialGroups()
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}

func GetNestedMaterials(c echo.Context) error {
	groupNames := new(model.Pids)
	err := c.Bind(groupNames)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err.Error())
	}
	// groupCode := c.Param("groupCode")
	materialsList, ServiceError := services.GetNestedMaterials(groupNames)
	if ServiceError != nil {
		return c.JSON(http.StatusInternalServerError, ServiceError.Error())
	}

	return c.JSON(http.StatusOK, materialsList)
}

func GetAllCategories(c echo.Context) error {

	tradeCategoriesList, ServiceError := services.GetAllCategories()
	if ServiceError != nil {
		return c.String(http.StatusInternalServerError, tradeCategoriesList)
	}
	return c.String(http.StatusOK, tradeCategoriesList)
}

//search and get materials if exist
func GetMaterialIfExist(c echo.Context) error {
	searchMaterialNameList := new([]string)
	bindErr := c.Bind(searchMaterialNameList)
	if bindErr != nil {
		fmt.Println(bindErr)
		return c.String(http.StatusInternalServerError, bindErr.Error())
	}

	tradeCategoriesList, ServiceError := services.SearchAndGetMaterialExist(*searchMaterialNameList)
	if ServiceError != nil {
		return c.String(http.StatusInternalServerError, tradeCategoriesList)
	}
	return c.String(http.StatusOK, tradeCategoriesList)
}
