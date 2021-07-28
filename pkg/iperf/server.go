package iperf

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/chez-shanpu/traffic-generator/pkg/tg"
)

func RunServer() (res *tg.Result, err error) {
	args := []string{
		"-s",
		"-1",
		"-J",
	}

	out, err := exec.Command(iperfCmd, args...).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Exec command: %s %s, output: %s, %s", iperfCmd, args, out, err)
	}

	sb, ss, err := parseIperfServerOutput(out)
	res = &tg.Result{
		SendByte:   sb,
		SendSecond: ss,
	}
	return res, err
}

func parseIperfServerOutput(out []byte) (sb int64, ss float64, err error) {
	var i interface{}
	err = json.Unmarshal(out, &i)
	if err != nil {
		return 0, 0, err
	}

	sb = calcReceiveBytesFromResultJson(i)
	ss = i.(map[string]interface{})["end"].(map[string]interface{})["sum"].(map[string]interface{})["seconds"].(float64)
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
