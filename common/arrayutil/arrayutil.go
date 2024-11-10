package arrayutil

// Contains 判断数组是否包含某个值
// 支持基本数据类型: string, int, int64, float64 等
func Contains[T comparable](arr []T, target T) bool {
	for _, item := range arr {
		if item == target {
			return true
		}
	}
	return false
}

// ContainsAny 判断数组是否包含另一个数组中的任意一个值
func ContainsAny[T comparable](arr []T, targets []T) bool {
	for _, target := range targets {
		if Contains(arr, target) {
			return true
		}
	}
	return false
}

// ContainsAll 判断数组是否包含另一个数组中的所有值
func ContainsAll[T comparable](arr []T, targets []T) bool {
	for _, target := range targets {
		if !Contains(arr, target) {
			return false
		}
	}
	return true
}
