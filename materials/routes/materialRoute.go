package route

import (
	"genMaterials/controllers"

	"github.com/labstack/echo"
)

func MaterialRouteService(e *echo.Echo) {

	e.POST("/mm/material/csv", controllers.AddMaterial)
	e.GET("/mm/health", controllers.Health)
	e.GET("/mm/material/listall", controllers.GetMaterial)
	e.GET("/mm/material/groupList/:groupCode", controllers.GetMaterialOnGroup)
	e.GET("/mm/material/groups", controllers.GetMaterialGroups)
	e.GET("/mm/material/materialsSpecs", controllers.GetNestedMaterials)
	e.GET("/mm/material/getall/trades", controllers.GetAllMaterialWithTrades)
	e.GET("/mm/material/get/trades", controllers.GetMaterialWithTrades)
	e.GET("/mm/topmaterial/get/trades", controllers.GetTopMaterialWithTrades)
	e.GET("/mm/trade/get/categories/:tradename", controllers.GetTradesCategories)
	e.GET("/mm/material/search/materialname", controllers.SearchMaterialExist)
	e.GET("/mm/trade/get/all/categories", controllers.GetAllCategories)
	e.GET("/mm/material/get/searched/material", controllers.GetMaterialIfExist)

}
