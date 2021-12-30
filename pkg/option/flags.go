package option

import "github.com/spf13/viper"

const (
	Bitrate       = "bitrate"
	BitrateLambda = "bitrate-lambda"
	BitrateUnit   = "bitrate-unit"
	Cycle         = "cycle"
	DstAddr       = "dst-addr"
	DstPort       = "dst-port"
	Flowlabel     = "flowlabel"
	IPv6          = "ipv6"
	Mss           = "mss"
	Out           = "out"
	Param         = "param"
	Seed          = "seed"
	SendLambda    = "send-lambda"
	SendSeconds   = "send-seconds"
	UDP           = "udp"
	WaitLambda    = "wait-lambda"
	WaitSeconds   = "wait-seconds"
	WindowSize    = "window"
)

type Config struct {
	Bitrate       float64
	BitrateLambda float64
	BitrateUnit   string
	Cycle         int
	DstAddr       string
	DstPort       string
	Flowlabel     int64
	IPv6          bool
	Mss           int64
	Out           string
	Param         string
	Seed          uint64
	SendLambda    float64
	SendSeconds   int64
	UDP           bool
	WaitLambda    float64
	WaitSeconds   int64
	WindowSize    string
}

func (c *Config) Populate() {
	c.Bitrate = viper.GetFloat64(Bitrate)
	c.BitrateLambda = viper.GetFloat64(BitrateLambda)
	c.BitrateUnit = viper.GetString(BitrateUnit)
	c.Cycle = viper.GetInt(Cycle)
	c.DstAddr = viper.GetString(DstAddr)
	c.DstPort = viper.GetString(DstPort)
	c.Flowlabel = viper.GetInt64(Flowlabel)
	c.IPv6 = viper.GetBool(IPv6)
	c.Mss = viper.GetInt64(Mss)
	c.Out = viper.GetString(Out)
	c.Param = viper.GetString(Param)
	c.Seed = viper.GetUint64(Seed)
	c.SendLambda = viper.GetFloat64(SendLambda)
	c.SendSeconds = viper.GetInt64(SendSeconds)
	c.UDP = viper.GetBool(UDP)
	c.WaitLambda = viper.GetFloat64(WaitLambda)
	c.WaitSeconds = viper.GetInt64(WaitSeconds)
	c.WindowSize = viper.GetString(WindowSize)
}
