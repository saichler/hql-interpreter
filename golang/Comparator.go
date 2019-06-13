package hqlinterpreter

import (
	"bytes"
	"errors"
	"github.com/saichler/hql-interpreter/golang/comparators"
	hqlparser "github.com/saichler/hql-parser/golang"
	. "github.com/saichler/hql-schema/golang"
	"reflect"
)

type Comparator struct {
	left             string
	leftSchemaField  *AttributeID
	op               hqlparser.ComparatorOperation
	right            string
	rightSchemaField *AttributeID
}

type Comparable interface {
	Compare([]reflect.Value, []reflect.Value) bool
}

var comparables = make(map[hqlparser.ComparatorOperation]Comparable)

func initComparables() {
	if len(comparables) == 0 {
		comparables[hqlparser.Equal] = comparators.NewEqual()
	}
}

func (comparator *Comparator) String() string {
	buff := bytes.Buffer{}
	if comparator.leftSchemaField != nil {
		buff.WriteString(comparator.leftSchemaField.ID())
	} else {
		buff.WriteString(comparator.left)
	}
	buff.WriteString(string(comparator.op))
	if comparator.rightSchemaField != nil {
		buff.WriteString(comparator.rightSchemaField.ID())
	} else {
		buff.WriteString(comparator.right)
	}
	return buff.String()
}

func CreateComparator(schema *Schema, mainTable *SchemaNode, c *hqlparser.Comparator) (*Comparator, error) {
	initComparables()
	ormComp := &Comparator{}
	ormComp.op = c.Operation()
	ormComp.left = c.Left()
	ormComp.right = c.Right()
	ormComp.leftSchemaField = schema.CreateAttributeID(mainTable.CreateFeildID(ormComp.left))
	ormComp.rightSchemaField = schema.CreateAttributeID(mainTable.CreateFeildID(ormComp.right))

	if ormComp.leftSchemaField == nil && ormComp.rightSchemaField == nil {
		return nil, errors.New("No Field was found for comparator:" + c.String())
	}
	return ormComp, nil
}

func (comparator *Comparator) Match(value reflect.Value) (bool, error) {
	var leftValue []reflect.Value
	var rightValue []reflect.Value
	if comparator.leftSchemaField != nil {
		leftValue = comparator.leftSchemaField.ValueOf(value)
	} else {
		leftValue = []reflect.Value{reflect.ValueOf(comparator.left)}
	}
	if comparator.rightSchemaField != nil {
		rightValue = comparator.rightSchemaField.ValueOf(value)
	} else {
		rightValue = []reflect.Value{reflect.ValueOf(comparator.right)}
	}
	matcher := comparables[comparator.op]
	if matcher == nil {
		panic("No Matcher for: " + comparator.op + " operation.")
	}
	return matcher.Compare(leftValue, rightValue), nil
}
