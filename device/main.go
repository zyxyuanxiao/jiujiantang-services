package main

import (
	"os"

	"github.com/jinmukeji/jiujiantang-services/device/config"
	handler "github.com/jinmukeji/jiujiantang-services/device/handler"
	"github.com/jinmukeji/jiujiantang-services/device/mysqldb"
	logger "github.com/jinmukeji/jiujiantang-services/pkg/rpc"
	dbutilmysql "github.com/jinmukeji/plat-pkg/v2/dbutil/mysql"
	storemysql "github.com/jinmukeji/plat-pkg/v2/store/mysql"
	proto "github.com/jinmukeji/proto/v3/gen/micro/idl/partner/xima/device/v1"
	"github.com/micro/cli/v2"
	micro "github.com/micro/go-micro/v2"
)

func main() {
	versionMeta := config.NewVersionMetadata()
	// Create a new service. Optionally include some options here.
	authenticationWrapper := new(handler.AuthenticationWrapper)
	service := micro.NewService(
		// Service Basic Info
		micro.Name(config.FullServiceName()),
		micro.Version(config.ProductVersion),

		// Fault Tolerance - Heartbeating
		// 	 See also: https://micro.mu/docs/fault-tolerance.html#heartbeating
		micro.RegisterTTL(config.DefaultRegisterTTL),
		micro.RegisterInterval(config.DefaultRegisterInterval),

		// Setup wrappers
		micro.WrapHandler(logger.LogWrapper, authenticationWrapper.HandleWrapper()),

		// // Setup runtime flags
		dbClientOptions(),

		// Setup --version flag
		defaultVersionFlags(),

		// Setup metadata
		micro.Metadata(versionMeta),
	)

	// optionally setup command line usage
	service.Init(
		micro.Action(func(c *cli.Context) error {
			if c.Bool("version") {
				config.PrintFullVersionInfo()
				os.Exit(0)
			}
			return nil
		}),
	)

	log.Infof("Starting service: %s", config.FullServiceName())
	log.Infof("Product Version: %s", config.ProductVersion)
	log.Infof("Git SHA: %s", config.GitSHA)
	log.Infof("Git Branch: %s", config.GitBranch)
	log.Infof("Go Version: %s", config.GoVersion)
	log.Infof("Build Version: %s", config.BuildVersion)
	log.Infof("Build Time: %s", config.BuildTime)

	// Register handler
	server := service.Server()
	db, err := newDbClient()
	if err != nil {
		log.Panicf("Failed to connect to MySQL instance at %s. Error: %v", dbAddress, err)
	}
	log.Infoln("Connected to MySQL instance at", dbAddress)
	authenticationWrapper.SetDataStore(db)
	deviceManagerService := handler.NewDeviceManagerService(db)
	if err := proto.RegisterDeviceManagerAPIHandler(server, deviceManagerService); err != nil {
		log.Fatalln(err)
	}
	// Run the server
	if err := service.Run(); err != nil {
		log.Fatalln(err)
	}
}

// 数据库和邮件服务器连接信息
var (
	dbAddress   string
	dbUsername  string
	dbPassword  string
	dbDatabase  string
	dbEnableLog = false
	dbMaxConns  = 1
)

func defaultVersionFlags() micro.Option {
	return micro.Flags(
		&cli.BoolFlag{
			Name:  "version",
			Usage: "Show version information",
		},
	)
}

// dbClientOptions 构建命令行启动参数
func dbClientOptions() micro.Option {
	return micro.Flags(
		&cli.StringFlag{
			Name:        "x_db_address",
			Value:       "localhost:3306",
			Usage:       "MySQL instance `ADDRESS` - [host]:[port]",
			EnvVars:     []string{"X_DB_ADDRESS"},
			Destination: &dbAddress,
		},
		&cli.StringFlag{
			Name:        "x_db_username",
			Usage:       "MySQL login `USERNAME`",
			EnvVars:     []string{"X_DB_USERNAME"},
			Destination: &dbUsername,
		},
		&cli.StringFlag{
			Name:        "x_db_password",
			Usage:       "MySQL login `PASSWORD`",
			EnvVars:     []string{"X_DB_PASSWORD"},
			Destination: &dbPassword,
		},
		&cli.StringFlag{
			Name:        "x_db_database",
			Usage:       "MySQL database name",
			EnvVars:     []string{"X_DB_DATABASE"},
			Destination: &dbDatabase,
		},
		&cli.BoolFlag{
			Name:        "x_db_enable_log",
			Usage:       "Enable MySQL client log",
			EnvVars:     []string{"X_DB_ENABLE_LOG"},
			Destination: &dbEnableLog,
		},
		&cli.IntFlag{
			Name:        "x_db_max_connections",
			Usage:       "Max connections of MySQL client",
			EnvVars:     []string{"X_DB_MAX_CONNECTIONS"},
			Value:       1,
			Destination: &dbMaxConns,
		},
	)
}

// newDbClient 创建一个 DbClient
func newDbClient() (*mysqldb.DbClient, error) {
	cfg := dbutilmysql.NewConfig()
	cfg.User = dbUsername
	cfg.Passwd = dbPassword
	cfg.Net = "tcp"
	cfg.Addr = dbAddress
	cfg.DBName = dbDatabase
	db, err := dbutilmysql.OpenGormDB(
		dbutilmysql.WithMySQLConfig(cfg),
	)
	if err != nil {
		return nil, err
	}

	return mysqldb.NewDbClient(*storemysql.NewStore(db))
}
