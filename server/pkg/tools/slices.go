package tools

// RemoveDuplication [T any]
//
//	@Description: 删除重复的值，并且保持原有的顺序
//	@param arr
//	@return []T
func RemoveDuplication[T comparable](arr []T) []T {
	return RemoveDuplicationByKey(arr, func(element T) T {
		return element
	})
}

func RemoveDuplicationByKey[T any, K comparable](arr []T, getKeyFunc func(element T) K) []T {
	set := make(map[K]struct{})
	res := make([]T, 0, len(arr))
	for _, v := range arr {
		key := getKeyFunc(v)
		_, ok := set[key]
		if ok {
			continue
		}
		set[key] = struct{}{}
		res = append(res, v)
	}
	return res
}

func GetSlice[T, V any](elements []V, fn func(e V) T) []T {
	var v = make([]T, 0, len(elements))
	for _, e := range elements {
		v = append(v, fn(e))
	}
	return v
}
func GetNoRepeatSlice[T any, V comparable](elements []T, getSliceElementFn func(e T) V) []V {
	var v = make([]V, 0, len(elements))
	var existValues = make(map[V]struct{})
	for _, e := range elements {
		key := getSliceElementFn(e)
		// 已存在，不添加
		if _, ok := existValues[key]; ok {
			continue
		}
		v = append(v, key)
		existValues[key] = struct{}{}
	}
	return v
}

type TP interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int64 | int32 | int16 | int8 | string | float32 | float64
}

func InSlice[T comparable](el T, arr []T) bool {
	for _, v := range arr {
		if v == el {
			return true
		}
	}
	return false
}

func Slice2Map[K comparable, V any](elements []V, getKeyFunc func(element V) K) map[K]V {
	var res = make(map[K]V)
	for _, element := range elements {
		res[getKeyFunc(element)] = element
	}
	return res
}

func Slice2MapSlice[K comparable, V any](elements []V, getKeyFunc func(element V) K) map[K][]V {
	var res = make(map[K][]V)
	for _, element := range elements {
		key := getKeyFunc(element)
		res[key] = append(res[key], element)
	}
	return res
}

func FullEqualFunc[V any, T comparable](items []V, getValue func(item V) T) bool {
	if len(items) <= 1 {
		return true
	}
	firstValue := getValue(items[0])
	for _, info := range items[1:] {
		if getValue(info) != firstValue {
			return false
		}
	}
	return true
}

func Intersect[T comparable](a []T, b []T) []T {
	inter := make([]T, 0)
	mp := make(map[T]bool)

	for _, s := range a {
		if _, ok := mp[s]; !ok {
			mp[s] = true
		}
	}
	for _, s := range b {
		if _, ok := mp[s]; ok {
			inter = append(inter, s)
		}
	}

	return inter
}
func GetMap[T any, K comparable](items []T, getKeyFunc func(element T) K) map[K]struct{} {
	mp := make(map[K]struct{})
	for _, item := range items {
		mp[getKeyFunc(item)] = struct{}{}
	}
	return mp
}
func SliceFilter[T any](arr []T, predicate func(i T) bool) []T {
	var result = make([]T, 0, len(arr))
	for _, v := range arr {
		if predicate(v) {
			result = append(result, v)
		}
	}
	return result
}

type Slice[T any, K comparable] []T

func NewSlice[T any, K comparable](data []T) *Slice[T, K] {
	return (*Slice[T, K])(&data)
}
func (s *Slice[T, K]) Filter(predicate func(T) bool) *Slice[T, K] {
	filterd := SliceFilter(*s, predicate)
	return (*Slice[T, K])(&filterd)
}
func (s *Slice[T, K]) RemoveDuplicationByKey(getKeyFunc func(element T) K) *Slice[T, K] {
	removed := RemoveDuplicationByKey(*s, getKeyFunc)
	return (*Slice[T, K])(&removed)
}
func (s *Slice[T, K]) Data() []T {
	return *s
}

func SliceFullEqual[T comparable](a []T, b []T) bool {
	if a == nil && b != nil {
		return false
	}
	if a != nil && b == nil {
		return false
	}
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// ChunkSlice 将切片按指定大小分块
func ChunkSlice[T any](items []T, size int) [][]T {
	if size <= 0 {
		return [][]T{items}
	}
	n := len(items)
	if n == 0 {
		return nil
	}
	chunks := make([][]T, 0, (n+size-1)/size)
	for i := 0; i < n; i += size {
		end := i + size
		if end > n {
			end = n
		}
		chunks = append(chunks, items[i:end])
	}
	return chunks
}

// SliceDiff 返回两个切片的差异,无序
func SliceDiff[T comparable](s1 []T, s2 []T) (onlyInS1 []T, onlyInS2 []T) {
	s1Map := GetMap(s1, func(e T) T { return e })
	s2Map := GetMap(s2, func(e T) T { return e })
	// s1中有但s2中没有的
	for v, _ := range s1Map {
		if _, exists := s2Map[v]; !exists {
			onlyInS1 = append(onlyInS1, v)
		}
	}
	// s2中有但s1中没有的
	for v, _ := range s2Map {
		if _, exists := s1Map[v]; !exists {
			onlyInS2 = append(onlyInS2, v)
		}
	}
	return onlyInS1, onlyInS2
}
