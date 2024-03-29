package comparators

import (
	"reflect"
	"strings"
)

type LessThanOrEqual struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewLessThanOrEqual() *LessThanOrEqual {
	c := &LessThanOrEqual{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = lteqStringMatcher
	c.compares[reflect.Int] = lteqIntMatcher
	c.compares[reflect.Int8] = lteqIntMatcher
	c.compares[reflect.Int16] = lteqIntMatcher
	c.compares[reflect.Int32] = lteqIntMatcher
	c.compares[reflect.Int64] = lteqIntMatcher
	c.compares[reflect.Uint] = lteqUintMatcher
	c.compares[reflect.Uint8] = lteqUintMatcher
	c.compares[reflect.Uint16] = lteqUintMatcher
	c.compares[reflect.Uint32] = lteqUintMatcher
	c.compares[reflect.Uint64] = lteqUintMatcher
	return c
}

func (lt *LessThanOrEqual) Compare(left, right []reflect.Value) bool {
	kind := getKind(left, right)
	comp := lt.compares[kind]
	if comp == nil {
		panic("No LT or EQ Comparator for kind:" + kind.String())
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

func lteqStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zside := removeSingleQuote(strings.ToLower(right.String()))
	return aside <= zside
}

func lteqIntMatcher(left, right reflect.Value) bool {
	aside, ok := getInt64(left)
	if !ok {
		return false
	}
	zside, ok := getInt64(right)
	if !ok {
		return false
	}
	return aside <= zside
}

func lteqUintMatcher(left, right reflect.Value) bool {
	aside, ok := getUint64(left)
	if !ok {
		return false
	}
	zside, ok := getUint64(right)
	if !ok {
		return false
	}
	return aside <= zside
}
