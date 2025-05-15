package utils

func NilIfZeroUint(input *uint) *uint {
	if input != nil && *input == 0 {
		return nil
	}
	return input
}
