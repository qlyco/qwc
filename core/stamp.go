package core

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/qlyco/qwc/internal/res"
)

type buildStamp struct {
	path   string
	logger *log.Logger
}

func (stamp buildStamp) validFrom() int64 {
	source, readErr := os.ReadFile(stamp.path)
	value, parseErr := strconv.ParseInt(strings.TrimSpace(string(source)), res.TimestampBase, res.TimestampBits)

	if readErr != nil || parseErr != nil {
		return stamp.renew()
	}

	return value
}

func (stamp buildStamp) renew() int64 {
	now := time.Now()
	err := writeAtomic(stamp.path, strings.NewReader(strconv.FormatInt(now.UnixNano(), res.TimestampBase)), now)

	if err != nil {
		stamp.logger.Printf(res.FileLogFormat, stamp.path, err)
	} else {
		stamp.logger.Print(res.ReloadedMessage)
	}

	return now.UnixNano()
}
