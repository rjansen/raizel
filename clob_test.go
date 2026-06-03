package raizel

import (
	"strings"
	"testing"

	go_ora "github.com/sijms/go-ora/v2"
)

func TestPromoteOracleClobs(t *testing.T) {
	big := strings.Repeat("x", oracleVarchar2BindLimit+1)
	atLimit := strings.Repeat("y", oracleVarchar2BindLimit)

	args := []any{
		1,               // non-string: untouched
		"short",         // small string: stays VARCHAR2
		atLimit,         // exactly at the limit: stays VARCHAR2
		big,             // over the limit: promoted to CLOB
		[]byte("bytes"), // non-string: untouched
	}
	promoteOracleClobs(args)

	if _, ok := args[0].(int); !ok {
		t.Errorf("args[0]: int was altered: %T", args[0])
	}
	if s, ok := args[1].(string); !ok || s != "short" {
		t.Errorf("args[1]: short string was altered: %#v", args[1])
	}
	if s, ok := args[2].(string); !ok || len(s) != oracleVarchar2BindLimit {
		t.Errorf("args[2]: at-limit string should stay VARCHAR2, got %T", args[2])
	}
	clob, ok := args[3].(go_ora.Clob)
	if !ok {
		t.Fatalf("args[3]: oversized string not promoted to go_ora.Clob, got %T", args[3])
	}
	if !clob.Valid || len(clob.String) != oracleVarchar2BindLimit+1 {
		t.Errorf("args[3]: clob content wrong: valid=%v len=%d", clob.Valid, len(clob.String))
	}
	if _, ok := args[4].([]byte); !ok {
		t.Errorf("args[4]: []byte was altered: %T", args[4])
	}
}
