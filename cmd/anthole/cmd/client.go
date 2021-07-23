/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

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
	"Anthole/pkg/client"
	"Anthole/pkg/common"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Anthole Client Option",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("client called")

		antholeConf, err := common.GetConfig(cfgFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "config error:", err)
			os.Exit(1)
		}
		err = client.AntClientStart(antholeConf)
		if err != nil {
			fmt.Fprintln(os.Stderr, "init client error:", err)
			os.Exit(1)
		}
	},
}

func init() {
	cfgFile = *clientCmd.Flags().StringP("config", "c", "", "config file (default is /.anthole.yaml)")
	rootCmd.AddCommand(clientCmd)
}
