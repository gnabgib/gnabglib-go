package bytes

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/gnabgib/gnablib-go/encoding/hex"
	"github.com/gnabgib/gnablib-go/test"
)

var uint16Tests = []struct {
	a, b, sum uint16
	carry     byte
}{
	{0, 0, 0, 0},
	{0, 1, 1, 0},
	{1, 1, 2, 0},
	{0xffff, 0, 0xffff, 0},
	{0xffff, 1, 0, 1},
	{19997, 19379, 39376, 0},
	{0xffff, 0xffff, 0xfffe, 1},
}

var uint32Tests = []struct {
	a, b, sum uint32
	carry     byte
}{
	{0, 0, 0, 0},
	{0, 1, 1, 0},
	{0xffffffff, 0, 0xffffffff, 0},
	{0x11223344, 0, 0x11223344, 0},
	{0xffffffff, 1, 0, 1},
	{0xffffffff, 0xffffffff, 0xfffffffe, 1},
}

var uint8AddEqTests = []struct {
	a, b, carry, sumExpect, carryExpect byte
}{
	{0, 0, 0, 0, 0},
	{0, 0, 1, 1, 0},
	{0, 1, 0, 1, 0},
	{1, 0, 0, 1, 0},
	{0, 0, 2, 0, 0}, //Only carry bit 1 is observed
	{254, 1, 0, 255, 0},
	{254, 2, 0, 0, 1},
	{254, 1, 1, 0, 1},
	{0xff, 0xff, 0, 0xfe, 1},
	{0xff, 0xff, 1, 0xff, 1}, //Max possible output
}

var uint16AddEqTests = []struct {
	a, b        uint16
	carry       byte
	sumExpect   uint16
	carryExpect byte
}{
	{0, 0, 0, 0, 0},
	{0, 0, 1, 1, 0},
	{0, 1, 0, 1, 0},
	{1, 0, 0, 1, 0},
	{0, 0, 2, 0, 0}, //Only carry bit 1 is observed
	{254, 1, 0, 255, 0},
	{254, 2, 0, 256, 0},
	{254, 1, 1, 256, 0},
	{0xff, 0xff, 0, 0x1fe, 0},
	{0xff, 0xff, 1, 0x1ff, 0},
	{65534, 1, 0, 65535, 0},
	{65534, 2, 0, 0, 1},
	{65534, 1, 1, 0, 1},
	{65535, 256, 0, 255, 1},
	{0xffff, 0, 0, 0xffff, 0},
	{0xf0f0, 0x0f0f, 0, 0xffff, 0},
	{0xffff, 0xffff, 0, 0xfffe, 1},
	{0xffff, 0xffff, 1, 0xffff, 1}, //Max possible output
}

var uint32AddEqTests = []struct {
	a, b        uint32
	carry       byte
	sumExpect   uint32
	carryExpect byte
}{
	{0, 0, 0, 0, 0},
	{0, 0, 1, 1, 0},
	{0, 1, 0, 1, 0},
	{1, 0, 0, 1, 0},
	{0, 0, 2, 0, 0}, //Only carry bit 1 is observed
	{254, 1, 0, 255, 0},
	{254, 2, 0, 256, 0},
	{254, 1, 1, 256, 0},
	{0xff, 0xff, 0, 0x1fe, 0},
	{0xff, 0xff, 1, 0x1ff, 0},
	{65534, 1, 0, 65535, 0},
	{65534, 2, 0, 65536, 0},
	{65534, 1, 1, 65536, 0},
	{65535, 256, 0, 65791, 0},
	{0xffff, 0, 0, 0xffff, 0},
	{0xf0f0, 0x0f0f, 0, 0xffff, 0},
	{0xffff, 0xffff, 0, 0x1fffe, 0},
	{0xffff, 0xffff, 1, 0x1ffff, 0},
	{0xffffffff, 0, 0, 0xffffffff, 0},
	{0xf0f0f0f0, 0x0f0f0f0f, 0, 0xffffffff, 0},
	{0xffffffff, 0xffffffff, 0, 0xfffffffe, 1},
	{0xffffffff, 0xffffffff, 1, 0xffffffff, 1}, //Max possible output
}

func bytesEqualAndCarry(t *testing.T, a, b, expect, found []byte, expectCarry, carry byte) {
	if !bytes.Equal(found, expect) {
		test.StringMatchTitle(
			t,
			fmt.Sprintf("\n %v\n+%v\n=", hex.FromBytes(a), hex.FromBytes(b)),
			" ",
			hex.FromBytes(expect),
			hex.FromBytes(found))
	}
	if expectCarry != carry {
		t.Errorf("Expecting carry %d, got %d", expectCarry, carry)
	}
}
func bytesEqualAndCarryHex(t *testing.T, a, b, expectHex string, found []byte, expectCarry, carry byte) {
	expect := hex.ToBytesFast(expectHex)
	if !bytes.Equal(found, expect) {
		test.StringMatchTitle(
			t,
			fmt.Sprintf("\n %s\n+%s\n=", a, b),
			" ",
			hex.FromBytes(expect),
			hex.FromBytes(found))
	}
	if expectCarry != carry {
		t.Errorf("Expecting carry %d, got %d", expectCarry, carry)
	}
}
func sumEqualAndCarry(t *testing.T, a, b uint32, expect, found []byte, expectCarry, carry byte) {
	if !bytes.Equal(found, expect) {
		test.StringMatchTitle(
			t,
			fmt.Sprintf("\n %d\n+%d\n=", a, b),
			" ",
			hex.FromBytes(expect),
			hex.FromBytes(found))
	}
	if expectCarry != carry {
		t.Errorf("Expecting carry %d, got %d", expectCarry, carry)
	}
}
func sumEqualAndCarryHex(t *testing.T, a, b string, expect, found []byte, expectCarry, carry byte) {
	if !bytes.Equal(found, expect) {
		test.StringMatchTitle(
			t,
			fmt.Sprintf("\n %s\n+%s\n=", a, b),
			" ",
			hex.FromBytes(expect),
			hex.FromBytes(found))
	}
	if expectCarry != carry {
		t.Errorf("Expecting carry %d, got %d", expectCarry, carry)
	}
}

func TestAddEq8(t *testing.T) {
	sum := make([]byte, 1)
	sumExpect := make([]byte, 1)
	b := make([]byte, 1)
	var carry byte
	for _, rec := range uint8AddEqTests {
		sumExpect[0] = rec.sumExpect

		sum[0] = rec.a
		b[0] = rec.b
		carry = addEq8(sum, b, rec.carry)
		sumEqualAndCarry(t, uint32(rec.a), uint32(rec.b), sumExpect, sum, rec.carryExpect, carry)

		//Test commutative
		sum[0] = rec.b
		b[0] = rec.a
		carry = addEq8(sum, b, rec.carry)
		sumEqualAndCarry(t, uint32(rec.b), uint32(rec.a), sumExpect, sum, rec.carryExpect, carry)

	}
}

func TestAddEq16LE(t *testing.T) {
	sum := make([]byte, bytes16)
	sumExpect := make([]byte, bytes16)
	b := make([]byte, bytes16)
	var carry byte
	for _, rec := range uint16AddEqTests {
		binary.LittleEndian.PutUint16(sumExpect, rec.sumExpect)

		binary.LittleEndian.PutUint16(sum, rec.a)
		binary.LittleEndian.PutUint16(b, rec.b)
		carry = addEq16LE(sum, b, rec.carry)
		sumEqualAndCarry(t, uint32(rec.a), uint32(rec.b), sumExpect, sum, rec.carryExpect, carry)

		//Test commutative
		binary.LittleEndian.PutUint16(sum, rec.b)
		binary.LittleEndian.PutUint16(b, rec.a)
		carry = addEq16LE(sum, b, rec.carry)
		sumEqualAndCarry(t, uint32(rec.b), uint32(rec.a), sumExpect, sum, rec.carryExpect, carry)
	}

	var hexTests = []struct {
		a, b        string
		carry       byte
		sumExpect   uint16
		carryExpect byte
	}{
		{"FFFF", "", 0, 0xffff, 0},
		{"FFFF", "00", 0, 0xffff, 0},
		{"FFFF", "0000", 0, 0xffff, 0},
		{"FFFF", "000001", 0, 0xffff, 0}, //Last byte of b ignored
		{"FFFF", "01", 0, 0, 1},
		{"FFFF", "0001", 0, 0x00ff, 1}, //65535+256=65791
		{"FFFF", "0100", 0, 0, 1},
	}
	for _, rec := range hexTests {
		binary.LittleEndian.PutUint16(sumExpect, rec.sumExpect)

		padLE(sum, hex.ToBytesFast(rec.a), bytes16)
		//Note we want a fresh b2 because we want to test irregular length byte slices
		b2 := hex.ToBytesFast(rec.b)
		carry = addEq16LE(sum, b2, rec.carry)
		sumEqualAndCarryHex(t, rec.a, rec.b, sumExpect, sum, rec.carryExpect, carry)

		//Test commutative, although sum must always be the right size, so use padLE
		padLE(sum, hex.ToBytesFast(rec.b), bytes16)
		b2 = hex.ToBytesFast(rec.a)
		carry = addEq16LE(sum, b2, rec.carry)
		sumEqualAndCarryHex(t, rec.b, rec.a, sumExpect, sum, rec.carryExpect, carry)
	}
}

func TestAddEq32LE(t *testing.T) {
	sum := make([]byte, bytes32)
	sumExpect := make([]byte, bytes32)
	b := make([]byte, bytes32)
	var carry byte
	for _, rec := range uint32AddEqTests {
		binary.LittleEndian.PutUint32(sumExpect, rec.sumExpect)

		binary.LittleEndian.PutUint32(sum, rec.a)
		binary.LittleEndian.PutUint32(b, rec.b)
		carry = addEq32LE(sum, b, rec.carry)
		sumEqualAndCarry(t, rec.a, rec.b, sumExpect, sum, rec.carryExpect, carry)

		//Test commutative
		binary.LittleEndian.PutUint32(sum, rec.b)
		binary.LittleEndian.PutUint32(b, rec.a)
		carry = addEq32LE(sum, b, rec.carry)
		sumEqualAndCarry(t, rec.b, rec.a, sumExpect, sum, rec.carryExpect, carry)
	}

	var hexTests = []struct {
		a, b        string
		carry       byte
		sumExpect   uint32
		carryExpect byte
	}{
		{"FFFFFFFF", "", 0, 0xffffffff, 0},
		{"FFFFFFFF", "00", 0, 0xffffffff, 0},
		{"FFFFFFFF", "0000", 0, 0xffffffff, 0},
		{"FFFFFFFF", "000000", 0, 0xffffffff, 0},
		{"FFFFFFFF", "00000000", 0, 0xffffffff, 0},
		{"FFFFFFFF", "0000000001", 0, 0xffffffff, 0}, //Last byte of b ignored
		{"FFFFFFFF", "01", 0, 0, 1},
		{"FFFFFFFF", "0001", 0, 0x000000ff, 1},
		{"FFFFFFFF", "000001", 0, 0x0000ffff, 1},
		{"FFFFFFFF", "00000001", 0, 0x00ffffff, 1},
		{"FFFFFFFF", "0100", 0, 0, 1},
	}
	for _, rec := range hexTests {
		binary.LittleEndian.PutUint32(sumExpect, rec.sumExpect)

		padLE(sum, hex.ToBytesFast(rec.a), bytes32)
		//Note we want a fresh b2 because we want to test irregular length byte slices
		b2 := hex.ToBytesFast(rec.b)
		carry = addEq32LE(sum, b2, rec.carry)
		sumEqualAndCarryHex(t, rec.a, rec.b, sumExpect, sum, rec.carryExpect, carry)

		//Test commutative, although sum must always be the right size, so use padLE
		padLE(sum, hex.ToBytesFast(rec.b), bytes32)
		b2 = hex.ToBytesFast(rec.a)
		carry = addEq32LE(sum, b2, rec.carry)
		sumEqualAndCarryHex(t, rec.b, rec.a, sumExpect, sum, rec.carryExpect, carry)
	}
}

func TestAdd16LE(t *testing.T) {
	sum := make([]byte, bytes16)
	a := make([]byte, bytes16)
	b := make([]byte, bytes16)
	expectSum := make([]byte, bytes16)
	var carry byte
	for _, rec := range uint16Tests {
		binary.LittleEndian.PutUint16(a, rec.a)
		binary.LittleEndian.PutUint16(b, rec.b)
		binary.LittleEndian.PutUint16(expectSum, rec.sum)

		carry, _ = Add16LE(sum, a, b)
		bytesEqualAndCarry(t, a, b, expectSum, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add16LE(sum, b, a)
		bytesEqualAndCarry(t, b, a, expectSum, sum, rec.carry, carry)
	}

	var byteTests = []struct {
		a, b, sum []byte
		carry     byte
	}{
		{[]byte{}, []byte{}, []byte{0, 0}, 0},
		{[]byte{1}, []byte{}, []byte{1, 0}, 0},
		{[]byte{}, []byte{1}, []byte{1, 0}, 0},
		{[]byte{1, 2}, []byte{}, []byte{1, 2}, 0},
		{[]byte{1, 2}, []byte{1}, []byte{2, 2}, 0},
		{[]byte{0xff, 1}, []byte{1}, []byte{0, 2}, 0},
		{[]byte{0xff, 0xff}, []byte{1}, []byte{0, 0}, 1},

		//Max overflow
		{[]byte{0xff, 0xff}, []byte{0xff, 0xff}, []byte{0xfe, 0xff}, 1},
	}
	for _, rec := range byteTests {
		carry, _ = Add16LE(sum, rec.a, rec.b)
		bytesEqualAndCarry(t, a, b, rec.sum, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add16LE(sum, rec.b, rec.a)
		bytesEqualAndCarry(t, b, a, rec.sum, sum, rec.carry, carry)
	}
}

func TestAdd32LE(t *testing.T) {
	sumExpect := make([]byte, bytes32)

	sum := make([]byte, bytes32)
	a := make([]byte, bytes32)
	b := make([]byte, bytes32)
	var carry byte
	for _, rec := range uint32Tests {
		binary.LittleEndian.PutUint32(a, rec.a)
		binary.LittleEndian.PutUint32(b, rec.b)
		binary.LittleEndian.PutUint32(sumExpect, rec.sum)

		carry, _ = Add32LE(sum, a, b)
		bytesEqualAndCarry(t, a, b, sumExpect, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add32LE(sum, b, a)
		bytesEqualAndCarry(t, b, a, sumExpect, sum, rec.carry, carry)
	}

	var hexTests = []struct {
		a, b  string
		sum   uint32
		carry byte
	}{
		{"", "", 0, 0},
		{"00", "", 0, 0},
		{"00", "00", 0, 0},
		{"0000", "", 0, 0},
		{"0000", "00", 0, 0},
		{"0000", "0000", 0, 0},
		{"01", "", 1, 0},
		{"0102", "", 0x0201, 0},
		{"0102", "01", 0x0202, 0},
		{"FF01", "01", 0x0200, 0},
		{"FFFF", "01", 0x010000, 0},
		{"FFFF", "FFFF", 0x1fffe, 0},

		{"FFFFFFFF", "", 0xffffffff, 0},
		{"FFFFFFFF", "00", 0xffffffff, 0},
		{"FFFFFFFF", "0000", 0xffffffff, 0},
		{"FFFFFFFF", "000000", 0xffffffff, 0},
		{"FFFFFFFF", "00000000", 0xffffffff, 0},
		{"FFFFFFFF", "0000000001", 0xffffffff, 0}, //Last byte of b ignored
		{"FFFFFFFF", "01", 0, 1},
		{"FFFFFFFF", "0001", 0x000000ff, 1},
		{"FFFFFFFF", "000001", 0x0000ffff, 1},
		{"FFFFFFFF", "00000001", 0x00ffffff, 1},
		{"FFFFFFFF", "0100", 0, 1},
		{"FFFFFFFF", "010000", 0, 1},
		{"FFFFFFFF", "01000000", 0, 1},
		{"FFFFFFFF", "FFFFFFFF", 0xfffffffe, 1}, //Nax sum
	}
	for _, rec := range hexTests {
		binary.LittleEndian.PutUint32(sumExpect, rec.sum)

		a2 := hex.ToBytesFast(rec.a)
		b2 := hex.ToBytesFast(rec.b)
		carry, _ = Add32LE(sum, a2, b2)
		sumEqualAndCarryHex(t, rec.a, rec.b, sumExpect, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add32LE(sum, b2, a2)
		sumEqualAndCarryHex(t, rec.b, rec.a, sumExpect, sum, rec.carry, carry)
	}
}

// 128,256,512 use similar processing

func TestAdd128LE(t *testing.T) {
	sum := make([]byte, bytes128)
	var carry byte

	var hexTests = []struct {
		a, b, sum string
		carry     byte
	}{
		{
			"",
			"",
			"00000000000000000000000000000000",
			0},
		{
			"00",
			"",
			"00000000000000000000000000000000",
			0},
		{
			"01",
			"",
			"01000000000000000000000000000000",
			0},
		{
			"0001",
			"",
			"00010000000000000000000000000000",
			0},
		{
			"0100",
			"",
			"01000000000000000000000000000000",
			0},
		{
			"000102",
			"",
			"00010200000000000000000000000000",
			0},
		{
			"00010203",
			"",
			"00010203000000000000000000000000",
			0},
		{
			"0001020304",
			"",
			"00010203040000000000000000000000",
			0},
		{
			"000102030405",
			"",
			"00010203040500000000000000000000",
			0},
		{
			"00010203040506",
			"",
			"00010203040506000000000000000000",
			0},
		{
			"0001020304050607",
			"",
			"00010203040506070000000000000000",
			0},
		{
			"00010203040506070809",
			"",
			"00010203040506070809000000000000",
			0},
		{
			"000102030405060708090A0B",
			"",
			"000102030405060708090A0B00000000",
			0},
		{
			"000102030405060708090A0B0C0D",
			"",
			"000102030405060708090A0B0C0D0000",
			0},
		{
			"000102030405060708090A0B0C0D0E",
			"",
			"000102030405060708090A0B0C0D0E00",
			0},
		{
			"000102030405060708090A0B0C0D0E0F",
			"",
			"000102030405060708090A0B0C0D0E0F",
			0},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			0},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"01",
			"00000000000000000000000000000000",
			1},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"0001",
			"FF000000000000000000000000000000",
			1},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"000001",
			"FFFF0000000000000000000000000000",
			1},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"FEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			1},
		{
			"FF0000000000000000000000000000FF",
			"00000000000000000000000000000001",
			"FF000000000000000000000000000000",
			1},
		{
			"FF0000000000000000000000000000FF",
			"01000000000000000000000000000000",
			"000100000000000000000000000000FF",
			0},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF00",
			"00FF00FF00FF00FF00FF00FF00FF00FF",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			0},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF00",
			"01FF00FF00FF00FF00FF00FF00FF00FF",
			"00000000000000000000000000000000",
			1},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF00",
			"00FF00FF00FF00FF01FF00FF00FF00FF",
			"FFFFFFFFFFFFFFFF0000000000000000",
			1},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF01",
			"00FF00FF00FF00FF00FF00FF00FF00FF",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00",
			1},
	}
	for _, rec := range hexTests {
		a := hex.ToBytesFast(rec.a)
		b := hex.ToBytesFast(rec.b)
		carry, _ = Add128LE(sum, a, b)
		bytesEqualAndCarryHex(t, rec.a, rec.b, rec.sum, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add128LE(sum, b, a)
		bytesEqualAndCarryHex(t, rec.b, rec.a, rec.sum, sum, rec.carry, carry)
	}
}

func TestAdd256LE(t *testing.T) {
	sum := make([]byte, bytes256)
	var carry byte

	var hexTests = []struct {
		a, b, sum string
		carry     byte
	}{
		{
			"",
			"",
			"0000000000000000000000000000000000000000000000000000000000000000",
			0},
		{
			"03",
			"05",
			"0800000000000000000000000000000000000000000000000000000000000000",
			0},
	}
	for _, rec := range hexTests {
		a := hex.ToBytesFast(rec.a)
		b := hex.ToBytesFast(rec.b)
		carry, _ = Add256LE(sum, a, b)
		bytesEqualAndCarryHex(t, rec.a, rec.b, rec.sum, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add256LE(sum, b, a)
		bytesEqualAndCarryHex(t, rec.b, rec.a, rec.sum, sum, rec.carry, carry)
	}
}

func TestAdd512LE(t *testing.T) {
	sum := make([]byte, bytes512)
	var carry byte

	var hexTests = []struct {
		a, b, sum string
		carry     byte
	}{
		{
			"",
			"",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			0},
		{
			"01",
			"",
			"01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 0},
		{
			"00",
			"01",
			"01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 0},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"01",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 1},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00", 1},
		{ //Same as +1
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 1},
		{
			"0123456789ABCDEF",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000123456789ABCDEF",
			"0123456789ABCDEF0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000123456789ABCDEF", 0},
		{
			"0123456789ABCDEF",
			"0123456789ABCDEF0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"02468ACE12579BDF0100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 0},
		{
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001",
			"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000FF",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 1},
		{
			"01000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"FF000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			"00010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 0},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00",
			"00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 0},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00",
			"01FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF",
			"00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 1},
		{
			"FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF01",
			"00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF00FF",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF00", 1},
		{
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF",
			"FEFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFF", 1},
		{
			"03030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303",
			"05050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505",
			"08080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808", 0},
		{
			"03030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303",
			"050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505",
			"08080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080803", 0},
		{
			"03030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303",
			"050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505",
			"08080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080803030303", 0},
		{
			"03030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303",
			"0505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505",
			"08080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080303030303", 0},
		{
			"03030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303030303",
			"0505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505050505",
			"08080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080808080303030303030303", 0},
		{
			"0100",
			"22",
			"23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 0},
		{
			"0001",
			"22",
			"22010000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", 0},
	}
	for _, rec := range hexTests {
		a := hex.ToBytesFast(rec.a)
		b := hex.ToBytesFast(rec.b)
		carry, _ = Add512LE(sum, a, b)
		bytesEqualAndCarryHex(t, rec.a, rec.b, rec.sum, sum, rec.carry, carry)

		//Test commutative
		carry, _ = Add512LE(sum, b, a)
		bytesEqualAndCarryHex(t, rec.b, rec.a, rec.sum, sum, rec.carry, carry)
	}
}
