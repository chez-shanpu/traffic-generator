package traffic

import (
	"encoding/csv"
	"os"
	"strconv"
)

type Bitrate string
type Second int64
type MilliSecond int64

type Param struct {
	Bitrate          Bitrate
	SendSeconds      Second
	WaitMilliSeconds MilliSecond
}

type Params []*Param

func (ps Params) Output(out string) error {
	var f *os.File
	var err error

	if out == "" {
		f = os.Stdout
	} else {
		if f, err = os.Create(out); err != nil {
			return err
		}
	}
	return ps.OutputCSV(f)
}

func (ps Params) OutputCSV(f *os.File) error {
	writer := csv.NewWriter(f)
	defer writer.Flush()

	csvHead := []string{"Cycle", "Bitrate", "SendSeconds", "WaitMilliSeconds"}
	if err := writer.Write(csvHead); err != nil {
		return err
	}

	for i, p := range ps {
		var line []string
		line = append(line, strconv.Itoa(i))
		line = append(line, string(p.Bitrate))
		line = append(line, strconv.FormatInt(int64(p.SendSeconds), 10))
		line = append(line, strconv.FormatInt(int64(p.WaitMilliSeconds), 10))
		if err := writer.Write(line); err != nil {
			return err
		}
	}
	return nil
}
