package hqlinterpreter

import (
	"bytes"
	"errors"
	hqlparser "github.com/saichler/hql-parser/golang"
	. "github.com/saichler/hql-schema/golang"
	. "github.com/saichler/utils/golang"
	"reflect"
	"strings"
)

type Query struct {
	tables  map[string]*SchemaNode
	columns map[string]*AttributeID
	where   *Expression
}

func (query *Query) String() string {
	buff := bytes.Buffer{}
	buff.WriteString("Select ")
	first := true

	for _, column := range query.columns {
		if !first {
			buff.WriteString(", ")
		}
		buff.WriteString(column.ID())
		first = false
	}

	buff.WriteString(" From ")

	first = true
	for _, table := range query.tables {
		if !first {
			buff.WriteString(", ")
		}
		buff.WriteString(table.ID())
		first = false
	}

	if query.where != nil {
		buff.WriteString(" Where ")
		buff.WriteString(query.where.String())
	}
	return buff.String()
}

func (query *Query) Tables() map[string]*SchemaNode {
	return query.tables
}

func (query *Query) Columns() map[string]*AttributeID {
	return query.columns
}

func (query *Query) OnlyTopLevel() bool {
	return true
}

func (query *Query) initTables(provider SchemaProvider, pq *hqlparser.Query) error {
	for _, tableName := range pq.Tables() {
		found := false
		for _, name := range provider.Tables() {
			if strings.ToLower(name) == tableName {
				query.tables[tableName], _ = provider.Schema().SchemaNode(name)
				found = true
				break
			}
		}
		if !found {
			return errors.New("Could not find Struct " + tableName + " in Orm Registry.")
		}
	}
	return nil
}

func (query *Query) initColumns(provider SchemaProvider, pq *hqlparser.Query) error {
	mainTable, e := query.MainTable()
	if e != nil {
		return e
	}
	for _, col := range pq.Columns() {
		sf := provider.Schema().CreateAttributeID(mainTable.CreateFeildID(col))
		if sf == nil {
			return errors.New("Cannot find query field: " + col)
		}
		query.columns[col] = sf
	}
	return nil
}

func NewQuery(provider SchemaProvider, sql string) (*Query, error) {

	qp, err := hqlparser.NewQuery(sql)
	if err != nil {
		return nil, err
	}
	ormQuery := &Query{}
	ormQuery.tables = make(map[string]*SchemaNode)
	ormQuery.columns = make(map[string]*AttributeID)

	err = ormQuery.initTables(provider, qp)
	if err != nil {
		return nil, err
	}

	err = ormQuery.initColumns(provider, qp)
	if err != nil {
		return nil, err
	}

	mainTable, err := ormQuery.MainTable()
	if err != nil {
		return nil, err
	}
	if qp.Where()!=nil {
		expr, err := CreateExpression(provider.Schema(), mainTable, qp.Where())
		if err != nil {
			return nil, err
		}
		ormQuery.where = expr
	}
	return ormQuery, nil
}

func (query *Query) MainTable() (*SchemaNode, error) {
	for _, t := range query.tables {
		return t, nil
	}
	return nil, errors.New("No tables in query")
}

func (query *Query) match(value reflect.Value) (bool, error) {
	if !value.IsValid() {
		return false, nil
	}
	if value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return false, nil
		} else {
			value = value.Elem()
		}
	}
	tableName := strings.ToLower(value.Type().Name())
	table := query.tables[tableName]
	if table == nil {
		return false, nil
	}
	return query.where.Match(value)
}

func (query *Query) Filter(list []interface{}) []interface{} {
	result := make([]interface{}, 0)
	for _, i := range list {
		if query.Match(i) {
			result = append(result, i)
		}
	}
	return result
}

func (query *Query) Match(any interface{}) bool {
	val := reflect.ValueOf(any)
	m, e := query.match(val)
	if e != nil {
		Error(e)
	}
	return m
}
