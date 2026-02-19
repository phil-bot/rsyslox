package filters

import (
	"fmt"
	"strings"
	"time"
)

// Builder constructs SQL WHERE clauses and arguments for filtering
type Builder struct {
	conditions []string
	args       []interface{}
}

// New creates a new filter builder
func New() *Builder {
	return &Builder{
		conditions: []string{},
		args:       []interface{}{},
	}
}

// AddDateRange adds a date range filter
func (b *Builder) AddDateRange(start, end time.Time) {
	b.conditions = append(b.conditions, "ReceivedAt BETWEEN ? AND ?")
	b.args = append(b.args, start, end)
}

// AddSeverityFilter adds a severity filter using Priority MOD 8.
//
// Using MOD 8 makes the filter universally correct regardless of whether the
// database uses legacy format (Priority = Severity 0-7) or modern format
// (Priority = RFC PRI = Facility*8 + Severity):
//   - Legacy: Priority is already 0-7, so Priority MOD 8 = Priority ✓
//   - Modern: Priority MOD 8 extracts the Severity component ✓
//   - Mixed:  both cases are handled correctly by the same expression ✓
func (b *Builder) AddSeverityFilter(values []int) {
	if len(values) == 0 {
		return
	}

	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	b.conditions = append(b.conditions,
		fmt.Sprintf("Priority MOD 8 IN (%s)", strings.Join(placeholders, ",")))
	for _, v := range values {
		b.args = append(b.args, v)
	}
}

// AddMultiValueFilter adds a multi-value IN filter for a column
func (b *Builder) AddMultiValueFilter(column string, values []interface{}) {
	if len(values) == 0 {
		return
	}

	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	b.conditions = append(b.conditions,
		fmt.Sprintf("%s IN (%s)", column, strings.Join(placeholders, ",")))
	b.args = append(b.args, values...)
}

// AddStringMultiValue adds a multi-value string filter
func (b *Builder) AddStringMultiValue(column string, values []string) {
	if len(values) == 0 {
		return
	}

	interfaceValues := make([]interface{}, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}

	b.AddMultiValueFilter(column, interfaceValues)
}

// AddIntMultiValue adds a multi-value integer filter
func (b *Builder) AddIntMultiValue(column string, values []int) {
	if len(values) == 0 {
		return
	}

	interfaceValues := make([]interface{}, len(values))
	for i, v := range values {
		interfaceValues[i] = v
	}

	b.AddMultiValueFilter(column, interfaceValues)
}

// AddMessageSearch adds a message search filter with OR logic for multiple terms
func (b *Builder) AddMessageSearch(terms []string) {
	if len(terms) == 0 {
		return
	}

	messageConditions := []string{}
	for _, term := range terms {
		messageConditions = append(messageConditions, "Message LIKE ?")
		b.args = append(b.args, "%"+term+"%")
	}

	b.conditions = append(b.conditions,
		fmt.Sprintf("(%s)", strings.Join(messageConditions, " OR ")))
}

// Build returns the final WHERE clause and arguments
func (b *Builder) Build() (string, []interface{}) {
	if len(b.conditions) == 0 {
		return "1=1", []interface{}{}
	}

	whereClause := strings.Join(b.conditions, " AND ")
	return whereClause, b.args
}

// HasConditions returns true if any conditions have been added
func (b *Builder) HasConditions() bool {
	return len(b.conditions) > 0
}

// ConditionCount returns the number of conditions
func (b *Builder) ConditionCount() int {
	return len(b.conditions)
}
