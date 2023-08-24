package slice

func DoesExist(list []uint, item uint) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}

func MapUint64ToUint(l []uint64) []uint {
	r := make([]uint, 0)

	for _, item := range l {
		r = append(r, uint(item))
	}

	return r
}

func MapUintToUint64(l []uint) []uint64 {
	r := make([]uint64, 0)

	for _, item := range l {
		r = append(r, uint64(item))
	}

	return r
}
