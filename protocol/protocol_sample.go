package protocol

import (
	"github.com/giskook/inl2_das_sample/base"
)

type SamplePacket struct {
	TagMac  string
	Rssi    int8
	RingMac string
	DegreeX int8
	DegreeY int8
	DegreeZ int8
	Bett    uint8
}

func (p *SamplePacket) Serialize() []byte {
	return nil
}

func ParseSample(buffer []byte) *SamplePacket {
	reader, _, _ := ParseHeader(buffer)
	ring_mac := base.ReadMac(reader)
	_degree_x, _ := reader.ReadByte()
	degree_x := int8(_degree_x)
	_degree_y, _ := reader.ReadByte()
	degree_y := int8(_degree_y)
	_degree_z, _ := reader.ReadByte()
	degree_z := int8(_degree_z)
	bett, _ := reader.ReadByte()
	reader.ReadByte()
	reader.ReadByte()
	tag_mac := base.ReadMac(reader)
	_rssi, _ := reader.ReadByte()
	rssi := int8(_rssi)

	return &SamplePacket{
		TagMac:  tag_mac,
		Rssi:    rssi,
		RingMac: ring_mac,
		DegreeX: degree_x,
		DegreeY: degree_y,
		DegreeZ: degree_z,
		Bett:    bett,
	}
}
