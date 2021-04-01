package gva

import (
	"github.com/gookit/color"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
)

var cfgFile string
var rootCmd = &cobra.Command{
	Use:   "gva",
	Short: "这是一款amazing的终端工具",
	Long: `欢迎使用gva终端工具
 ________ ____   ____   _____   
 /  _____/ \   \ /   /  /  _  \  
/   \  ___  \   Y   /  /  /_\  \ 
\    \_\  \  \     /  /    |    \
 \______  /   \___/   \____|__  /
        \/                    \/ 
`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		color.Warn.Println(err)
		os.Exit(1)
	}

}

func init() {

	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gva.yaml)")
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			color.Warn.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".gva" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".gva")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		color.Warn.Println("Using config file:", viper.ConfigFileUsed())
	}
}
