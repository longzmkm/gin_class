package initialize

import (
	"fmt"
	"gin_class/global"
	"gin_class/initialize/internal"
	"gin_class/model"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func Gorm() *gorm.DB {
	fmt.Println(global.GVA_CONFIG.System.DbType)
	switch global.GVA_CONFIG.System.DbType {
	case "mysql":
		return GormMysql()
	default:
		return GormMysql()
	}

}

func MysqlTables(db *gorm.DB) {

	err := db.AutoMigrate(
		model.SysUser{},
	)

	if err != nil {
		fmt.Println("register table failed", zap.Any("err", err))
	}
}

func GormMysql() *gorm.DB {
	m := global.GVA_CONFIG.Mysql
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	// TODO 这个地方可以优化
	//dsn := m.Username + ":" + m.Password + "@tcp(" + m.Path + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		//DriverName                string
		DSN: dsn, // DSN data source name
		//Conn                      gorm.ConnPool
		SkipInitializeWithVersion: false, // 根据版本自动配置
		DefaultStringSize:         191,
		//DefaultDatetimePrecision  *int
		DisableDatetimePrecision: false, // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:   true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:  true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		//DontSupportForShareClause bool
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), gormConfig(m.LogMode)); err != nil {
		global.GVA_LOG.Error("MySQL启动异常", zap.Any("err", err))
		os.Exit(0)
		return nil
	} else {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxIdleConns(m.MaxIdleConns)
		sqlDB.SetMaxOpenConns(m.MaxOpenConns)
		return db
	}

}

//@author: SliverHorn
//@function: gormConfig
//@description: 根据配置决定是否开启日志
//@param: mod bool
//@return: *gorm.Config

func gormConfig(mod bool) *gorm.Config {
	var config = &gorm.Config{DisableForeignKeyConstraintWhenMigrating: true}

	switch global.GVA_CONFIG.Mysql.LogZap {

	case "silent", "Silent":
		config.Logger = internal.Default.LogMode(logger.Silent)
	case "error", "Error":
		config.Logger = internal.Default.LogMode(logger.Error)
	case "warn", "Warn":
		config.Logger = internal.Default.LogMode(logger.Warn)
	case "info", "Info":
		config.Logger = internal.Default.LogMode(logger.Info)
	case "zap", "Zap":
		config.Logger = internal.Default.LogMode(logger.Info)
	default:
		if mod {
			config.Logger = internal.Default.LogMode(logger.Info)
			break
		}
		config.Logger = internal.Default.LogMode(logger.Silent)
	}

	return config
}
