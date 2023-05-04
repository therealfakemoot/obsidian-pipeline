package gloss

import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"text/template"
	"unicode"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

func CharmUsage(c *cobra.Command) error {
	var b bytes.Buffer
	err := tmpl(&b, c.UsageTemplate(), c)
	if err != nil {
		c.PrintErrln(err)
		return err
	}
	pretty(c.ErrOrStderr(), b.String())
	return nil
}

func CharmHelp(c *cobra.Command, a []string) {
	var b bytes.Buffer
	// The help should be sent to stdout
	// See https://github.com/spf13/cobra/issues/1002
	err := tmpl(&b, c.HelpTemplate(), c)
	if err != nil {
		c.PrintErrln(err)
	}
	pretty(c.ErrOrStderr(), b.String())
}

var BaseStyle = lipgloss.NewStyle().Bold(true).Width(80)

func pretty(w io.Writer, s string) {
	fmt.Fprintf(w, "%s\n", BaseStyle.Render(s))
}

var templateFuncs = template.FuncMap{
	"trim":                    strings.TrimSpace,
	"trimRightSpace":          trimRightSpace,
	"trimTrailingWhitespaces": trimRightSpace,
	"appendIfNotPresent":      appendIfNotPresent,
	"rpad":                    rpad,
	"gt":                      Gt,
	"eq":                      Eq,
}

func Gt(a interface{}, b interface{}) bool {
	var left, right int64
	av := reflect.ValueOf(a)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		left = int64(av.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		left = av.Int()
	case reflect.String:
		left, _ = strconv.ParseInt(av.String(), 10, 64)
	}

	bv := reflect.ValueOf(b)

	switch bv.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		right = int64(bv.Len())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		right = bv.Int()
	case reflect.String:
		right, _ = strconv.ParseInt(bv.String(), 10, 64)
	}

	return left > right
}

// FIXME Eq is unused by cobra and should be removed in a version 2. It exists only for compatibility with users of cobra.

// Eq takes two types and checks whether they are equal. Supported types are int and string. Unsupported types will panic.
func Eq(a interface{}, b interface{}) bool {
	av := reflect.ValueOf(a)
	bv := reflect.ValueOf(b)

	switch av.Kind() {
	case reflect.Array, reflect.Chan, reflect.Map, reflect.Slice:
		panic("Eq called on unsupported type")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return av.Int() == bv.Int()
	case reflect.String:
		return av.String() == bv.String()
	}
	return false
}

func trimRightSpace(s string) string {
	return strings.TrimRightFunc(s, unicode.IsSpace)
}

// appendIfNotPresent will append stringToAppend to the end of s, but only if it's not yet present in s.
func appendIfNotPresent(s, stringToAppend string) string {
	if strings.Contains(s, stringToAppend) {
		return s
	}
	return s + " " + stringToAppend
}

// rpad adds padding to the right of a string.
func rpad(s string, padding int) string {
	formattedString := fmt.Sprintf("%%-%ds", padding)
	return fmt.Sprintf(formattedString, s)
}

// tmpl executes the given template text on data, writing the result to w.
func tmpl(w io.Writer, text string, data interface{}) error {
	t := template.New("top")
	t.Funcs(templateFuncs)
	template.Must(t.Parse(text))
	return t.Execute(w, data)
}
