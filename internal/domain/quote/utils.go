package quote

import "fmt"

func CheckIsValidTimeframe(timeframeSec int) error {
	if timeframeSec != 0 && 86400%timeframeSec != 0 {
		return fmt.Errorf("one day is not a multiple of requested timeframe")
	}

	return nil
}
