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
	"fmt"
	"github.com/chez-shanpu/traffic-generator/pkg/iperf"
	"github.com/chez-shanpu/traffic-generator/pkg/sts"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tg",
	Short: "tg is a traffic generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		dstAddr := viper.GetString("tg.dstAddr")
		dstPort := viper.GetString("tg.dstPort")
		c := iperf.NewIperfClient(dstAddr, dstPort)

		cycle := viper.GetInt("tg.cycle")
		seed := viper.GetUint64("tg.seed")
		lambda := viper.GetFloat64("tg.lambda")
		rate := viper.GetFloat64("tg.rate")
		p := sts.NewPlanner(cycle, seed, lambda, rate)

		c.BlockCounts = p.CalcPacketCounts()
		c.WaitDurations = p.CalcWaitDurations()

		res, err := c.GenerateTraffic()
		if err != nil {
			return err
		}

		out := viper.GetString("tg.out")
		err = c.OutputResults(res, out)

		return err
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	flags := rootCmd.Flags()

	flags.Int("cycle", 0, "number of traffic generation cycles")
	flags.StringP("out", "o", "", "output file path")
	flags.StringP("dst-addr", "d", "", "destination ip address")
	flags.StringP("dst-port", "p", "", "destination port number")
	flags.Float64("lambda", 0, "lambda for poisson distribution")
	flags.Float64("rate", 0, "rate for exponential distribution")
	flags.Uint64("seed", uint64(time.Now().UnixNano()), "seed for random values")

	_ = viper.BindPFlag("tg.cycle", flags.Lookup("cycle"))
	_ = viper.BindPFlag("tg.out", flags.Lookup("out"))
	_ = viper.BindPFlag("tg.dstAddr", flags.Lookup("dst-addr"))
	_ = viper.BindPFlag("tg.dstPort", flags.Lookup("dst-port"))
	_ = viper.BindPFlag("tg.lambda", flags.Lookup("labda"))
	_ = viper.BindPFlag("tg.rate", flags.Lookup("rate"))
	_ = viper.BindPFlag("tg.seed", flags.Lookup("seed"))

	_ = rootCmd.MarkFlagRequired("duration")
	_ = rootCmd.MarkFlagRequired("dst-addr")
	_ = rootCmd.MarkFlagRequired("lambda")
	_ = rootCmd.MarkFlagRequired("rate")

}

