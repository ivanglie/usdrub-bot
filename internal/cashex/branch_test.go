package cashex

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func Test_newBranch(t *testing.T) {
	got := newBranch("b", "a", "s", "c", 12.34, 56.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }())
	want := newBranch("b", "a", "s", "c", 12.34, 56.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }())
	if !reflect.DeepEqual(got, want) {
		t.Errorf("newBranch() = %v, want %v", got, want)
	}
}

func Test_BySellSorter(t *testing.T) {
	got := []branch{
		newBranch("bank", "address", "subway", "currency", 12.34, 58.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 56.75, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 78.56, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 56.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 52.64, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
	}
	sort.Sort(BySellSorter(got))

	want := []branch{
		newBranch("bank", "address", "subway", "currency", 12.34, 52.64, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 56.75, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 56.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 58.78, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
		newBranch("bank", "address", "subway", "currency", 12.34, 78.56, func() time.Time { t, _ := time.Parse("02.01.2006 15:04", "01.02.2018 12:35"); return t }()),
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("branch = %v, want %v", got, want)
	}
}
