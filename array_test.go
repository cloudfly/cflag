package confo

import (
	"flag"
	"os"
	"reflect"
	"testing"
)

var (
	fooFlag         Array
	fooFlagDuration ArrayDuration
	fooFlagBool     ArrayBool
	fooFlagInt      ArrayInt
)

func init() {
	os.Args = append(os.Args, "--fooFlag=foo", "--fooFlag=bar")
	os.Args = append(os.Args, "--fooFlagDuration=10s", "--fooFlagDuration=5m")
	os.Args = append(os.Args, "--fooFlagBool=true", "--fooFlagBool=false,true")
	os.Args = append(os.Args, "--fooFlagInt=1", "--fooFlagInt=2,3")
	flag.Var(&fooFlag, "fooFlag", "test")
	flag.Var(&fooFlagDuration, "fooFlagDuration", "test")
	flag.Var(&fooFlagBool, "fooFlagBool", "test")
	flag.Var(&fooFlagInt, "fooFlagInt", "test")
}

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestArray(t *testing.T) {
	expected := []string{
		"foo",
		"bar",
	}
	if len(expected) != len(fooFlag) {
		t.Errorf("len array flag (%d) is not equal to %d", len(fooFlag), len(expected))
	}
	for i, v := range fooFlag {
		if v != expected[i] {
			t.Errorf("unexpected item in array %q", v)
		}
	}
}

func TestArraySet(t *testing.T) {
	f := func(s string, expectedValues []string) {
		t.Helper()
		var a Array
		_ = a.Set(s)
		if !reflect.DeepEqual([]string(a), expectedValues) {
			t.Fatalf("unexpected values parsed;\ngot\n%q\nwant\n%q", a, expectedValues)
		}
	}
	f("", nil)
	f(`foo`, []string{`foo`})
	f(`foo,b ar,baz`, []string{`foo`, `b ar`, `baz`})
	f(`foo,b\"'ar,"baz,d`, []string{`foo`, `b\"'ar`, `"baz,d`})
	f(`,foo,,ba"r,`, []string{``, `foo`, ``, `ba"r`, ``})
	f(`""`, []string{``})
	f(`"foo,b\nar"`, []string{`foo,b` + "\n" + `ar`})
	f(`"foo","bar",baz`, []string{`foo`, `bar`, `baz`})
	f(`,fo,"\"b, a'\\",,r,`, []string{``, `fo`, `"b, a'\`, ``, `r`, ``})
}

func TestArrayString(t *testing.T) {
	f := func(s string) {
		t.Helper()
		var a Array
		_ = a.Set(s)
		result := a.String()
		if result != s {
			t.Fatalf("unexpected string;\ngot\n%s\nwant\n%s", result, s)
		}
	}
	f("")
	f("foo")
	f("foo,bar")
	f(",")
	f(",foo,")
	f(`", foo","b\"ar",`)
	f(`,"\nfoo\\",bar`)
}

func TestArrayDuration(t *testing.T) {
	expected := []Duration{
		{Msecs: int64(10 * msecsPerSecond)},
		{Msecs: int64(5 * msecsPerMinute)},
	}
	if len(expected) != len(fooFlagDuration) {
		t.Errorf("len array flag (%d) is not equal to %d", len(fooFlag), len(expected))
	}
	for i, v := range fooFlagDuration {
		if v != expected[i] {
			t.Errorf("unexpected item in array %v", v)
		}
	}
}

func TestArrayDurationSet(t *testing.T) {
	f := func(s string, expectedValues []Duration) {
		t.Helper()
		var a ArrayDuration
		_ = a.Set(s)
		if !reflect.DeepEqual([]Duration(a), expectedValues) {
			t.Fatalf("unexpected values parsed;\ngot\n%q\nwant\n%q", a, expectedValues)
		}
	}
	f("", nil)
	f(`1m`, []Duration{{Msecs: int64(msecsPerMinute)}})
	f(`5m,1s,1h`, []Duration{{Msecs: int64(5 * msecsPerMinute)}, {Msecs: int64(msecsPerSecond)}, {Msecs: int64(msecsPerHour)}})
}

func TestArrayDurationString(t *testing.T) {
	f := func(s string) {
		t.Helper()
		var a ArrayDuration
		_ = a.Set(s)
		result := a.String()
		if result != s {
			t.Fatalf("unexpected string;\ngot\n%s\nwant\n%s", result, s)
		}
	}
	f("")
	f("10s,1m0s")
	f("5m0s,1s")
}

func TestArrayBool(t *testing.T) {
	expected := []bool{
		true, false, true, true,
	}
	if len(expected) != len(fooFlagBool) {
		t.Errorf("len array flag (%d) is not equal to %d", len(fooFlag), len(expected))
	}
	for i, v := range fooFlagBool {
		if v != expected[i] {
			t.Errorf("unexpected item in array index=%v,value=%v,want=%v", i, v, expected[i])
		}
	}
}

func TestArrayBoolSet(t *testing.T) {
	f := func(s string, expectedValues []bool) {
		t.Helper()
		var a ArrayBool
		_ = a.Set(s)
		if !reflect.DeepEqual([]bool(a), expectedValues) {
			t.Fatalf("unexpected values parsed;\ngot\n%v\nwant\n%v", a, expectedValues)
		}
	}
	f("", nil)
	f(`true`, []bool{true})
	f(`false,True,False`, []bool{false, true, false})
}

func TestArrayBoolString(t *testing.T) {
	f := func(s string) {
		t.Helper()
		var a ArrayBool
		_ = a.Set(s)
		result := a.String()
		if result != s {
			t.Fatalf("unexpected string;\ngot\n%s\nwant\n%s", result, s)
		}
	}
	f("")
	f("true")
	f("true,false")
	f("false,true")
}

func TestArrayInt(t *testing.T) {
	expected := []int{1, 2, 3}
	if len(expected) != len(fooFlagInt) {
		t.Errorf("len array flag (%d) is not equal to %d", len(fooFlag), len(expected))
	}
	for i, n := range fooFlagInt {
		if n != expected[i] {
			t.Errorf("unexpected item in array %d", n)
		}
	}
}

func TestArrayIntSet(t *testing.T) {
	f := func(s string, expectedValues []int) {
		t.Helper()
		var a ArrayInt
		_ = a.Set(s)
		if !reflect.DeepEqual([]int(a), expectedValues) {
			t.Fatalf("unexpected values parsed;\ngot\n%q\nwant\n%q", a, expectedValues)
		}
	}
	f("", nil)
	f(`1`, []int{1})
	f(`-2,3,-64`, []int{-2, 3, -64})
}

func TestArrayIntString(t *testing.T) {
	f := func(s string) {
		t.Helper()
		var a ArrayInt
		_ = a.Set(s)
		result := a.String()
		if result != s {
			t.Fatalf("unexpected string;\ngot\n%s\nwant\n%s", result, s)
		}
	}
	f("")
	f("10,1")
	f("-5,1,123")
}
