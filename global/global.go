package global

import (
	"gin_class/config"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"github.com/spf13/viper"
)

var (
	GVA_DB     *gorm.DB
	GVA_CONFIG config.Server
	GVA_LOG    *zap.Logger
	GVA_VP     *viper.Viper
)
