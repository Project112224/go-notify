package extension

import (
	"log"
	"time"
)

type DateString string

func (ds DateString) FormatToHMS(format ...string) string {
	finalFormat := "15:04:05"
	if len(format) > 0 {
		finalFormat = format[0]
	}

	t, err := time.Parse(time.RFC3339Nano, string(ds))
	if err != nil {
		log.Println("時間解析失敗:", err)
		return string(ds)
	}
	return t.Format(finalFormat)
}
