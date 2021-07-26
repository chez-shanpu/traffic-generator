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

package sts

import (
	"encoding/csv"
	"os"
	"strconv"

	"github.com/chez-shanpu/traffic-generator/pkg/tg"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Planner struct {
	CycleNum   int
	Seed       uint64
	SendLambda float64
	WaitLambda float64
	Bitrate    tg.Bitrate
}

func NewPlanner(c int, s uint64, sl, wl float64, b string) *Planner {
	return &Planner{
		CycleNum:   c,
		Seed:       s,
		SendLambda: sl,
		WaitLambda: wl,
		Bitrate:    tg.Bitrate(b),
	}
}

func (p Planner) CalcBitrates() []*tg.Bitrate {
	var bs []*tg.Bitrate
	for i := 0; i < p.CycleNum; i++ {
		b := p.Bitrate
		bs = append(bs, &b)
	}
	return bs
}

func (p Planner) CalcSendSeconds() []*tg.SendSeconds {
	ps := distuv.Exponential{
		Rate: p.SendLambda,
		Src:  rand.NewSource(p.Seed),
	}

	var sds []*tg.SendSeconds
	for i := 0; i < p.CycleNum; i++ {
		sd := tg.SendSeconds(ps.Rand())
		sds = append(sds, &sd)
	}
	return sds
}

func (p Planner) CalcWaitMilliSeconds() []*tg.WaitMilliSeconds {
	e := distuv.Exponential{
		Rate: p.WaitLambda,
		Src:  rand.NewSource(p.Seed),
	}

	var wds []*tg.WaitMilliSeconds
	for i := 0; i < p.CycleNum-1; i++ {
		wd := tg.WaitMilliSeconds(e.Rand() * 1000)
		wds = append(wds, &wd)
	}
	zero := tg.WaitMilliSeconds(0)
	wds = append(wds, &zero)
	return wds
}

func (p Planner) OutputParams(b []*tg.Bitrate, s []*tg.SendSeconds, w []*tg.WaitMilliSeconds, out string) error {
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

	err = p.outputParamsCSV(b, s, w, f)
	return err
}

func (p Planner) outputParamsCSV(b []*tg.Bitrate, s []*tg.SendSeconds, w []*tg.WaitMilliSeconds, f *os.File) error {
	writer := csv.NewWriter(f)
	defer writer.Flush()

	csvHead := []string{"Cycle", "Bitrate", "SendSeconds", "WaitMilliSeconds"}
	err := writer.Write(csvHead)
	if err != nil {
		return err
	}

	for i := 0; i < p.CycleNum; i++ {
		var line []string
		line = append(line, strconv.Itoa(i))
		line = append(line, string(*b[i]))
		line = append(line, strconv.FormatInt(int64(*s[i]), 10))
		line = append(line, strconv.FormatInt(int64(*w[i]), 10))
		err = writer.Write(line)
		if err != nil {
			return err
		}
	}
	return err
}
