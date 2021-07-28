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
	"github.com/chez-shanpu/traffic-generator/pkg/consts"
	"github.com/chez-shanpu/traffic-generator/pkg/iperf"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run traffic generator and out put its results",
	RunE: func(cmd *cobra.Command, args []string) error {
		paramFile := viper.GetString(consts.RunCmdParamKey)
		dstAddr := viper.GetString(consts.RunCmdDstAddrKey)
		dstPort := viper.GetString(consts.RunCmdDstPortKey)
		mss := viper.GetInt64(consts.RunCmdMssKey)
		udp := viper.GetBool(consts.RunCmdUdpKey)
		ipv6 := viper.GetBool(consts.RunCmdIPv6Key)
		flowlabel := viper.GetInt64(consts.RunCmdFlowlabelKey)
		c, err := iperf.NewIperfClientFromParamsFile(dstAddr, dstPort, mss, udp, ipv6, flowlabel, paramFile)
		if err != nil {
			return err
		}

		rs, err := c.GenerateTraffic()
		if err != nil {
			return err
		}

		out := viper.GetString(consts.RunCmdOutKey)
		err = c.OutputResults(rs, out)

		return err
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	flags := runCmd.Flags()
	flags.String(consts.ParamFlag, "", "path to the param file")
	flags.StringP(consts.DstAddrFlag, "a", "", "destination ip address")
	flags.StringP(consts.DstPortFlag, "p", "", "destination port number")
	flags.StringP(consts.OutFlag, "o", "", "path to the result file (if this value is empty the results will be output to stdout)")
	flags.Int64P(consts.MssFlag, "m", 0, "TCP/SCTP maximum segment size")
	flags.Bool(consts.UdpFlag, false, "Run iperf3 client with udp option")
	flags.Bool(consts.IPv6Flag, false, "only ipv6")
	flags.Int64(consts.FlowlabelFlag, -1, "ipv6 flow label")

	_ = viper.BindPFlag(consts.RunCmdParamKey, flags.Lookup(consts.ParamFlag))
	_ = viper.BindPFlag(consts.RunCmdDstAddrKey, flags.Lookup(consts.DstAddrFlag))
	_ = viper.BindPFlag(consts.RunCmdDstPortKey, flags.Lookup(consts.DstPortFlag))
	_ = viper.BindPFlag(consts.RunCmdOutKey, flags.Lookup(consts.OutFlag))
	_ = viper.BindPFlag(consts.RunCmdMssKey, flags.Lookup(consts.MssFlag))
	_ = viper.BindPFlag(consts.RunCmdUdpKey, flags.Lookup(consts.UdpFlag))
	_ = viper.BindPFlag(consts.RunCmdIPv6Key, flags.Lookup(consts.IPv6Flag))
	_ = viper.BindPFlag(consts.RunCmdFlowlabelKey, flags.Lookup(consts.FlowlabelFlag))

	_ = rootCmd.MarkFlagRequired(consts.ParamFlag)
	_ = rootCmd.MarkFlagRequired(consts.DstAddrFlag)
}
