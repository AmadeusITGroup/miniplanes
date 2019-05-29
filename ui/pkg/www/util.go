package www

import (
	"fmt"
	"strconv"
	"strings"
)

func stringToInt32(v string) (int32, error) {
	dest, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, nil
	}
	return int32(dest), err
}

func stringToInt64(v string) (int64, error) {
	dest, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0, nil
	}
	return dest, err
}

func stringToFloat64(v string) (float64, error) {
	dest, err := strconv.ParseFloat(v, 64)
	if err != nil {
		return 0, nil
	}
	return dest, err
}

func stringToBool(v string) (bool, error) {
	if isTrue := strings.EqualFold(v, "true"); isTrue {
		return true, nil
	}
	if isFalse := strings.EqualFold(v, "false"); isFalse {
		return false, nil
	}
	return false, fmt.Errorf("supplied value %q, admitted values are 'true' or 'false'", v)
}
