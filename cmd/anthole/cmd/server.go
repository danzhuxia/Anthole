package cmd

import (
	"Anthole/pkg/common"
	"Anthole/pkg/server"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string
var confChan = make(chan *common.AntHoleConfig)
var errChan = make(chan error, 1)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Anthole Server",
	Run: func(cmd *cobra.Command, args []string) {
		antAscii:=`
  ___        _   _           _      
 / _ \      | | | |         | |     
/ /_\ \_ __ | |_| |__   ___ | | ___ 
|  _  | '_ \| __| '_ \ / _ \| |/ _ \
| | | | | | | |_| | | | (_) | |  __/
\_| |_/_| |_|\__|_| |_|\___/|_|\___|
`
		fmt.Println(antAscii)
		fmt.Println("Server Called...")
		fmt.Println("Use server config filePath: ", cfgFile)
		antholeConf, err := common.GetConfig(cfgFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, "[init config error] ", err)
			os.Exit(1)
		}
		go antServerStart()
		confChan <- antholeConf
		for {
			if err, ok := <-errChan; ok {
				fmt.Fprintln(os.Stderr, "[start server error] ", err)

			} else {
				break
			}
		}
	},
}

func init() {
	// cfgFile = *serverCmd.Flags().StringP("config", "c", "", "config file (default is $HOME/.anthole.yaml)")
	serverCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.anthole.yaml)")
	rootCmd.AddCommand(serverCmd)
}

func antServerStart() {
	logrus.Info("Server is Starting...")
	if conf, ok := <-confChan; ok {
		close(confChan)
		for serviceid, service := range conf.Client.Services {
			go func(serveid, remoteport int) {
				err := server.CreateWorker(serveid, remoteport)
				if err != nil {
					errChan <- err
				}
			}(serviceid, service.RemotePort)
		}

		master := server.StartMasterInstance()

		go func(port int) {
			err := master.Run(port)
			if err != nil {
				errChan <- err
			}
		}(conf.Server.Port)
	}
}
