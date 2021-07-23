package cmd

import (
	"Anthole/pkg/common"
	"Anthole/pkg/server"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var confChan = make(chan *common.AntHoleConfig, 1)
var errChan chan error

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Anthole Server",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		antholeConf, err := common.GetConfig(cfgFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		confChan <- antholeConf
		err = antServerStart()
		if err != nil {
			cobra.CheckErr(err)
		}
	},
}

func init() {
	cfgFile = *serverCmd.Flags().StringP("config", "c", "", "config file (default is $HOME/.anthole.yaml)")
	// serverCmd.Flags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.anthole.yaml)")
	rootCmd.AddCommand(serverCmd)
}

func antServerStart() (err error) {

	for {
		select {
		case conf := <-confChan:
			close(confChan)
			// 启动worker
			for serviceid, service := range conf.Client.Services {
				go func(serveid, remoteport int) {
					err := server.CreateWorker(serveid, remoteport)
					if err != nil {
						errChan <- err
					}
				}(serviceid, service.RemotePort)
			}

			master := server.StartMasterInstance()
			err = master.Run(conf.Server.Port)
			if err != nil {
				errChan <- err
			}
		case err := <-errChan:
			return err
		}
	}
}
