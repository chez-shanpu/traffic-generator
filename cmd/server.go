package cmd

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/chez-shanpu/traffic-generator/pkg/consts"
	"github.com/chez-shanpu/traffic-generator/pkg/iperf"
	"github.com/chez-shanpu/traffic-generator/pkg/tg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server which handle only one client connection",
	RunE: func(cmd *cobra.Command, args []string) error {
		res, err := iperf.RunServer()
		if err != nil {
			return err
		}

		out := viper.GetString(consts.ServerCmdOutKey)
		return outputServerResult(res, out)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	flags := serverCmd.Flags()
	flags.StringP(consts.OutFlag, "o", "", "path to the result file (if this value is empty the results will be output to stdout)")

	_ = viper.BindPFlag(consts.ServerCmdOutKey, flags.Lookup(consts.OutFlag))
}

func outputServerResult(res *tg.Result, out string) error {
	var f *os.File
	var err error

	if out == "" {
		f = os.Stdout
	} else {
		f, err = os.Create(out)
		if err != nil {
			return err
		}
	}

	err = outputServerResultsCSV(res, f)
	return err
}

func outputServerResultsCSV(res *tg.Result, f *os.File) error {
	w := csv.NewWriter(f)
	defer w.Flush()

	csvHead := []string{"TotalReceiveBytes", "SendSeconds"}
	err := w.Write(csvHead)
	if err != nil {
		return err
	}

	var line []string
	line = append(line, strconv.FormatInt(res.SendByte, 10))
	line = append(line, strconv.FormatFloat(res.SendSecond, 'f', -1, 64))
	err = w.Write(line)

	return err
}
