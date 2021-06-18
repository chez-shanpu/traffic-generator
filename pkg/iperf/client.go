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

package iperf

import (
	"encoding/csv"
	"encoding/json"
	"github.com/chez-shanpu/traffic-generator/pkg/tg"
	"os"
	"os/exec"
	"strconv"
	"time"
)

const iperf = "iperf3"

type Client struct {
	DstAddr       string
	DstPort       string
	BlockCounts   []tg.PacketCount
	WaitDurations []tg.WaitDuration
}

func NewIperfClient(dstAddr, dstPort string) *Client {
	return &Client{
		DstAddr:       dstAddr,
		DstPort:       dstPort,
		BlockCounts:   nil,
		WaitDurations: nil,
	}
}

func (c Client) TotalBlockCounts() tg.PacketCount {
	var total tg.PacketCount
	for _, b := range c.BlockCounts {
		total += b
	}
	return total
}

func (c Client) TotalWaitDurations() tg.WaitDuration {
	var total tg.WaitDuration
	for _, w := range c.WaitDurations {
		total += w
	}
	return total
}

func (c Client) GenerateTraffic() (*tg.Result, error) {
	res := tg.Result{}
	for i := 0; i < len(c.WaitDurations); i++ {
		tb, ts, err := c.execIperf(c.BlockCounts[i])
		if err != nil {
			return nil, err
		}
		res.TransferBytes = append(res.TransferBytes, tb)
		res.TransferSeconds = append(res.TransferSeconds, ts)
		time.Sleep(time.Duration(c.WaitDurations[i]) * time.Millisecond)
	}

	res.TotalTransferBytes = calcTotalTransferBytes(res.TransferBytes)
	res.TotalTransferSeconds = calcTotalTransferSeconds(res.TransferSeconds)

	return &res, nil
}

func (c Client) execIperf(pc tg.PacketCount) (tb int64, ts float64, err error) {
	args := []string{
		"-c",
		c.DstAddr,
		"-f", "K",
		"-k", strconv.FormatInt(int64(pc), 10),
		"-J",
	}

	if c.DstPort != "" {
		args = append(args, "-p")
		args = append(args, c.DstPort)
	}

	out, err := exec.Command(iperf, args...).Output()
	if err != nil {
		return 0, 0, nil
	}

	tb, ts, err = parseIperfOutput(out)
	return tb, ts, err
}

func parseIperfOutput(out []byte) (tb int64, ts float64, err error) {
	var i interface{}
	err = json.Unmarshal(out, &i)
	if err != nil {
		return 0, 0, err
	}

	tb = int64(i.(map[string]interface{})["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["bytes"].(float64))
	ts = i.(map[string]interface{})["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["seconds"].(float64)

	return tb, ts, nil
}

func calcTotalTransferBytes(tbs []int64) int64 {
	total := int64(0)
	for _, tb := range tbs {
		total += tb
	}
	return total
}

func calcTotalTransferSeconds(tss []float64) float64 {
	total := 0.0
	for _, ts := range tss {
		total += ts
	}
	return total
}

func (c Client) OutputResults(res *tg.Result, out string) error {
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

	err = c.outputResultsCSV(res, f)
	return err
}

func (c Client) outputResultsCSV(res *tg.Result, f *os.File) error {
	w := csv.NewWriter(f)
	defer w.Flush()

	csvHead := []string{"Cycle", "Transfer", "Bitrate", "Packets", "WaitDuration"}
	err := w.Write(csvHead)
	if err != nil {
		return err
	}

	for i := 0; i < len(res.TransferBytes); i++ {
		var line []string
		line = append(line, strconv.Itoa(i))
		line = append(line, strconv.FormatInt(res.TransferBytes[i], 10))
		line = append(line, strconv.FormatFloat(res.TotalTransferSeconds, 'f', -1, 64))
		line = append(line, strconv.FormatInt(int64(c.BlockCounts[i]), 10))
		line = append(line, strconv.FormatFloat(float64(c.WaitDurations[i]), 'f', -1, 64))
		err := w.Write(line)
		if err != nil {
			return err
		}
	}
	var line []string
	line = append(line, "Total")
	line = append(line, strconv.FormatInt(res.TotalTransferBytes, 10))
	line = append(line, strconv.FormatFloat(res.TotalTransferSeconds, 'f', -1, 64))
	line = append(line, strconv.FormatInt(int64(c.TotalBlockCounts()), 10))
	line = append(line, strconv.FormatFloat(float64(c.TotalWaitDurations()), 'f', -1, 64))
	err = w.Write(line)

	return err
}
