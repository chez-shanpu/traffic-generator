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
	"os"
	"os/signal"
	"syscall"

	"github.com/chez-shanpu/traffic-generator/pkg/iperf3"
	"github.com/chez-shanpu/traffic-generator/pkg/option"
	"github.com/chez-shanpu/traffic-generator/pkg/traffic"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run server which handle only one client connection",
	Run: func(cmd *cobra.Command, args []string) {
		cfg := option.Config{}
		cfg.Populate()
		s := iperf3.NewServer()
		errCh := make(chan error, 1)

		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		defer func() {
			signal.Stop(sigs)
		}()

		var ress traffic.Results
		go func() {
			for {
				res, err := s.Run()
				if err != nil {
					errCh <- err
				}
				ress = append(ress, res)
			}
		}()

		select {
		case <-sigs:
			if err := s.OutputResult(ress, cfg.Out); err != nil {
				fmt.Printf("[ERROR]: %v", err)
			}
		case err := <-errCh:
			fmt.Printf("[ERROR]: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
