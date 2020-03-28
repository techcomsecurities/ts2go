/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/techcomsecurities/ts2go/convert"
)

var opt convert.I2SOpt

// i2sCmd represents the i2s command
var i2sCmd = &cobra.Command{
	Use:   "i2s",
	Short: "Convert TS interface to Golang struct",
	Long:  `Convert TS interface to Golang struct.`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infof("start convert TS interface in %s to Golang struct", opt.FilePath)
		_, err := convert.NewI2SConverter(opt).Run()
		if err != nil {
			log.WithError(err).Fatalf("fail to convert i2s of %s", opt.FilePath)
		}
		log.Info("Done")
	},
}

func init() {
	rootCmd.AddCommand(i2sCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// i2sCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	i2sCmd.Flags().StringVarP(&opt.FilePath, "filePath", "f", "", "Path to .ts interface file")
	i2sCmd.Flags().StringVarP(&opt.UnionTypePriority, "utp", "u", "f", "Type priority in case of TS prop is union type, f is first, l is last")
	i2sCmd.Flags().StringVarP(&opt.OutputPath, "output", "o", "", "Path to .go output file")
}
