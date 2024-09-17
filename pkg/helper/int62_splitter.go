package helper

import "errors"

func SplitInt62(version int64) (int32, int32, error) {
	if version < 0 || (version<<1) < 0 {
		return 0, 0, errors.New("NOT SUPPORTING USAGE OF MORE THAN 62 BITS OF THE 64 INTEGER")
	}
	low := int32(version & 0x7FFF_FFFF)
	high := int32(version >> 31)
	return high, low, nil
}

func MergeInt62(high int32, low int32) (int64, error) {
	if high < 0 || low < 0 {
		return 0, errors.New("NOT SUPPORTING USAGE OF NEGATIVE INTEGERS")
	}

	high64 := int64(high)
	high64 = high64 << 31
	low64 := int64(low)
	value := high64 + low64

	return value, nil
}
