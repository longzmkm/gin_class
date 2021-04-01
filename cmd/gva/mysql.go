package gva

import (
	"fmt"
	data "gin_class/cmd/information/system"
	"gin_class/global"
	"github.com/gookit/color"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strings"
)

var Mysql = &_mysql{_config: &gorm.Config{}}

type DatabaseInfo struct {
	Value        string `gorm:"column:Value"`
	VariableName string `gorm:"column:Variable_name"`
}

type _mysql struct {
	db      *gorm.DB
	err     error
	_config *gorm.Config

	old       string
	input     string
	version   string
	character string
}

// 初始化数据库
func (m *_mysql) Init() {
	if global.GVA_CONFIG.Mysql.LogMode {
		m._config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		m._config.Logger = logger.Default.LogMode(logger.Silent)
	}
	m._config.DisableForeignKeyConstraintWhenMigrating = true

	m.db, m.err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       global.GVA_CONFIG.Mysql.Dsn(), // DSN data source name
		DefaultStringSize:         191,                           // string 类型字段的默认长度
		DisableDatetimePrecision:  true,                          // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,                          // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,                          // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false,                         // 根据当前 MySQL 版本自动配置
	}), m._config)

	global.GVA_DB = m.db

}

// 检测数据库
func (m *_mysql) CheckDatabase() {
	var unknownDatabase = fmt.Sprintf("Unknown database '%v'", global.GVA_CONFIG.Mysql.Dbname)
	m.Init()
	if m.err != nil {
		if strings.Split(m.err.Error(), ": ")[1] == unknownDatabase {
			color.Debug.Print("\n[Mysql] -->配置文件的数据库名为:")
			color.LightGreen.Printf(" {%v} ", global.GVA_CONFIG.Mysql.Dbname)
			color.Debug.Println("不存在!\n")
			color.Debug.Println("您的配置文件所配置的数据库不存在,请选择:")
			color.Debug.Print("0:请自行创建配置文件所配置的数据库名为:")
			color.LightGreen.Printf(" {%v} \n", global.GVA_CONFIG.Mysql.Dbname)
			color.Debug.Print("1:尝试使用sql为您创建配置文件所配置的数据库名为:")
			color.LightGreen.Printf(" {%v} \n", global.GVA_CONFIG.Mysql.Dbname)
			color.Debug.Println("2:忽略错误! 注意: 如果不修复, 将会退出初始化数据的进程!")
			color.Warn.Println("\n注意!!!!!!!")
			color.Warn.Println("输入1之后,如果配置文件的mysql用户名为root才会有百分百的权限去创建数据库,不是root的话就会跳过创建数据库步骤!\n")
			color.Debug.Println("请输入指令:")
			if n, _ := fmt.Scanln(&m.input); n != 0 {
				if m.input == "1" {
					if global.GVA_CONFIG.Mysql.Username == "root" {
						m.database()
					} else {
						color.Debug.Print("\n很抱歉,您的配置文件的mysql用户名配置不是root,不确定你有无权限创建数据库,为您跳过创建数据库操作,请自行创建配置文件所配置的数据库名为:")
						color.LightGreen.Printf(" {%v} \n", global.GVA_CONFIG.Mysql.Dbname)
					}
				} else if m.input == "2" {
					os.Exit(0)
				} else {
					color.Warn.Println("[Mysql] --> 请自行创建数据库!")
					os.Exit(0)
				}
			}
		}
	}
}

// 初始化数据
func (m *_mysql) InitData() {
	if m.err = data.Admin.Init(); m.err != nil {
		color.Warn.Println("\n[Mysql] --> sys_users 表初始数据失败, err: %v", m.err)
		os.Exit(0)
	}
}

// 获取数据库版本
func (m *_mysql) Version() {
	color.Debug.Println("[Mysql] -->获取数据库版本中.......")
	if err := global.GVA_DB.Raw("SELECT VERSION() AS version;").Scan(&m.version).Error; err != nil {
		color.Info.Printf("[Mysql] -->获取数据库版本失败! err: %v", err)
		m.version = "未知版本~~~"
	}
	color.Debug.Printf("\n[Mysql] -->获取数据库版本成功!\n")
}

// 获取数据信息
func (m *_mysql) info() {

	m.Version()
	color.Debug.Print("\n您当前的数据库版本: ")
	color.LightGreen.Printf(" {%v} ", m.version)
	color.Debug.Print(", 使用的数据库是: ")
	color.LightGreen.Printf(" {%v} ", global.GVA_CONFIG.Mysql.Dbname)
	color.Debug.Print(", 数据库编码是: ")
	color.LightGreen.Printf(" {%v} \n\n", m.character)

}

//  创建配置文件的数据库
func (m *_mysql) database() {
	m.old = global.GVA_CONFIG.Mysql.Dbname
	global.GVA_CONFIG.Mysql.Dbname = "mysql"
	color.Debug.Printf("\n[Mysql] --> 正在连接 mysql 数据库中.......\n")

	m.Init()

	if m.err != nil {
		color.Error.Printf("\n[Mysql] --> 链接 mysql 数据库失败!, err: %v\n", m.err)
		color.Error.Printf("[Mysql] --> 请自行创建配置文件所需的数据库!\n")
		os.Exit(0)
	}
	color.Debug.Printf("\n[Mysql] --> 连接 mysql 数据库成功\n")
	global.GVA_CONFIG.Mysql.Dbname = m.old
	color.Debug.Print("\n[Mysql] --> 正在为您创建数据库名为:")
	color.LightGreen.Printf(" {%v} ", global.GVA_CONFIG.Mysql.Dbname)
	color.Debug.Print("中.......\n")
	if m.err = global.GVA_DB.Exec("CREATE DATABASE IF NOT EXISTS " + global.GVA_CONFIG.Mysql.Dbname + " DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_general_ci;").Error; m.err != nil {
		color.Debug.Print("\n[Mysql] --> 创建数据库名为:")
		color.LightGreen.Printf(" {%v} ", global.GVA_CONFIG.Mysql.Dbname)
		color.Debug.Print("失败!请手动修改数据库名为")
		color.LightGreen.Printf(" {%v} \n", global.GVA_CONFIG.Mysql.Dbname)
		os.Exit(0)
		return
	}
	color.Debug.Print("\n[Mysql] --> 正在为您创建数据库名为:")
	color.LightGreen.Printf(" {%v} ", global.GVA_CONFIG.Mysql.Dbname)
	color.Debug.Print("成功!\n")

}
