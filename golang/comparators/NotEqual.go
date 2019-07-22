package comparators

import (
	"reflect"
	"strings"
)

type NotEqual struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewNotEqual() *NotEqual {
	c := &NotEqual{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = noteqStringMatcher
	c.compares[reflect.Int] = noteqIntMatcher
	c.compares[reflect.Int8] = noteqIntMatcher
	c.compares[reflect.Int16] = noteqIntMatcher
	c.compares[reflect.Int32] = noteqIntMatcher
	c.compares[reflect.Int64] = noteqIntMatcher
	c.compares[reflect.Uint] = noteqUintMatcher
	c.compares[reflect.Uint8] = noteqUintMatcher
	c.compares[reflect.Uint16] = noteqUintMatcher
	c.compares[reflect.Uint32] = noteqUintMatcher
	c.compares[reflect.Uint64] = noteqUintMatcher
	return c
}

func (notequal *NotEqual) Compare(left, right []reflect.Value) bool {
	kind := getKind(left, right)
	comp := notequal.compares[kind]
	if comp == nil {
		panic("No NotEqual Comparator for kind:" + kind.String())
	}
	for _, aside := range left {
		for _, zside := range right {
			if comp(aside, zside) {
				return true
			}
		}
	}
	return false
}

func noteqStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zside := removeSingleQuote(strings.ToLower(right.String()))
	return aside != zside
}

func noteqIntMatcher(left, right reflect.Value) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}
	zside, ok := getInt64(right)
	if !ok {
		return false
	}
	return aside != zside
}

func noteqUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside != zside
}
