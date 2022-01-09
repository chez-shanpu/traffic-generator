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

	"github.com/chez-shanpu/traffic-generator/pkg/option"

	"github.com/chez-shanpu/traffic-generator/pkg/sts"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate traffic data and output",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := option.Config{}
		cfg.Populate()

		p := sts.NewPlanner(cfg)
		ps := p.GenerateTrafficParams()

		return ps.Output(cfg.Out)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	flags := initCmd.Flags()
	flags.Int(option.Cycle, 0, "number of traffic generation cycles")
	flags.Uint64(option.Seed, uint64(time.Now().UnixNano()), "seed for random values")
	flags.Float64(option.SendLambda, 0, "lambda of exponential distribution for send duration")
	flags.Int64(option.SendSeconds, 0, "send duration seconds")
	flags.Float64(option.WaitLambda, 0, "lambda of exponential distribution for wait duration")
	flags.Int64(option.WaitSeconds, 0, "wait duration seconds")
	flags.String(option.Bitrate, "", "traffic bitrate")
	flags.Float64(option.BitrateLambda, 0, "lambda of poisson distribution for bitrate")
	flags.String(option.BitrateUnit, "", "bitrate unit (e.g. K,M,G)")

	_ = viper.BindPFlags(flags)

	_ = rootCmd.MarkFlagRequired(option.Cycle)
	_ = rootCmd.MarkFlagRequired(option.Bitrate)
}
