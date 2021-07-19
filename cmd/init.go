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
	"time"

	"github.com/chez-shanpu/traffic-generator/pkg/consts"
	"github.com/chez-shanpu/traffic-generator/pkg/sts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate traffic data and output",
	RunE: func(cmd *cobra.Command, args []string) error {
		cycle := viper.GetInt(consts.InitCmdCycleKey)
		seed := viper.GetUint64(consts.InitCmdSeedKey)
		lambda := viper.GetFloat64(consts.InitCmdLambdaKey)
		rate := viper.GetFloat64(consts.InitCmdRateKey)
		bitrate := viper.GetString(consts.InitCmdBitrateKey)
		p := sts.NewPlanner(cycle, seed, lambda, rate, bitrate)

		bits := p.CalcBitrates()
		sends := p.CalcSendSeconds()
		waits := p.CalcWaitSeconds()

		out := viper.GetString(consts.InitCmdOutKey)
		err := p.OutputParams(bits, sends, waits, out)

		return err
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.Int(consts.CycleFlag, 0, "number of traffic generation cycles")
	flags.Uint64(consts.SeedFlag, uint64(time.Now().UnixNano()), "seed for random values")
	flags.Float64(consts.LambdaFlag, 0, "lambda for poisson distribution")
	flags.Float64(consts.RateFlag, 0, "rate for exponential distribution")
	flags.String(consts.BitrateFlag, "0", "traffic bitrate")
	flags.StringP(consts.OutFlag, "o", "", "output file path")

	_ = viper.BindPFlag(consts.InitCmdCycleKey, flags.Lookup(consts.CycleFlag))
	_ = viper.BindPFlag(consts.InitCmdSeedKey, flags.Lookup(consts.SeedFlag))
	_ = viper.BindPFlag(consts.InitCmdLambdaKey, flags.Lookup(consts.LambdaFlag))
	_ = viper.BindPFlag(consts.InitCmdRateKey, flags.Lookup(consts.RateFlag))
	_ = viper.BindPFlag(consts.InitCmdBitrateKey, flags.Lookup(consts.BitrateFlag))
	_ = viper.BindPFlag(consts.InitCmdOutKey, flags.Lookup(consts.OutFlag))

	_ = rootCmd.MarkFlagRequired(consts.CycleFlag)
	_ = rootCmd.MarkFlagRequired(consts.LambdaFlag)
	_ = rootCmd.MarkFlagRequired(consts.RateFlag)
	_ = rootCmd.MarkFlagRequired(consts.BitrateFlag)
}
