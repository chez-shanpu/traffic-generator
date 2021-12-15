/*
Copyright © 2021 Tomoki Sugiura <cheztomo513@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package cmd

import (
	"github.com/chez-shanpu/traffic-generator/pkg/iperf3"
	"github.com/chez-shanpu/traffic-generator/pkg/option"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run traffic generator and out put its results",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := option.Config{}
		cfg.Populate()

		c, err := iperf3.NewIperfClientFromParamsFile(cfg)
		if err != nil {
			return err
		}

		rs, err := c.GenerateTraffic()
		if err != nil {
			return err
		}

		return c.OutputResults(rs, cfg.Out)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()
	flags.String(option.Param, "", "path to the param file")
	flags.StringP(option.DstAddr, "a", "", "destination ip address")
	flags.StringP(option.DstPort, "p", "", "destination port number")
	flags.Int64P(option.Mss, "m", 0, "TCP/SCTP maximum segment size")
	flags.Bool(option.UDP, false, "Run iperf3 client with udp option")
	flags.Bool(option.IPv6, false, "only ipv6")
	flags.Int64(option.Flowlabel, -1, "ipv6 flow label")
	flags.StringP(option.WindowSize, "w", "", "window size / socket buffer size")

	_ = viper.BindPFlags(flags)

	_ = runCmd.MarkFlagRequired(option.Param)
	_ = runCmd.MarkFlagRequired(option.DstAddr)
}
