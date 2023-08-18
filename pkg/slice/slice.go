package slice

func DoesExist(list []uint, item uint) bool {
	for _, i := range list {
		if i == item {
			return true
		}
	}
	return false
}
