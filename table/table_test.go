package table

import (
	"testing"
	"woqutech.com/cloudctl/cmd"
)

func TestTable(t *testing.T) {
	var data = []cmd.TableService{
		{"phoenix.service", "1234", "active (running)", "3days ago", "Yes"},
		{"phoenix.service", "1234", "active (running)", "3days ago", "Yes"},
		{"phoenix.service", "1234", "active (running)", "3days ago", "Yes"},
		{"phoenix.service", "1234", "active (running)", "3days ago", "Yes"},
	}
	Output(data)
}
