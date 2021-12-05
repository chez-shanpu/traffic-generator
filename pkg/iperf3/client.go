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

package iperf3

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/chez-shanpu/traffic-generator/pkg/option"

	"github.com/chez-shanpu/traffic-generator/pkg/traffic"

	"github.com/gocarina/gocsv"
)

type Param struct {
	Bitrate          traffic.Bitrate     `csv:"Bitrate"`
	SendSeconds      traffic.Second      `csv:"SendSeconds"`
	WaitMilliSeconds traffic.MilliSecond `csv:"WaitMilliSeconds"`
}

type Client struct {
	Params             []*Param
	DstAddr            string
	DstPort            string
	MaximumSegmentSize int64
	UdpFlag            bool
	IPv6Flag           bool
	Flowlabel          int64
}

func NewIperfClientFromParamsFile(cfg option.Config) (*Client, error) {
	ps, err := parseParamsFile(cfg.Param)
	if err != nil {
		return nil, err
	}

	return NewIperfClient(cfg, ps), nil
}

func NewIperfClient(cfg option.Config, params []*Param) *Client {
	return &Client{
		DstAddr:            cfg.DstAddr,
		DstPort:            cfg.DstPort,
		MaximumSegmentSize: cfg.Mss,
		UdpFlag:            cfg.UDP,
		IPv6Flag:           cfg.IPv6,
		Flowlabel:          cfg.Flowlabel,
		Params:             params,
	}
}

func parseParamsFile(paramFilePath string) ([]*Param, error) {
	f, err := os.Open(paramFilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var params []*Param
	err = gocsv.UnmarshalFile(f, &params)
	return params, err
}

func (c Client) GenerateTraffic() (traffic.Results, error) {
	var rs traffic.Results

	for i, p := range c.Params {
		fmt.Printf("Run %d: Bitrate %s, SendSeconds %d\n", i, p.Bitrate, p.SendSeconds)
		r, err := c.execIperf3(c.makeIperf3Args(p))
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
		fmt.Printf("sleep %d msec\n", p.WaitMilliSeconds)
		time.Sleep(time.Duration(p.WaitMilliSeconds) * time.Millisecond)
	}
	return rs, nil
}

func (c Client) OutputResults(rs traffic.Results, out string) error {
	var f *os.File
	var err error

	if out == "" {
		f = os.Stdout
	} else {
		if f, err = os.Create(out); err != nil {
			return err
		}
	}

	return c.OutputResultsCSV(rs, f)
}

func (c Client) OutputResultsCSV(rs traffic.Results, f *os.File) error {
	w := csv.NewWriter(f)
	defer w.Flush()

	csvHead := []string{"Cycle", "SendByte", "Bitrate", "SendSecond", "WaitMilliSecond"}
	if err := w.Write(csvHead); err != nil {
		return err
	}

	for i, r := range rs {
		var line []string
		line = append(line, strconv.Itoa(i))
		line = append(line, strconv.FormatInt(r.SendByte, 10))
		line = append(line, string(c.Params[i].Bitrate))
		line = append(line, strconv.FormatFloat(r.SendSecond, 'f', -1, 64))
		line = append(line, strconv.FormatInt(int64(c.Params[i].WaitMilliSeconds), 10))
		if err := w.Write(line); err != nil {
			return err
		}
	}
	var line []string
	line = append(line, "Total")
	line = append(line, strconv.FormatInt(rs.TotalSendBytes(), 10))
	line = append(line, "-")
	line = append(line, strconv.FormatInt(int64(c.TotalSendSeconds()), 10))
	line = append(line, strconv.FormatFloat(float64(c.TotalWaitMilliSeconds()), 'f', -1, 64))
	return w.Write(line)
}

func (c Client) TotalSendSeconds() traffic.Second {
	res := traffic.Second(0)
	for _, p := range c.Params {
		res += p.SendSeconds
	}
	return res
}

func (c Client) TotalWaitMilliSeconds() traffic.MilliSecond {
	res := traffic.MilliSecond(0)
	for _, p := range c.Params {
		res += p.WaitMilliSeconds
	}
	return res
}

func (c Client) makeIperf3Args(p *Param) []string {
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
	if c.UdpFlag {
		args = append(args, "-u")
	}
	if c.IPv6Flag {
		args = append(args, "-6")
	}
	if c.Flowlabel > 0 {
		args = append(args, "-L")
		args = append(args, strconv.FormatInt(c.Flowlabel, 10))
	}
	return args
}

func (c Client) execIperf3(args []string) (res *traffic.Result, err error) {
	out, err := exec.Command(iperf3, args...).CombinedOutput()
	if err != nil {
		fmt.Printf("[ERROR] Exec command: %s %s, output: %s, %s", iperf3, args, out, err)
		return &traffic.Result{
			SendByte:   0,
			SendSecond: 0,
		}, nil
	}

	sb, ss, err := parseIperfOutput(out)
	res = &traffic.Result{
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
