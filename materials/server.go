package main

import (
	"fmt"
	"genMaterials/common"
	"genMaterials/db"
	"genMaterials/log"
	route "genMaterials/routes"

	"github.com/casbin/casbin"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	_ "github.com/lib/pq"
)

const ymlPrefix = "config" //config.yml

var yml common.Yaml

var appLogger *log.AppLogger

// func init() {
// 	yml := common.InitYamlConfig()
// 	appLogger := log.GetApplogger()
// 	viper.SetConfigName(ymlPrefix) // actually 'fileName + yml'
// 	viper.AddConfigPath(".")

// 	if err := viper.ReadInConfig(); err != nil {
// 		appLogger.Fatal("Error reading config file : ")
// 	}

// 	err := viper.Unmarshal(&yml)
// 	if err != nil {
// 		appLogger.Fatal("unable to decode into struct : ")
// 	}

// }

type Enforcer struct {
	enforcer casbin.Enforcer
}

func main() {

	yml := common.InitYamlConfig()
	fmt.Println(yml)
	appLogger.Info("HTTP Server Initialize...")

	/* sqlx start */
	_, err := db.DBConnect()
	if err != nil {
		panic(err)
	}
	_, er := common.ReadPropertiesFile("sample_test.properties")

	if er != nil {
		panic(er)
	}

	e := echo.New()
	e.Use(middleware.CORS())
	cE := casbin.NewEnforcer("model.conf", "policy.csv")
	enforcer := Enforcer{enforcer: *cE}
	e.Use(enforcer.Enforce)
	e.Use(common.TraceLogInterceptor())

	// route.ProjectRouteService(e)
	route.MaterialRouteService(e)
	// route.MsUserActivate(e)

	e.Static("/", yml.Echo.Static)
	e.Logger.Fatal(e.Start(yml.Echo.Port))
}

func (e *Enforcer) Enforce(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Request().Header.Get("Authorization")
		method := c.Request().Method
		path := c.Request().URL.Path
		result := e.enforcer.Enforce(user, path, method)

		if result {
			return next(c)
		}
		return echo.ErrForbidden
	}
}
