/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	rootPath string
)

type FileInfo struct {
	Folder, Name string
	Size         int64
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "direct-analyzer",
	Short: "Утилита для подсчёта статистики в директории (количество файлов, размер, строки).",
	Long:  "Утилита для подсчёта статистики в директории (количество файлов, размер, строки).",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		infos := []FileInfo{}
		rootPathLen := len(rootPath)
		filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			if !d.IsDir() {
				fileName := info.Name()
				folder := path[rootPathLen:]
				folder = folder[:len(folder)-len(fileName)]
				infos = append(infos, FileInfo{folder, fileName, info.Size()})
			}

			return nil
		})
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		for _, info := range infos {
			fmt.Fprintf(writer, "%s\t%s\t%d\n", info.Folder, info.Name, info.Size)
		}
		writer.Flush()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.direct-analyzer.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().StringVar(&rootPath, "path", ".", "Direction to analyze")
}
