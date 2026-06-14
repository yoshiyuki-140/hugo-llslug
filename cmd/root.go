package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yoshiyuki-140/hugo-llslug/internal/adapter/cli"
	hugoAdapter "github.com/yoshiyuki-140/hugo-llslug/internal/adapter/hugo"
	"github.com/yoshiyuki-140/hugo-llslug/internal/adapter/ollama"
	"github.com/yoshiyuki-140/hugo-llslug/internal/usecase"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "hugo-llslug",
	Short: "ローカルLLMを使ってHugoのSlugを自動生成するCLIツール",
	Long: `llslug は、ローカルLLM（Ollama）を活用して、日本語の記事タイトルから
URLフレンドリーな英語のケバブケース（Slug）を自動生成する、Hugo専用のCLI拡張ツールです。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		modelName, _ := cmd.Flags().GetString("model")
		client := ollama.NewClient(modelName)
		executor := hugoAdapter.NewExecutor()
		uc := usecase.NewSlugUsecase(client, executor)
		runner := cli.NewRunner(uc)
		return runner.Run()
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.hugo-llslug.yaml)")
	rootCmd.Flags().String("model", "", "使用するOllamaモデル名 (default: qwen3.5:0.8b)")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".hugo-llslug")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
