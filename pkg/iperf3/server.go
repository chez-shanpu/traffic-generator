package iperf3

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/chez-shanpu/traffic-generator/pkg/file"

	"github.com/chez-shanpu/traffic-generator/pkg/traffic"
)

type Server struct{}

func NewServer() *Server {
	return &Server{}
}

func (s *Server) Run() (res *traffic.Result, err error) {
	args := makeServerArgs()
	out, err := exec.Command(iperf3, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Exec command: %s %s, output: %s, %s", iperf3, args, out, err)
	}

	sb, ss, err := parseIperfServerOutput(out)
	res = &traffic.Result{
		SendByte:   sb,
		SendSecond: ss,
	}
	return res, err
}

func (s *Server) OutputResult(ress traffic.Results, out string) error {
	f, err := file.Create(out)
	if err != nil {
		return err
	}
	return s.OutputResultCSV(ress, f)
}

func (s *Server) OutputResultCSV(ress traffic.Results, f *os.File) error {
	w := csv.NewWriter(f)
	defer w.Flush()

	var totalReceiveByte int64
	var totalSendSeconds float64
	for _, res := range ress {
		totalReceiveByte += res.SendByte
		totalSendSeconds += res.SendSecond
	}

	csvHead := []string{"TotalReceiveBytes", "SendSeconds"}
	if err := w.Write(csvHead); err != nil {
		return err
	}

	var line []string
	line = append(line, strconv.FormatInt(totalReceiveByte, 10))
	line = append(line, strconv.FormatFloat(totalSendSeconds, 'f', -1, 64))
	return w.Write(line)
}

func makeServerArgs() []string {
	return []string{
		"-s",
		"-1",
		"-J",
	}
}

func parseIperfServerOutput(out []byte) (sb int64, ss float64, err error) {
	var i interface{}

	if err = json.Unmarshal(out, &i); err != nil {
		return 0, 0, err
	}

	sb = calcReceiveBytesFromResultJson(i)
	if _, ok := i.(map[string]interface{})["end"].(map[string]interface{})["sum_received"].(map[string]interface{}); ok {
		// tcp
		ss = i.(map[string]interface{})["end"].(map[string]interface{})["sum_received"].(map[string]interface{})["seconds"].(float64)
	} else {
		// udp
		ss = i.(map[string]interface{})["end"].(map[string]interface{})["sum"].(map[string]interface{})["seconds"].(float64)
	}
	return
}

func calcReceiveBytesFromResultJson(i interface{}) (bytes int64) {
	intervals := i.(map[string]interface{})["intervals"].([]interface{})
	for _, interval := range intervals {
		b := interval.(map[string]interface{})["sum"].(map[string]interface{})["bytes"].(float64)
		bytes += int64(b)
	}
	return bytes
}
