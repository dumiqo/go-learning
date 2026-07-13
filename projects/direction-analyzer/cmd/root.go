/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

var (
	rootPath string
)

type FileNode struct {
	Path     string
	Name     string
	Size     uint64
	Children []*FileNode
	IsFolder bool
}

type FileInfo struct {
	Path, Name string
	Size       uint64
	IsFolder   bool
}

func (f *FileNode) Find(path string) *FileNode {
	if f == nil {
		return nil
	}
	if !f.IsFolder {
		return nil
	}
	fPath := strings.TrimSuffix(f.Path, "\\")
	path = strings.TrimSuffix(path, "\\")
	if fPath+"\\"+f.Name == path {
		return f
	}
	for _, v := range f.Children {
		x := v.Find(path)
		if x != nil {
			return x
		}
	}
	return nil
}

func (f *FileNode) TotalSize() uint64 {
	if f == nil {
		return 0
	}
	if !f.IsFolder {
		return f.Size
	}
	totalSize := f.Size
	for _, v := range f.Children {
		totalSize += v.TotalSize()
	}
	return totalSize
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "direct-analyzer",
	Short: "Утилита для подсчёта статистики в директории (количество файлов, размер, строки).",
	Long:  "Утилита для подсчёта статистики в директории (количество файлов, размер, строки).",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		infos := make([]FileInfo, 0)
		filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			info, err := d.Info()
			if err != nil {
				return err
			}
			if path == rootPath {
				return nil
			}
			fileName := info.Name()
			path = strings.TrimSuffix(path, "\\")
			folderPath := path[:len(path)-len(fileName)]
			infos = append(infos, FileInfo{folderPath, fileName, uint64(info.Size()), d.IsDir()})

			return nil
		})
		root := FileNode{rootPath, rootPath, 0, make([]*FileNode, 0), true}

		sort.Slice(infos, func(i, j int) bool {
			// Сначала папки, потом файлы
			if infos[i].IsFolder != infos[j].IsFolder {
				return infos[i].IsFolder // true идёт раньше
			}
			return strings.Count(infos[i].Path, "\\") < strings.Count(infos[j].Path, "\\")
		})
		for _, info := range infos {
			if info.Path == rootPath {
				root.Children = append(root.Children, &FileNode{info.Path, info.Name, info.Size, make([]*FileNode, 0), info.IsFolder})
			} else {
				node := root.Find(info.Path)

				if node == nil {
					break
				}

				node.Children = append(node.Children, &FileNode{info.Path, info.Name, info.Size, make([]*FileNode, 0), info.IsFolder})
			}
		}
		writer := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
		defer writer.Flush()
		fmt.Fprintf(writer, "Name\tSize\tIsFolder\n")
		for _, v := range root.Children {
			fmt.Fprintf(writer, "%s\t%d\t%t\n",
				v.Name,
				v.TotalSize(),
				v.IsFolder)
		}
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
