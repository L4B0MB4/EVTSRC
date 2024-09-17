package helper_test

import (
	"fmt"
	"math"
	"testing"
	"time"

	"github.com/L4B0MB4/EVTSRC/pkg/helper"
)

func TestInt62Split64Bit(t *testing.T) {
	var b1 int64 = int64(math.Pow(2, 64)) - 1

	v := time.Now().UnixMicro()
	fmt.Print(v)
	_, _, err := helper.SplitInt62(b1)
	if err == nil {
		t.Error("Should have thrown an error with 64 bit integer")
		t.Fail()
	}
}
func TestInt62Split63Bit(t *testing.T) {
	var b1 int64 = int64(math.Pow(2, 63)) - 1

	_, _, err := helper.SplitInt62(b1)
	if err == nil {
		t.Error("Should have thrown an error with 63 bit integer")
		t.Fail()
	}
}

func TestInt62Split62Bit(t *testing.T) {

	var b1 int64 = int64(math.Pow(2, 62)) - 1

	_, _, err := helper.SplitInt62(b1)
	if err != nil {
		t.Error("Should be able to handle 62 bit integer")
		t.Fail()
	}
}

func TestInt62Split(t *testing.T) {

	var b1 int64 = int64(math.Pow(2, 62)) - 1

	i1, i2, err := helper.SplitInt62(b1)
	if err != nil {
		t.Error("Should be able to handle 62 bit integer")
		t.Fail()
	}
	b_created, err := helper.MergeInt62(i1, i2)

	if err != nil {
		t.Error("Should be able to handle two normal integers")
		t.Fail()
	}
	if b_created != b1 {
		t.Error("big integers should have the same values")
		t.Fail()
	}
}
func TestInt62MergeTwoNegatives(t *testing.T) {

	_, err := helper.MergeInt62(int32(math.Pow(2, 32)-1), int32(math.Pow(2, 32)-1))
	if err == nil {
		t.Error("Should have caused an error")
		t.Fail()
	}
}

func TestInt62MergeOneNegativeOnePositive(t *testing.T) {

	_, err := helper.MergeInt62(int32(math.Pow(2, 32)-1), int32(math.Pow(2, 1)))
	if err == nil {
		t.Error("Should have caused an error")
		t.Fail()
	}
}

func TestInt62MergeOnePositiveOneNegative(t *testing.T) {

	_, err := helper.MergeInt62(int32(math.Pow(2, 1)), int32(math.Pow(2, 32)-1))
	if err == nil {
		t.Error("Should have caused an error")
		t.Fail()
	}
}

func TestInt62MergeTwoPositive(t *testing.T) {

	big, err := helper.MergeInt62(int32(math.Pow(2, 31)-1), int32(math.Pow(2, 31)-1))
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	if big != 0x3FFF_FFFF_FFFF_FFFF {
		t.Errorf("Should have been equal to %v but was %v", 0x3FFF_FFFF_FFFF_FFFF, big)
		t.Fail()
	}
}

func TestInt62MergeTwoPositiveEx2(t *testing.T) {

	big, err := helper.MergeInt62(int32(math.Pow(2, 30)-1), int32(math.Pow(2, 31)-1))
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	if big != 0x1FFF_FFFF_FFFF_FFFF {
		t.Errorf("Should have been equal to %v but was %v", 0x1FFF_FFFF_FFFF_FFFF, big)
		t.Fail()
	}
}

func TestInt62MergeTwoPositiveEx3(t *testing.T) {

	big, err := helper.MergeInt62(int32(math.Pow(2, 30)-1), int32(math.Pow(2, 30)-1))
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	if big != 0x1FFF_FFFF_BFFF_FFFF {
		t.Errorf("Should have been equal to %v but was %v", 0x1FFF_FFFF_BFFF_FFFF, big)
		t.Fail()
	}
}

func TestInt62SplitAndMergeEx1(t *testing.T) {

	var b1 int64 = int64(math.Pow(2, 57)) - 1

	i1, i2, err := helper.SplitInt62(b1)
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	big, err := helper.MergeInt62(i1, i2)
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	if big != b1 {
		t.Errorf("Should have been equal to %v but was %v", 0x1FFF_FFFF_BFFF_FFFF, big)
		t.Fail()
	}
}

func TestInt62SplitAndMergeEx2(t *testing.T) {

	var b1 int64 = int64(math.Pow(2, 57)) - 2

	i1, i2, err := helper.SplitInt62(b1)
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	big, err := helper.MergeInt62(i1, i2)
	if err != nil {
		t.Error("Should have been handled without error")
		t.Fail()
	}
	if big != b1 {
		t.Errorf("Should have been equal to %v but was %v", 0x1FFF_FFFF_BFFF_FFFF, big)
		t.Fail()
	}
}
