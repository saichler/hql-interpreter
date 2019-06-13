package comparators

import (
	"reflect"
	"strconv"
	"strings"
)

type Equal struct {
	compares map[reflect.Kind]func(reflect.Value, reflect.Value) bool
}

func NewEqual() *Equal {
	c := &Equal{}
	c.compares = make(map[reflect.Kind]func(reflect.Value, reflect.Value) bool)
	c.compares[reflect.String] = eqStringMatcher
	c.compares[reflect.Int] = eqIntMatcher
	c.compares[reflect.Int8] = eqIntMatcher
	c.compares[reflect.Int16] = eqIntMatcher
	c.compares[reflect.Int32] = eqIntMatcher
	c.compares[reflect.Int64] = eqIntMatcher
	c.compares[reflect.Uint] = eqUintMatcher
	c.compares[reflect.Uint8] = eqUintMatcher
	c.compares[reflect.Uint16] = eqUintMatcher
	c.compares[reflect.Uint32] = eqUintMatcher
	c.compares[reflect.Uint64] = eqUintMatcher
	return c
}

func (equal *Equal) Compare(left, right []reflect.Value) bool {
	kind := getKind(left, right)
	comp := equal.compares[kind]
	if comp == nil {
		panic("No Equal Comparator for kind:" + kind.String())
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

func removeSingleQuote(value string) string {
	if strings.Contains(value, "'") {
		return value[1 : len(value)-1]
	}
	return value
}

func eqStringMatcher(left, right reflect.Value) bool {
	aside := removeSingleQuote(strings.ToLower(left.String()))
	zside := removeSingleQuote(strings.ToLower(right.String()))
	return aside == zside
}

func eqIntMatcher(left, right reflect.Value) bool {
	var aside int64
	var zside int64
	if left.Kind() != reflect.String {
		aside = left.Int()
	} else {
		i, e := strconv.Atoi(left.String())
		if e != nil {
			return false
		}
		aside = int64(i)
	}

	if right.Kind() != reflect.String {
		aside = right.Int()
	} else {
		i, e := strconv.Atoi(right.String())
		if e != nil {
			return false
		}
		zside = int64(i)
	}

	return aside == zside
}

func eqUintMatcher(left, right reflect.Value) bool {
	var aside uint64
	var zside uint64
	if left.Kind() != reflect.String {
		aside = left.Uint()
	} else {
		i, e := strconv.Atoi(left.String())
		if e != nil {
			return false
		}
		aside = uint64(i)
	}

	if right.Kind() != reflect.String {
		aside = right.Uint()
	} else {
		i, e := strconv.Atoi(right.String())
		if e != nil {
			return false
		}
		zside = uint64(i)
	}

	return aside == zside
}

func getKind(aside, zside []reflect.Value) reflect.Kind {
	aSideKind := reflect.String
	zSideKind := reflect.String
	if len(aside) > 0 {
		aSideKind = aside[0].Kind()
	}
	if len(zside) > 0 {
		zSideKind = zside[0].Kind()
	}
	if aSideKind != reflect.String {
		return aSideKind
	} else if zSideKind != reflect.String {
		return zSideKind
	}
	return aSideKind
}
