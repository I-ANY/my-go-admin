package tools

import (
	"github.com/jinzhu/copier"
	"time"
)

func GetTime2StrConvert() copier.TypeConverter {
	return copier.TypeConverter{SrcType: time.Time{},
		DstType: "",
		Fn: func(src interface{}) (dst interface{}, err error) {
			d := src.(time.Time)
			if d.IsZero() {
				return "", nil
			}
			return d.Format(time.DateTime), nil
		}}
}

func GetTime2StrPtrConvert() copier.TypeConverter {
	return copier.TypeConverter{SrcType: time.Time{},
		DstType: ToPointer(""),
		Fn: func(src interface{}) (dst interface{}, err error) {
			d := src.(time.Time)
			if d.IsZero() {
				return nil, nil
			}
			return ToPointer(d.Format(time.DateTime)), nil
		}}
}

func WithConverts(converts ...copier.TypeConverter) []copier.TypeConverter {
	return converts
}

func GetPtrTime2PtrStrConvert() copier.TypeConverter {
	return copier.TypeConverter{
		SrcType: &time.Time{},
		DstType: ToPointer(""),
		Fn: func(src interface{}) (interface{}, error) {
			d := src.(*time.Time)
			if d == nil || d.IsZero() {
				return nil, nil
			}
			return ToPointer(d.Format(time.DateTime)), nil
		}}
}
func Time2StrConverts() []copier.TypeConverter {
	return WithConverts(GetTime2StrConvert(), GetTime2StrPtrConvert(), GetPtrTime2PtrStrConvert())
}
