package tracer

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/sony/sonyflake"
)

var gen = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: func() (uint16, error) {
		if id := os.ExpandEnv("QMACHINEID"); id != "" {
			if id, err := strconv.ParseUint(id, 10, 16); err == nil {
				return uint16(id), nil
			}
		}

		return uint16(time.Now().UnixNano()), nil
	},
})

func GetId() uint64 {
	for {
		res, err := gen.NextID()
		if err == nil {
			return res
		}

		log.Printf("Cannot generate unique Id: %s", err.Error())
		time.Sleep(0)
	}
}
