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
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/chez-shanpu/traffic-generator/pkg/tg"
	"github.com/gocarina/gocsv"
)

const iperfCmd = "iperf3"

type Param struct {
	Bitrate          tg.Bitrate          `csv:"Bitrate"`
	SendSeconds      tg.SendSeconds      `csv:"SendSeconds"`
	WaitMilliSeconds tg.WaitMilliSeconds `csv:"WaitMilliSeconds"`
}

type Client struct {
	Params             []*Param
	DstAddr            string
	DstPort            string
	MaximumSegmentSize int64
}

func NewIperfClientFromParamsFile(da, dp string, mss int64, paramFile string) (*Client, error) {
	ps, err := parsePramsFile(paramFile)
	if err != nil {
		return nil, err
	}

	c := NewIperfClient(da, dp, mss, ps)
	return c, err
}

func parsePramsFile(paramFilePath string) ([]*Param, error) {
	pFile, err := os.Open(paramFilePath)
	if err != nil {
		return nil, err
	}
	defer pFile.Close()

	var params []*Param
	err = gocsv.UnmarshalFile(pFile, &params)
	return params, err
}

func NewIperfClient(dstAddr, dstPort string, mss int64, params []*Param) *Client {
	return &Client{
		DstAddr:            dstAddr,
		DstPort:            dstPort,
		MaximumSegmentSize: mss,
		Params:             params,
	}
}

func (c Client) GenerateTraffic() (tg.Results, error) {
	var rs tg.Results
	for i, p := range c.Params {
		fmt.Printf("Run %d: Bitrate %s, SendSeconds %d\n", i, p.Bitrate, p.SendSeconds)
		r, err := c.execIperf(p)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
		fmt.Printf("sleep %d msec\n", p.WaitMilliSeconds)
		time.Sleep(time.Duration(p.WaitMilliSeconds) * time.Millisecond)
	}

	return rs, nil
}

func (c Client) TotalSendSeconds() tg.SendSeconds {
	var res tg.SendSeconds
	for _, p := range c.Params {
		res += p.SendSeconds
	}
	return res
}

func (c Client) TotalWaitMilliSeconds() tg.WaitMilliSeconds {
	var res tg.WaitMilliSeconds
	for _, p := range c.Params {
		res += p.WaitMilliSeconds
	}
	return res
}

func (c Client) execIperf(p *Param) (res *tg.Result, err error) {
	args := []string{
		"-c",
		c.DstAddr,
		"-t", strconv.FormatInt(int64(p.SendSeconds), 10),
		"-b", string(p.Bitrate),
		"-J",
	}
	if c.DstPort != "" {
		args = append(args, "-p")
		args = append(args, c.DstPort)
	}
	if c.MaximumSegmentSize != 0 {
		args = append(args, "-M")
		args = append(args, strconv.FormatInt(c.MaximumSegmentSize, 10))
	}

	out, err := exec.Command(iperfCmd, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("[ERROR] Exec command: %s %s, output: %s, %s", iperfCmd, args, out, err)
		return &tg.Result{
			SendByte:   0,
			SendSecond: -1,
		}, nil
	}

	sb, ss, err := parseIperfOutput(out)
	res = &tg.Result{
		SendByte:   sb,
		SendSecond: ss,
	}
	return res, err
}

func parseIperfOutput(out []byte) (sb int64, ss float64, err error) {
	var i interface{}
	err = json.Unmarshal(out, &i)
	if err != nil {
		return 0, 0, err
	}

	sb = int64(i.(map[string]interface{})["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["bytes"].(float64))
	ss = i.(map[string]interface{})["end"].(map[string]interface{})["sum_sent"].(map[string]interface{})["seconds"].(float64)

	return sb, ss, nil
}

func (c Client) OutputResults(rs tg.Results, out string) error {
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

	err = c.outputResultsCSV(rs, f)
	return err
}

func (c Client) outputResultsCSV(rs tg.Results, f *os.File) error {
	w := csv.NewWriter(f)
	defer w.Flush()

	csvHead := []string{"Cycle", "SendByte", "Bitrate", "SendSecond", "WaitMilliSecond"}
	err := w.Write(csvHead)
	if err != nil {
		return err
	}

	for i, r := range rs {
		var line []string
		line = append(line, strconv.Itoa(i))
		line = append(line, strconv.FormatInt(r.SendByte, 10))
		line = append(line, string(c.Params[i].Bitrate))
		line = append(line, strconv.FormatFloat(r.SendSecond, 'f', -1, 64))
		line = append(line, strconv.FormatInt(int64(c.Params[i].WaitMilliSeconds), 10))
		err := w.Write(line)
		if err != nil {
			return err
		}
	}
	var line []string
	line = append(line, "Total")
	line = append(line, strconv.FormatInt(rs.TotalSendBytes(), 10))
	line = append(line, "-")
	line = append(line, strconv.FormatInt(int64(c.TotalSendSeconds()), 10))
	line = append(line, strconv.FormatFloat(float64(c.TotalWaitMilliSeconds()), 'f', -1, 64))
	err = w.Write(line)

	return err
}
