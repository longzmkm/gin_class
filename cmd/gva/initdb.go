package gva

import (
	"gin_class/core"
	"gin_class/global"
	"github.com/spf13/cobra"
)

var initdbCmd = &cobra.Command{
	Use:   "initdb",
	Short: "gin-vue-admin初始化数据",
	Long: `gin-vue-admin初始化数据适配数据库情况: 
1. mysql完美适配,
2. postgresql不能保证完美适配,
3. sqlite未适配,
4. sqlserver未适配`,
	Run: func(cmd *cobra.Command, args []string) {
		//frame, _ := cmd.Flags().GetString("frame")
		path, _ := cmd.Flags().GetString("path")
		global.GVA_VP = core.Viper(path)
		Mysql.CheckDatabase()
		Mysql.Init()
	},
}

func init() {
	rootCmd.AddCommand(initdbCmd)
	initdbCmd.Flags().StringP("path", "p", "./config.yaml", "自定配置文件路径(绝对路径)")
	initdbCmd.Flags().StringP("frame", "f", "gin", "可选参数为gin,gf")
	initdbCmd.Flags().StringP("type", "t", "mysql", "可选参数为mysql,postgresql,sqlite,sqlserver")
}
