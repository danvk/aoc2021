package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

var BITS map[string]string

func Init() {
	BITS = map[string]string{
		"0": "0000",
		"1": "0001",
		"2": "0010",
		"3": "0011",
		"4": "0100",
		"5": "0101",
		"6": "0110",
		"7": "0111",
		"8": "1000",
		"9": "1001",
		"A": "1010",
		"B": "1011",
		"C": "1100",
		"D": "1101",
		"E": "1110",
		"F": "1111",
	}
}

type BitString []string

func (b BitString) NumAt(start int, numBits int) int64 {
	bits := b[start : start+numBits]
	str := strings.Join(bits, "")
	num, err := strconv.ParseInt(str, 2, 64)
	if err != nil {
		panic(err)
	}
	return num
}

func (b *BitString) PopBits(numBits int64) BitString {
	oldB := *b
	*b = (*b)[numBits:]
	return oldB[:numBits]
}

func (b BitString) Decode() int64 {
	return b.NumAt(0, len(b))
}

func (b *BitString) PopAndDecode(numBits int64) int64 {
	leadBits := b.PopBits(numBits)
	return leadBits.Decode()
}

type Packet struct {
	version int64
	typeId  int64
	value   int64
	packets []Packet
}

func (p Packet) VersionSum() int64 {
	var sum int64 = p.version
	for _, packet := range p.packets {
		sum += packet.VersionSum()
	}
	return sum
}

func (b *BitString) DecodePacket() Packet {
	version := b.PopAndDecode(3)
	typeId := b.PopAndDecode(3)
	if typeId == 4 {
		var numBits BitString
		for {
			bits := b.PopBits(5)
			numBits = append(numBits, bits[1:]...)
			if bits[0] == "0" {
				break
			}
		}
		value := numBits.Decode()
		return Packet{
			version,
			typeId,
			value,
			nil,
		}
	} else {
		// operator
		lengthTypeId := b.PopBits(1)[0]
		fmt.Printf("Length type id: %s\n", lengthTypeId)
		var subPackets []Packet
		if lengthTypeId == "0" {
			// the next 15 bits are a number that represents the total length in bits of the sub-packets contained by this packet.
			lengthBits := b.PopAndDecode(15)
			subPacketBits := b.PopBits(lengthBits)
			for len(subPacketBits) > 0 {
				subPackets = append(subPackets, subPacketBits.DecodePacket())
			}
		} else {
			// the next 11 bits are a number that represents the number of sub-packets immediately contained by this packet.
			numSubPackets := b.PopAndDecode(11)
			var i int64 = 0
			for ; i < numSubPackets; i++ {
				subPackets = append(subPackets, b.DecodePacket())
			}
		}
		return Packet{
			version,
			typeId,
			-1,
			subPackets,
		}
	}
}

func main() {
	Init()
	hex := os.Args[1]
	chars := strings.Split(hex, "")
	var bits BitString
	for _, char := range chars {
		bit, ok := BITS[char]
		if !ok {
			panic(char)
		}
		bits = append(bits, strings.Split(bit, "")...)
	}

	packet := bits.DecodePacket()
	fmt.Printf("Read %s -> %#v\n", hex, packet)
	fmt.Printf("Version Sum: %d\n", packet.VersionSum())
}
