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
	"github.com/chez-shanpu/traffic-generator/pkg/tg"
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

type Planner struct {
	CycleNum int
	Seed     uint64
	Lambda   float64
	Rate     float64
}

func NewPlanner(c int, s uint64, l, r float64) *Planner {
	return &Planner{
		CycleNum: c,
		Seed:     s,
		Lambda:   l,
		Rate:     r,
	}
}

func (p Planner) CalcPacketCounts() []tg.PacketCount {
	ps := distuv.Poisson{
		Lambda: p.Lambda,
		Src:    rand.NewSource(p.Seed),
	}

	var pcs []tg.PacketCount
	for i := 0; i < p.CycleNum; i++ {
		pc := ps.Rand()
		pcs = append(pcs, tg.PacketCount(pc))
	}
	return pcs
}

func (p Planner) CalcWaitDurations() []tg.WaitDuration {
	e := distuv.Exponential{
		Rate: p.Rate,
		Src:  rand.NewSource(p.Seed),
	}

	var wds []tg.WaitDuration
	for i := 0; i < p.CycleNum-1; i++ {
		wd := tg.WaitDuration(e.Rand() * 1000)
		wds = append(wds, wd)
	}
	wds = append(wds, 0)
	return wds
}
