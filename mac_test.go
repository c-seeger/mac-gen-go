package macgen

import (
	"fmt"
	"net"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestPrefixValues struct {
	Input bool
	Test  bool
}

func TestGenerateRandomLocalPrefix(t *testing.T) {
	TestValues := []TestPrefixValues{
		TestPrefixValues{true, false},
		TestPrefixValues{false, true},
	}

	for _, test := range TestValues {
		prefix := GenerateRandomLocalMacPrefix(test.Input)

		ui, err := strconv.ParseUint(prefix[0:2], 16, 8)
		if err != nil {
			t.Error(err)
		}
		bin := fmt.Sprintf("%08b", ui)

		asciiRep := 48
		if test.Test {
			asciiRep = 49
		}

		if int(bin[7]) != asciiRep {
			t.Errorf("uni/multicast bit not set correctly, got %v expected first bit to be %d", bin, asciiRep-48)
		}

		if int(bin[6]) != 49 {
			t.Errorf("LAA bit not set correctly, got %v expected second bit to be 1", bin)
		}
	}
}

type TestSufixValues struct {
	Input  net.IP
	Output string
	Error  error
}

func TestCalculateNICSufix(t *testing.T) {

	TestValues := []TestSufixValues{
		TestSufixValues{net.ParseIP("10.0.0.0"), "00:00:00", nil},
		TestSufixValues{net.ParseIP("10.255.255.255"), "ff:ff:ff", nil},
		TestSufixValues{net.ParseIP("192.168.12.127"), "04:4a:7f", nil},
		TestSufixValues{net.ParseIP("::1"), "", fmt.Errorf("ip is not v4")},
		TestSufixValues{net.ParseIP("foo"), "", fmt.Errorf("ip is not v4")},
	}

	for _, test := range TestValues {
		mac, err := CalculateNICSufix(test.Input)
		assert.Equal(t, test.Output, mac)
		assert.Equal(t, test.Error, err)
	}
}
