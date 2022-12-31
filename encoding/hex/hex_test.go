package hex

import (
	"bytes"
	"testing"

	goHex "encoding/hex"
)

type hexBytesPair struct {
	hex   string
	bytes []byte
}

var hexBytesPairs = []hexBytesPair{
	{"", []byte{}},
	{"00", []byte{0}},
	{"01", []byte{1}},
	{"02", []byte{2}},
	{"03", []byte{3}},
	{"04", []byte{4}},
	{"05", []byte{5}},
	{"06", []byte{6}},
	{"07", []byte{7}},
	{"08", []byte{8}},
	{"09", []byte{9}},
	{"0A", []byte{10}},
	{"0B", []byte{11}},
	{"0C", []byte{12}},
	{"0D", []byte{13}},
	{"0E", []byte{14}},
	{"0F", []byte{15}},
	{"10", []byte{16}},
	{"1F", []byte{31}},
	{"20", []byte{32}},
	{"3F", []byte{63}},
	{"40", []byte{64}},
	{"7F", []byte{127}},
	{"80", []byte{128}},
	{"FF", []byte{255}},

	{"0000", []byte{0, 0}},
	{"00000000", []byte{0, 0, 0, 0}},                     //int32(0)
	{"0000000000000000", []byte{0, 0, 0, 0, 0, 0, 0, 0}}, //int64(0)
	{"0001020304050607", []byte{0, 1, 2, 3, 4, 5, 6, 7}},
	{"08090A0B0C0D0E0F", []byte{8, 9, 10, 11, 12, 13, 14, 15}},
	{"F0F1F2F3F4F5F6F7", []byte{0xf0, 0xf1, 0xf2, 0xf3, 0xf4, 0xf5, 0xf6, 0xf7}},
	{"F8F9FAFBFCFDFEFF", []byte{0xf8, 0xf9, 0xfa, 0xfb, 0xfc, 0xfd, 0xfe, 0xff}},

	//Everyone's favourite hex-word
	{"DEADBEEF", []byte{0xde, 0xad, 0xbe, 0xef}},

	//32 bit number
	{"B226F0C8",
		[]byte{178, 38, 240, 200}},
	//64 bit number
	{"9AB60B0C2E1409E7",
		[]byte{154, 182, 11, 12, 46, 20, 9, 231}},
	//128 bit number
	{"492DD7AB258F61437CC76070F5FE2759",
		[]byte{73, 45, 215, 171, 37, 143, 97, 67, 124, 199, 96, 112, 245, 254, 39, 89}},
	//256 bit number
	{"1BE5D27A94CF25273024ABD8649C196BA5141EA309D2A579AC96A9C4FC5C4179",
		[]byte{27, 229, 210, 122, 148, 207, 37, 39, 48, 36, 171, 216, 100, 156, 25, 107, 165, 20, 30, 163, 9, 210, 165, 121, 172, 150, 169, 196, 252, 92, 65, 121}},
	//512 bit number
	{"18B32CDAF3C99B992BA4A72FDF57817C4BF1395276F5B90625D0CB66A987AC5F91A0DEEC6F45F27509A8AC4E88DEC209EA86D36AD9F109BD0084E2DD7F30B88C",
		[]byte{24, 179, 44, 218, 243, 201, 155, 153, 43, 164, 167, 47, 223, 87, 129, 124, 75, 241, 57, 82, 118, 245, 185, 6, 37, 208, 203, 102, 169, 135, 172, 95, 145, 160, 222, 236, 111, 69, 242, 117, 9, 168, 172, 78, 136, 222, 194, 9, 234, 134, 211, 106, 217, 241, 9, 189, 0, 132, 226, 221, 127, 48, 184, 140}},
}

func TestToBytes(t *testing.T) {
	for _, rec := range hexBytesPairs {
		found, err := ToBytes(rec.hex)

		if err != nil {
			t.Errorf("Got an error decoding %v: %s", rec.hex, err)
		}
		if !bytes.Equal(found, rec.bytes) {
			t.Errorf("%v, expecting, got\n %v\n %v", rec.hex, rec.bytes, found)
		}
	}
}

func TestFromBytes(t *testing.T) {
	for _, rec := range hexBytesPairs {
		found := FromBytes(rec.bytes)

		if found != rec.hex {
			t.Errorf("Expecting, got\n %v\n %v", rec.hex, found)
		}
	}
}

func TestLowerCaseDecodes(t *testing.T) {
	str := "b226f0c8"
	found, err := ToBytes(str)
	if err != nil {
		t.Errorf("Got an error decoding %v: %s", str, err)
	}
	if !bytes.Equal(found, []byte{178, 38, 240, 200}) {
		t.Errorf("Got invalid byte-value for %s\n%v", str, found)
	}
}

var badHex = []struct {
	hex string
	err error
}{
	{"A", ErrLength},
	{"AaA", ErrLength},
	{"0Z", InvalidHexAt("Z"[0], 1)},
	{"Z0", InvalidHexAt("Z"[0], 0)},
	{"1.1", InvalidHexAt("."[0], 1)},
	{" A", InvalidHexAt(" "[0], 0)},
	//Col 0 chars
	{"\x00", InvalidHexAt("\x00"[0], 0)},
	{"\x1F", InvalidHexAt("\x1F"[0], 0)},
	//Col 1
	{" ", InvalidHexAt(" "[0], 0)}, //Before valid
	{"/", InvalidHexAt("/"[0], 0)}, //Before valid
	{":", InvalidHexAt(":"[0], 0)}, //After valid
	{"?", InvalidHexAt("?"[0], 0)}, //After valid
	//Col 2
	{"@", InvalidHexAt("@"[0], 0)}, //Before valid
	{"G", InvalidHexAt("G"[0], 0)}, //After valid
	{"Z", InvalidHexAt("Z"[0], 0)}, //After valid
	{"_", InvalidHexAt("_"[0], 0)}, //After valid
	//Col 3
	{"`", InvalidHexAt("`"[0], 0)}, //Before valid
	{"g", InvalidHexAt("g"[0], 0)}, //After valid
	{"z", InvalidHexAt("z"[0], 0)}, //After valid
	{"~", InvalidHexAt("~"[0], 0)}, //After valid
	//Col 8
	{"\xFF", InvalidHexAt("\xFF"[0], 0)}, //Invalid
}

func TestToBytesInvalidHex(t *testing.T) {
	for _, rec := range badHex {
		found, err := ToBytes(rec.hex)
		if err == nil {
			t.Errorf("No error decoding %s: %v", rec.hex, found)
		} else if err != rec.err {
			t.Errorf("Decoding %s, wrong error: %s", rec.hex, err)
		}
	}
}

const invalid = byte(20)

var decodePairs = []struct {
	char   byte
	expect byte
}{
	//Col 0 (ctrl)
	{'\x00', invalid},
	{'\x01', invalid},
	{'\x02', invalid},
	{'\x03', invalid},
	{'\x04', invalid},
	{'\x05', invalid},
	{'\x06', invalid},
	{'\x07', invalid},
	{'\x08', invalid},
	{'\x09', invalid},
	{'\x0A', invalid},
	{'\x0B', invalid},
	{'\x0C', invalid},
	{'\x0D', invalid},
	{'\x0E', invalid},
	{'\x0F', invalid},
	//Col 1 (ctrl)
	{'\x10', invalid},
	{'\x11', invalid},
	{'\x12', invalid},
	{'\x13', invalid},
	{'\x14', invalid},
	{'\x15', invalid},
	{'\x16', invalid},
	{'\x17', invalid},
	{'\x18', invalid},
	{'\x19', invalid},
	{'\x1A', invalid},
	{'\x1B', invalid},
	{'\x1C', invalid},
	{'\x1D', invalid},
	{'\x1E', invalid},
	{'\x1F', invalid},
	//Col 2 (sp-/)
	{' ', invalid},
	{'!', invalid},
	{'"', invalid},
	{'#', invalid},
	{'$', invalid},
	{'%', invalid},
	{'&', invalid},
	{'\'', invalid},
	{'(', invalid},
	{')', invalid},
	{'*', invalid},
	{'+', invalid},
	{',', invalid},
	{'-', invalid},
	{'.', invalid},
	{'/', invalid},
	//Col 3 (0-?)
	{'0', 0},
	{'1', 1},
	{'2', 2},
	{'3', 3},
	{'4', 4},
	{'5', 5},
	{'6', 6},
	{'7', 7},
	{'8', 8},
	{'9', 9},
	{':', invalid},
	{';', invalid},
	{'<', invalid},
	{'=', invalid},
	{'>', invalid},
	{'?', invalid},
	//Col 4 (@-O)
	{'@', invalid},
	{'A', 10},
	{'B', 11},
	{'C', 12},
	{'D', 13},
	{'E', 14},
	{'F', 15},
	{'G', invalid},
	{'H', invalid},
	{'I', invalid},
	{'J', invalid},
	{'K', invalid},
	{'L', invalid},
	{'M', invalid},
	{'N', invalid},
	{'O', invalid},
	//Col 5 (P-_)
	{'P', invalid},
	{'Q', invalid},
	{'R', invalid},
	{'S', invalid},
	{'T', invalid},
	{'U', invalid},
	{'V', invalid},
	{'W', invalid},
	{'X', invalid},
	{'Y', invalid},
	{'Z', invalid},
	{'[', invalid},
	{'\\', invalid},
	{']', invalid},
	{'^', invalid},
	{'_', invalid},
	//Col 6 (`-0)
	{'`', invalid},
	{'a', 10},
	{'b', 11},
	{'c', 12},
	{'d', 13},
	{'e', 14},
	{'f', 15},
	{'g', invalid},
	{'h', invalid},
	{'i', invalid},
	{'j', invalid},
	{'k', invalid},
	{'l', invalid},
	{'m', invalid},
	{'n', invalid},
	{'o', invalid},
	//Col 7 (p-DEL)
	{'p', invalid},
	{'q', invalid},
	{'r', invalid},
	{'s', invalid},
	{'t', invalid},
	{'u', invalid},
	{'v', invalid},
	{'w', invalid},
	{'x', invalid},
	{'y', invalid},
	{'z', invalid},
	{'{', invalid},
	{'|', invalid},
	{'}', invalid},
	{'~', invalid},
	{'\x7f', invalid},

	//Col 8-F (invalid)
	{'\x80', invalid},
	{'\x81', invalid},
	{'\x82', invalid},
	{'\x83', invalid},
	{'\x84', invalid},
	{'\x85', invalid},
	{'\x86', invalid},
	{'\x87', invalid},
	{'\x88', invalid},
	{'\x89', invalid},
	{'\x8A', invalid},
	{'\x8B', invalid},
	{'\x8C', invalid},
	{'\x8D', invalid},
	{'\x8E', invalid},
	{'\x8F', invalid},

	{'\x90', invalid},
	{'\x91', invalid},
	{'\x92', invalid},
	{'\x93', invalid},
	{'\x94', invalid},
	{'\x95', invalid},
	{'\x96', invalid},
	{'\x97', invalid},
	{'\x98', invalid},
	{'\x99', invalid},
	{'\x9A', invalid},
	{'\x9B', invalid},
	{'\x9C', invalid},
	{'\x9D', invalid},
	{'\x9E', invalid},
	{'\x9F', invalid},

	{'\xA0', invalid},
	{'\xA1', invalid},
	{'\xA2', invalid},
	{'\xA3', invalid},
	{'\xA4', invalid},
	{'\xA5', invalid},
	{'\xA6', invalid},
	{'\xA7', invalid},
	{'\xA8', invalid},
	{'\xA9', invalid},
	{'\xAA', invalid},
	{'\xAB', invalid},
	{'\xAC', invalid},
	{'\xAD', invalid},
	{'\xAE', invalid},
	{'\xAF', invalid},

	{'\xB0', invalid},
	{'\xB1', invalid},
	{'\xB2', invalid},
	{'\xB3', invalid},
	{'\xB4', invalid},
	{'\xB5', invalid},
	{'\xB6', invalid},
	{'\xB7', invalid},
	{'\xB8', invalid},
	{'\xB9', invalid},
	{'\xBA', invalid},
	{'\xBB', invalid},
	{'\xBC', invalid},
	{'\xBD', invalid},
	{'\xBE', invalid},
	{'\xBF', invalid},

	{'\xC0', invalid},
	{'\xC1', invalid},
	{'\xC2', invalid},
	{'\xC3', invalid},
	{'\xC4', invalid},
	{'\xC5', invalid},
	{'\xC6', invalid},
	{'\xC7', invalid},
	{'\xC8', invalid},
	{'\xC9', invalid},
	{'\xCA', invalid},
	{'\xCB', invalid},
	{'\xCC', invalid},
	{'\xCD', invalid},
	{'\xCE', invalid},
	{'\xCF', invalid},

	{'\xD0', invalid},
	{'\xD1', invalid},
	{'\xD2', invalid},
	{'\xD3', invalid},
	{'\xD4', invalid},
	{'\xD5', invalid},
	{'\xD6', invalid},
	{'\xD7', invalid},
	{'\xD8', invalid},
	{'\xD9', invalid},
	{'\xDA', invalid},
	{'\xDB', invalid},
	{'\xDC', invalid},
	{'\xDD', invalid},
	{'\xDE', invalid},
	{'\xDF', invalid},

	{'\xE0', invalid},
	{'\xE1', invalid},
	{'\xE2', invalid},
	{'\xE3', invalid},
	{'\xE4', invalid},
	{'\xE5', invalid},
	{'\xE6', invalid},
	{'\xE7', invalid},
	{'\xE8', invalid},
	{'\xE9', invalid},
	{'\xEA', invalid},
	{'\xEB', invalid},
	{'\xEC', invalid},
	{'\xED', invalid},
	{'\xEE', invalid},
	{'\xEF', invalid},

	{'\xF0', invalid},
	{'\xF1', invalid},
	{'\xF2', invalid},
	{'\xF3', invalid},
	{'\xF4', invalid},
	{'\xF5', invalid},
	{'\xF6', invalid},
	{'\xF7', invalid},
	{'\xF8', invalid},
	{'\xF9', invalid},
	{'\xFA', invalid},
	{'\xFB', invalid},
	{'\xFC', invalid},
	{'\xFD', invalid},
	{'\xFE', invalid},
	{'\xFF', invalid},
}

func TestDecode(t *testing.T) {
	for _, rec := range decodePairs {
		found := decode(rec.char)
		if rec.expect > 15 {
			//Found should also be >15
			if found <= 15 {
				t.Errorf("Decoding %s, %d expecting invalid got %d", string(rec.char), rec.char, found)
			}
		} else {
			if found != rec.expect {
				t.Errorf("Decoding %s, expecting %d got %d", string(rec.char), rec.expect, found)
			}
		}
	}
}

var bench = []string{
	"",
	"B226F0C8",
	"9AB60B0C2e1409E7",
	"492DD7AB258F61437CC76070F5fE2759",
	"1BE5D27A94cf25273024ABD8649C196BA5141EA309D2A579AC96A9C4fc5C4179",
	"18B32CDAF3C99B992BA4A72FDF57817C4BF1395276F5B90625D0CB66A987AC5F91A0DEEC6F45F27509A8AC4E88DEC209EA86D36AD9F109BD0084E2DD7F30B88C",
	"fbe2e5f0eee3c820fbeafaebef20fffbf0e1e0f0f520e0ed20e8ece0ebe5f0f2f120fff0eeec20f120faf2fee5e2202ce8f6f3ede220e8e6eee1e8f0f2d1202ce8f0f2e5e220e5d1",
	"fbe2e5f0eee3c820fbeafaebef20fffbf0e1e0f0f520e0ed20e8ece0ebe5f0f2f120fff0eeec20f120faf2fee5e2202ce8f6f3ede220e8e6eee1e8f0f2d1202ce8f0f2e5e220e5d1" +
		"fbe2e5f0eee3c820fbeafaebef20fffbf0e1e0f0f520e0ed20e8ece0ebe5f0f2f120fff0eeec20f120faf2fee5e2202ce8f6f3ede220e8e6eee1e8f0f2d1202ce8f0f2e5e220e5d1" +
		"fbe2e5f0eee3c820fbeafaebef20fffbf0e1e0f0f520e0ed20e8ece0ebe5f0f2f120fff0eeec20f120faf2fee5e2202ce8f6f3ede220e8e6eee1e8f0f2d1202ce8f0f2e5e220e5d1",
	"I'm invalid hex",
	"492DD7AB258F61437CC76070F5fE2759 Invalid",
}

// While we don't need to track error count, we need the benchmarks to have the same
// penalty when the input is invalid.  And we need the compiler not to optimize out the call
func BenchmarkGoNative(b *testing.B) {
	e := 0
	for _, v := range bench {
		for i := 0; i < b.N; i++ {
			_, err := goHex.DecodeString(v)
			if err != nil {
				e++
			}
		}
	}
}

func BenchmarkHex(b *testing.B) {
	e := 0
	for _, v := range bench {
		for i := 0; i < b.N; i++ {
			_, err := ToBytes(v)
			if err != nil {
				e++
			}
		}
	}
}
