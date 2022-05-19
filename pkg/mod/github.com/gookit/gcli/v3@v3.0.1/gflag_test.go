package gcli_test

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"testing"

	"github.com/gookit/gcli/v3"
	"github.com/gookit/goutil/dump"
	"github.com/stretchr/testify/assert"
)

func TestFlags_Basic(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	assert.Len(t, fs.Metas(), 0)
	assert.Equal(t, 0, fs.Len())
	assert.Equal(t, "testFlags", fs.Name())

	assert.Nil(t, fs.LookupFlag("opt1"))
	assert.Len(t, fs.ShortNames("opt"), 0)
	assert.False(t, fs.HasFlag("opt1"))
	assert.False(t, fs.HasOption("opt1"))

	var s1, s2 string
	fs.StrOpt(&s1, "str1", "", "", "desc")

	assert.True(t, fs.IsOption("str1"))
	assert.False(t, fs.IsShortOpt("str1"))
	assert.False(t, fs.IsShortName("str1"))
	assert.False(t, fs.IsOption("not-exist"))

	fs.StrOpt(&s2, "str2", "b", "", "desc")
	assert.True(t, fs.IsShortName("b"))

	buf := new(bytes.Buffer)

	fs.IterAll(func(f *flag.Flag, meta *gcli.FlagMeta) {
		_, _ = fmt.Fprintf(buf, "flag: %s, shorts: %s;", f.Name, meta.Shorts2String(","))
	})

	assert.Equal(t, "flag: str1, shorts: ;flag: str2, shorts: b;", buf.String())
}

func TestFlags_BoolOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var b1, b2 bool
	b0 := fs.Bool("bl0", "", false, "desc0")
	fs.BoolOpt(&b1, "bl1", "ab", false, "desc1")
	fs.BoolVar(&b2, &gcli.FlagMeta{
		Name: "bl2",
		Desc: "desc2",
	})

	assert.False(t, *b0)
	assert.False(t, b1)
	assert.NoError(t, fs.Parse([]string{"--bl0", "-a", "--bl2"}))
	assert.True(t, *b0)
	assert.True(t, b1)
}

func TestFlags_StrOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")
	assert.Len(t, fs.Metas(), 0)

	var str string
	fs.StrVar(&str, &gcli.FlagMeta{
		Name: "test",
		Desc: "test desc",
	})

	assert.True(t, fs.HasFlagMeta("test"))
	assert.False(t, fs.HasFlagMeta("not-exist"))
	assert.Len(t, fs.Metas(), 1)

	f := fs.LookupFlag("test")
	assert.NotEmpty(t, f)

	assert.Equal(t, "test", f.Name)
	assert.Equal(t, "test desc", f.Usage)

	ns := fs.FlagNames()
	assert.Len(t, ns, 1)

	f = fs.LookupFlag("not-exist")
	assert.Nil(t, f)

	err := fs.Parse([]string{})
	assert.NoError(t, err)
	assert.Equal(t, "", str)

	err = fs.Parse([]string{"--test", "value"})
	assert.NoError(t, err)
	assert.Equal(t, "value", str)
	assert.Len(t, fs.ShortNames("test"), 0)
}

func TestFlags_Float64Opt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var f1, f2 float64
	fs.Float64Opt(&f1, "f1", "ab", 0, "desc1")
	fs.Float64Var(&f2, &gcli.FlagMeta{
		Name:   "f2",
		Desc:   "desc2",
		DefVal: 3.14,
	})

	assert.Equal(t, float64(0), f1)
	assert.Equal(t, 3.14, f2)
	assert.NoError(t, fs.Parse([]string{"-a", "12.3", "--f2", "1.63"}))
	assert.Equal(t, 12.3, f1)
	assert.Equal(t, 1.63, f2)
}

func TestFlags_IntOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var int1, int2 int
	fs.IntOpt(&int1, "int1", "ab", 0, "desc1")
	fs.IntVar(&int2, &gcli.FlagMeta{
		Name:   "int2",
		Desc:   "desc2",
		DefVal: 314,
	})

	assert.Equal(t, 0, int1)
	assert.Equal(t, 314, int2)
	assert.NoError(t, fs.Parse([]string{"-a", "123", "--int2", "163"}))
	assert.Equal(t, 123, int1)
	assert.Equal(t, 163, int2)
}

func TestFlags_Int64Opt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var int1, int2 int64
	fs.Int64Opt(&int1, "int1", "ab", 0, "desc1")
	fs.Int64Var(&int2, &gcli.FlagMeta{
		Name:   "int2",
		Desc:   "desc2",
		DefVal: 314,
	})

	assert.Equal(t, int64(0), int1)
	assert.Equal(t, int64(314), int2)
	assert.NoError(t, fs.Parse([]string{"-a", "12", "--int2", "16"}))
	assert.Equal(t, int64(12), int1)
	assert.Equal(t, int64(16), int2)
}

func TestFlags_UintOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var int1, int2 uint
	fs.UintOpt(&int1, "int1", "ab", 0, "desc1")
	fs.UintVar(&int2, &gcli.FlagMeta{
		Name:   "c",
		Desc:   "desc2",
		DefVal: 314,
	})

	assert.Equal(t, uint(0), int1)
	assert.Equal(t, uint(314), int2)
	assert.NoError(t, fs.Parse([]string{"-a", "12", "-c", "16"}))
	assert.Equal(t, uint(12), int1)
	assert.Equal(t, uint(16), int2)
}

func TestFlags_Uint64Opt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var uint1, uint2 uint64
	fs.Uint64Opt(&uint1, "uint1", "ab", 0, "desc1")
	fs.Uint64Var(&uint2, &gcli.FlagMeta{
		Name:   "uint2",
		Desc:   "desc2",
		DefVal: 314,
		Shorts: []string{"c", "", "f"},
	})

	fm2 := fs.FlagMeta("uint2")
	assert.Len(t, fm2.Shorts, 2)
	assert.Equal(t, "c,f", fm2.Shorts2String())

	assert.Equal(t, uint64(0), uint1)
	assert.Equal(t, uint64(314), uint2)
	assert.NoError(t, fs.Parse([]string{"-a", "12", "--uint2", "16"}))
	assert.Equal(t, uint64(12), uint1)
	assert.Equal(t, uint64(16), uint2)
}

func TestFlags_VarOpt(t *testing.T) {
	fs := gcli.NewFlags("testFlags")

	var ints gcli.Ints
	fs.Var(&ints, &gcli.FlagMeta{Name: "ints", Desc: "desc"})
	assert.NoError(t, fs.Parse([]string{"--ints", "12", "--ints", "16"}))

	assert.Len(t, ints, 2)
	assert.Equal(t, "[12 16]", ints.String())

	var ss gcli.Strings
	fs.VarOpt(&ss, "names", "ns", "desc")
	assert.NoError(t, fs.Parse([]string{"--names", "abc", "-n", "def", "-s", "ghi"}))

	assert.Len(t, ss, 3)
	assert.Equal(t, "[abc def ghi]", ss.String())
}

func TestFlags_CheckName(t *testing.T) {
	assert.PanicsWithValue(t, "GCli: redefined option flag 'opt1'", func() {
		var i int64
		fs := gcli.NewFlags()
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
		fs.Int64Opt(&i, "opt1", "", 0, "desc")
	})

	assert.PanicsWithValue(t, "GCli: option flag name cannot be empty", func() {
		var b bool
		fs := gcli.NewFlags()
		fs.BoolOpt(&b, "", "", false, "desc")
	})

	assert.PanicsWithValue(t, "GCli: option flag name '+invalid' is invalid, must match: ^[a-zA-Z][\\w-]*$", func() {
		var fv uint
		fs := gcli.NewFlags()
		fs.UintOpt(&fv, "+invalid", "", 0, "desc")
	})

	assert.PanicsWithValue(t, "GCli: option flag name '78' is invalid, must match: ^[a-zA-Z][\\w-]*$", func() {
		var fv uint64
		fs := gcli.NewFlags()
		fs.Uint64Opt(&fv, "78", "", 0, "desc")
	})
}

func TestFlags_CheckShorts(t *testing.T) {
	assert.NotPanics(t, func() {
		var fv float64
		fs := gcli.NewFlags()

		// "+" has been filtered by func: splitShortcut()
		fs.Float64Opt(&fv, "float", "+", 0, "desc")

		fm := fs.FlagMeta("float")
		assert.Len(t, fm.Shorts, 0)
	})

	assert.PanicsWithValue(t, "GCli: short name only allow: A-Za-z given: '+'", func() {
		var fv float64
		fs := gcli.NewFlags()
		fs.Float64Var(&fv, &gcli.FlagMeta{
			Name:   "float",
			Shorts: []string{"+"},
		})
	})

	assert.PanicsWithValue(t, "GCli: short name 'a' has been used as the current option name", func() {
		var i int
		fs := gcli.NewFlags()
		fs.IntVar(&i, &gcli.FlagMeta{Name: "a", Shorts: []string{"a"}})
	})

	assert.PanicsWithValue(t, "GCli: short name 's' has been used as an option name", func() {
		var i int
		fs := gcli.NewFlags()
		fs.IntOpt(&i, "s", "", 0, "desc")
		fs.IntOpt(&i, "int1", "s", 0, "desc")
	})

	assert.PanicsWithValue(t, "GCli: short name 's' has been used by option 'str'", func() {
		var str string
		fs := gcli.NewFlags()
		fs.StrOpt(&str, "str", "s", "", "desc")
		fs.StrOpt(&str, "str1", "s", "", "desc")
	})
}

var flagOpts = struct {
	intv int
	strv string
}{}

func TestFlags_Run(t *testing.T) {
	is := assert.New(t)

	fg := gcli.NewFlags("test", "desc message")
	fg.ExitFunc = func(code int) {}

	fg.IntOpt(&flagOpts.intv, "intv", "i", 0, "desc message for intv")
	fg.StrOpt(&flagOpts.strv, "strv", "s", "", "desc message for strv")

	// parse
	fg.Run([]string{"./app", "-i", "23", "-s", "inhere"})
	is.Equal(23, flagOpts.intv)
	is.Equal("inhere", flagOpts.strv)

	// help
	fg.Run([]string{"./app", "-h"})
}

func TestFlags_Parse(t *testing.T) {
	var str string

	gf := gcli.NewFlags("test")
	gf.StrVar(&str, &gcli.FlagMeta{
		Name:     "opt1",
		Required: true,
		Validator: func(val string) error {
			if len(val) < 5 {
				return errors.New("flag value min len is 5")
			}

			return nil
		},
	})

	err := gf.Parse([]string{})
	assert.Error(t, err)

	err = gf.Parse([]string{"--opt1", ""})
	assert.Error(t, err)

	err = gf.Parse([]string{"--opt1", "val"})
	assert.Error(t, err)
	assert.Equal(t, "flag value min len is 5", err.Error())

	err = gf.Parse([]string{"--opt1", "value"})
	assert.NoError(t, err)
	assert.Equal(t, "value", str)
}

func TestFlags_Int_bindingNilPtr(t *testing.T) {
	type userOpts struct {
		Int *int
	}

	opt := userOpts{}
	dump.P(opt)

	// must init for an ptr value.
	assert.Panics(t, func() {
		fs := gcli.NewFlags("test")
		fs.IntOpt(opt.Int, "int", "i", 20, "")
	})

	aint := 23
	opt = userOpts{Int: &aint}
	dump.P(opt)
	fs := gcli.NewFlags("test")
	fs.IntOpt(opt.Int, "int", "i", 20, "")
	fs.PrintHelpPanel()
}

func TestFlags_FromStruct_simple(t *testing.T) {
	fs := gcli.NewFlags("test")

	type userOpts0 struct {
		Int int    `flag:"name=int0;shorts=i;required=true;desc=int option message"`
		Bol bool   `flag:"name=bol;shorts=b;default=true;desc=bool option message"`
		Str string `flag:"name=str1;shorts=o,h;required=true;desc=str1 message"`
	}

	opt := &userOpts0{}
	assert.False(t, opt.Bol)
	assert.Equal(t, 0, opt.Int)
	assert.Equal(t, "", opt.Str)

	err := fs.FromStruct(opt)
	assert.NoError(t, err)
	assert.True(t, opt.Bol)
	assert.True(t, fs.HasFlagMeta("bol"))

	fs.PrintHelpPanel()

	err = fs.Parse([]string{"--int0", "13", "--str1", "xyz"})
	assert.NoError(t, err)

	assert.Equal(t, 13, opt.Int)
	assert.Equal(t, "xyz", opt.Str)
}

func TestFlags_FromStruct_ptrField(t *testing.T) {
	type userOpts struct {
		Int *int    `flag:"shorts=i;default=13;desc=int option message"`
		Str *string `flag:"name=str2;required=true;desc=str2 message"`
	}

	opt := &userOpts{}
	dump.P(opt)

	fs := gcli.NewFlags("test1")
	opt = &userOpts{}
	err := fs.FromStruct(opt)
	assert.Error(t, err)
	assert.Equal(t, "field: Int - nil pointer dereference", err.Error())

	aint := 23
	astr := "xyz"
	opt = &userOpts{Int: &aint, Str: &astr}
	dump.P(opt)
	assert.Equal(t, 23, *opt.Int)
	assert.Equal(t, "xyz", *opt.Str)

	fs = gcli.NewFlags("test1")
	err = fs.FromStruct(opt)
	assert.NoError(t, err)
	assert.Equal(t, 13, *(opt.Int))
	// assert.Equal(t, "", *opt.Str)

	dump.P(opt)
	fmt.Println("Flag Help:")
	fs.PrintHelpPanel()
}

func TestFlags_FromStruct(t *testing.T) {
	type userOpts struct {
		Int  int    `flag:"name=int0;shorts=i;required=true;desc=int option message"`
		Bol  bool   `flag:"name=bol;shorts=b;desc=bool option message"`
		Str1 string `flag:"name=str1;shorts=o,h;required=true;desc=str1 message"`
		// use ptr
		Str2 *string `flag:"name=str2;required=true;desc=str2 message"`
		// custom type and implement flag.Value
		Verb0 gcli.VerbLevel `flag:"name=verb0;shorts=V;desc=verb0 message"`
		// use ptr
		Verb1 *gcli.VerbLevel `flag:"name=verb1;desc=verb1 message"`
	}

	astr := "xyz"
	verb := gcli.VerbWarn
	fs := gcli.NewFlags("test")
	err := fs.FromStruct(&userOpts{
		Str2:  &astr,
		Verb1: &verb,
	})
	assert.NoError(t, err)

	help := fs.String()
	assert.Contains(t, help, "-o, -h, --str1")

	fmt.Println("Flag Help:")
	fs.PrintHelpPanel()
}

// func TestFlags_FromText(t *testing.T) {
// }

func TestFlags_PrintHelpPanel(t *testing.T) {
	fs := gcli.NewFlags("test")

	testOpts := struct {
		opt1 int
		opt2 bool
		opt3 string
	}{}

	fs.IntVar(&testOpts.opt1, &gcli.FlagMeta{Name: "opt1"})
	fs.StrVar(&testOpts.opt3, &gcli.FlagMeta{
		Name: "test",
		Desc: "test desc",
		// required
		Required: true,
	})
	fs.BoolOpt(&testOpts.opt2, "bol", "ab", false, "opt2 desc")
	fs.PrintHelpPanel()
}

func TestFlagMeta_Validate(t *testing.T) {
	fm := gcli.FlagMeta{
		Name:     "opt1",
		Required: true,
		Validator: func(val string) error {
			if len(val) < 5 {
				return errors.New("flag value min len is 5")
			}

			return nil
		},
	}

	err := fm.Validate("")
	assert.Error(t, err)
	assert.Equal(t, "flag 'opt1' is required", err.Error())

	err = fm.Validate("val")
	assert.Error(t, err)
	assert.Equal(t, "flag value min len is 5", err.Error())

	err = fm.Validate("value")
	assert.NoError(t, err)
}
