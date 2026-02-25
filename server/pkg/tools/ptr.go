package tools

func ToPointer[T any](t T) *T {
	return &t
}

func ToValue[T any](t *T) T {
	return ToValueWithZero(t)
}

// ToValueWithZero 将指针转换为值，如果指针为空，则返回该类型的零值
func ToValueWithZero[T any](t *T) T {
	if t == nil {
		var zero T
		return zero
	}
	return *t
}

// ToValueWithDefault 将指针转换为值，如果指针为空，则返回默认值def
func ToValueWithDefault[T any](t *T, def T) T {
	if t == nil {
		return def
	}
	return *t
}

func PtrEqPtr[T comparable](a, b *T) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}
func PtrEqVal[T comparable](ptr *T, value T) bool {
	if ptr == nil {
		return false
	}
	return *ptr == value
}
func ValEqPtr[T comparable](value T, ptr *T) bool {
	return PtrEqVal(ptr, value)
}
func IsNil[T any](ptr *T) bool {
	return ptr == nil
}
func IsZero[T comparable](value T) bool {
	var zero T
	return value == zero
}
