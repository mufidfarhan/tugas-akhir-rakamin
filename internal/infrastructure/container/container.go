package container

import (
	"backend-evermos/internal/helper"
	"backend-evermos/internal/infrastructure/mysql"
	"backend-evermos/internal/infrastructure/restclient"
	"backend-evermos/internal/pkg/repository"
	"backend-evermos/internal/pkg/usecase"
	"backend-evermos/internal/utils"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var v *viper.Viper

type (
	Container struct {
		Mysqldb       *gorm.DB
		Apps          *Apps
		AuthUsc       usecase.AuthUseCase
		UsersUsc      usecase.UsersUseCase
		ShopsUsc      usecase.ShopsUseCase
		ProductsUsc   usecase.ProductsUseCase
		CategoriesUsc usecase.CategoriesUseCase
		TrxUsc        usecase.TrxUseCase
		ProvcityUsc   usecase.ProvcityUseCase
	}

	Apps struct {
		Name      string `mapstructure:"name"`
		Host      string `mapstructure:"host"`
		Version   string `mapstructure:"version"`
		Address   string `mapstructure:"address"`
		HttpPort  int    `mapstructure:"httpport"`
		SecretJwt string `mapstructure:"secretJwt"`
	}
)

func loadEnv() {
	projectDirName := "go-example-cruid"
	projectName := regexp.MustCompile(`^(.*` + projectDirName + `)`)
	currentWorkDirectory, _ := os.Getwd()
	rootPath := projectName.Find([]byte(currentWorkDirectory))

	v.SetConfigFile(string(rootPath) + `/.env`)
}

func init() {
	v = viper.New()

	v.AutomaticEnv()
	loadEnv()

	path, err := os.Executable()
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, fmt.Sprintf("os.Executable panic : %s", err.Error()), err)
	}

	dir := filepath.Dir(path)
	v.AddConfigPath(dir)

	if err := v.ReadInConfig(); err != nil {
		helper.Logger(helper.LoggerLevelPanic, fmt.Sprintf("failed read config : %s", err.Error()), err)
	}

	err = v.ReadInConfig()
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, fmt.Sprintf("failed init config : %s", err.Error()), err)
	}

	helper.Logger(helper.LoggerLevelInfo, "Succeed read configuration file", err)
}

func AppsInit(v *viper.Viper) (apps Apps) {
	err := v.Unmarshal(&apps)
	if err != nil {
		helper.Logger(helper.LoggerLevelPanic, fmt.Sprint("Error when unmarshal configuration file : ", err.Error()), err)
	}
	helper.Logger(helper.LoggerLevelInfo, "Succeed when unmarshal configuration file", err)
	return
}

func InitContainer() (cont *Container) {
	apps := AppsInit(v)
	utils.InitJWT(apps.SecretJwt)
	mysqldb := mysql.DatabaseInit(v)
	restClient := restclient.New()

	userRepo := repository.NewUsersRepository(mysqldb)
	shopRepo := repository.NewShopsRepository(mysqldb)
	addressRepo := repository.NewAddressRepository(mysqldb)
	productRepo := repository.NewProductsRepository(mysqldb)
	productImageRepo := repository.NewProductImagesRepository(mysqldb)
	categoryRepo := repository.NewCategoriesRepository(mysqldb)
	trxRepo := repository.NewTrxRepository(mysqldb)
	trxDetailRepo := repository.NewTrxDetailsRepository(mysqldb)
	productLogRepo := repository.NewProductLogsRepository(mysqldb)
	provcityRepo := repository.NewProvcityRepository(restClient)

	authUsc := usecase.NewAuthUseCase(userRepo, shopRepo, provcityRepo)
	userUsc := usecase.NewUsersUseCase(userRepo, addressRepo, provcityRepo)
	shopUsc := usecase.NewShopsUseCase(shopRepo)
	productUsc := usecase.NewProductsUseCase(productRepo, shopRepo, productImageRepo, categoryRepo)
	categoryUsc := usecase.NewCategoriesUseCase(categoryRepo, shopRepo, productRepo)
	trxUsc := usecase.NewTrxUseCase(trxRepo, trxDetailRepo, productLogRepo, productRepo, addressRepo, productImageRepo)
	provCityUsc := usecase.NewProvcityUseCase(provcityRepo)

	return &Container{
		Apps:          &apps,
		Mysqldb:       mysqldb,
		AuthUsc:       authUsc,
		UsersUsc:      userUsc,
		ShopsUsc:      shopUsc,
		ProductsUsc:   productUsc,
		CategoriesUsc: categoryUsc,
		TrxUsc:        trxUsc,
		ProvcityUsc:   provCityUsc,
	}
}
