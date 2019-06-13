package hqlinterpreter

import (
	"bytes"
	"errors"
	hqlparser "github.com/saichler/hql-parser/golang"
	. "github.com/saichler/hql-schema/golang"
	"reflect"
)

type Condition struct {
	comparator *Comparator
	op         hqlparser.ConditionOperation
	next       *Condition
}

func CreateCondition(schema *Schema, mainTable *SchemaNode, c *hqlparser.Condition) (*Condition, error) {
	ormCond := &Condition{}
	ormCond.op = c.Operation()
	comp, e := CreateComparator(schema, mainTable, c.Comparator())
	if e != nil {
		return nil, e
	}
	ormCond.comparator = comp
	if c.Next() != nil {
		next, e := CreateCondition(schema, mainTable, c.Next())
		if e != nil {
			return nil, e
		}
		ormCond.next = next
	}
	return ormCond, nil
}

func (condition *Condition) String() string {
	buff := &bytes.Buffer{}
	buff.WriteString("(")
	condition.toString(buff)
	buff.WriteString(")")
	return buff.String()
}

func (condition *Condition) toString(buff *bytes.Buffer) {
	if condition.comparator != nil {
		buff.WriteString(condition.comparator.String())
	}
	if condition.next != nil {
		buff.WriteString(string(condition.op))
		condition.next.toString(buff)
	}
}

func (condition *Condition) Match(value reflect.Value) (bool, error) {
	comp, e := condition.comparator.Match(value)
	if e != nil {
		return false, e
	}
	next := true
	if condition.op == hqlparser.Or {
		next = false
	}
	if condition.next != nil {
		next, e = condition.next.Match(value)
		if e != nil {
			return false, e
		}
	}
	if condition.op == "" {
		return next && comp, nil
	}
	if condition.op == hqlparser.And {
		return comp && next, nil
	}
	if condition.op == hqlparser.Or {
		return comp || next, nil
	}
	return false, errors.New("Unsupported operation in match:" + string(condition.op))
}
