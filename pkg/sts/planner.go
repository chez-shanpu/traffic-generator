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
	"math"
	"strconv"

	"github.com/chez-shanpu/traffic-generator/pkg/option"

	"github.com/chez-shanpu/traffic-generator/pkg/traffic"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Planner struct {
	CycleNum      int
	Seed          uint64
	SendLambda    float64
	SendSeconds   int64
	WaitLambda    float64
	WaitSeconds   int64
	Bitrate       string
	BitrateLambda float64
	BitrateUnit   string
}

func NewPlanner(cfg option.Config) *Planner {
	return &Planner{
		CycleNum:      cfg.Cycle,
		Seed:          cfg.Seed,
		SendLambda:    cfg.SendLambda,
		SendSeconds:   cfg.SendSeconds,
		WaitLambda:    cfg.WaitLambda,
		WaitSeconds:   cfg.WaitSeconds,
		Bitrate:       cfg.Bitrate,
		BitrateLambda: cfg.BitrateLambda,
		BitrateUnit:   cfg.BitrateUnit,
	}
}

func (p *Planner) GenerateTrafficParams() traffic.Params {
	var ts traffic.Params

	bits := p.GenerateBitrates()
	sends := p.GenerateSendSeconds()
	waits := p.GenerateWaitMilliSeconds()

	for i := 0; i < p.CycleNum; i++ {
		t := &traffic.Param{
			Bitrate:          bits[i],
			SendSeconds:      sends[i],
			WaitMilliSeconds: waits[i],
		}
		ts = append(ts, t)
	}
	return ts
}

func (p *Planner) GenerateBitrates() []traffic.Bitrate {
	var bs []traffic.Bitrate

	if len(p.Bitrate) != 0 {
		for i := 0; i < p.CycleNum; i++ {
			bs = append(bs, traffic.Bitrate(p.Bitrate))
		}
	} else {
		bs = p.GenerateRandomBitrates()
	}

	if len(p.BitrateUnit) != 0 {
		var newbs []traffic.Bitrate
		for _, b := range bs {
			b += traffic.Bitrate(p.BitrateUnit)
			newbs = append(newbs, b)
		}
		bs = newbs
	}

	return bs
}

func (p Planner) GenerateRandomBitrates() []traffic.Bitrate {
	ps := distuv.Poisson{
		Lambda: p.BitrateLambda,
		Src:    rand.NewSource(p.Seed),
	}

	var bs []traffic.Bitrate
	for i := 0; i < p.CycleNum; i++ {
		b := math.Round(ps.Rand())
		if b < 1 {
			b = 1
		}
		bs = append(bs, traffic.Bitrate(strconv.FormatFloat(b, 'f', -1, 64)))
	}

	return bs
}

func (p Planner) GenerateSendSeconds() []traffic.Second {
	var ss []traffic.Second

	if p.SendSeconds > 0 {
		for i := 0; i < p.CycleNum; i++ {
			ss = append(ss, traffic.Second(p.SendSeconds))
		}
	} else {
		ss = p.GenerateRandomSendSeconds()
	}

	return ss
}

func (p Planner) GenerateRandomSendSeconds() []traffic.Second {
	ps := distuv.Exponential{
		Rate: p.SendLambda,
		Src:  rand.NewSource(p.Seed),
	}

	var ss []traffic.Second
	for i := 0; i < p.CycleNum; i++ {
		s := traffic.Second(math.Ceil(ps.Rand()))
		ss = append(ss, s)
	}
	return ss
}

func (p Planner) GenerateWaitMilliSeconds() []traffic.MilliSecond {
	var ms []traffic.MilliSecond

	if p.WaitSeconds > 0 {
		for i := 0; i < p.CycleNum-1; i++ {
			ms = append(ms, traffic.MilliSecond(p.WaitSeconds*1000))
		}
	} else {
		ms = p.GenerateRandomWaitMilliSeconds()
	}

	zero := traffic.MilliSecond(0)
	ms = append(ms, zero)
	return ms
}

func (p Planner) GenerateRandomWaitMilliSeconds() []traffic.MilliSecond {
	e := distuv.Exponential{
		Rate: p.WaitLambda,
		Src:  rand.NewSource(p.Seed),
	}

	var ms []traffic.MilliSecond
	for i := 0; i < p.CycleNum-1; i++ {
		m := traffic.MilliSecond(e.Rand() * 1000)
		ms = append(ms, m)
	}
	return ms
}
