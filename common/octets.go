package common

import (
	"fmt"
	"math"
	"reflect"
)

type Octets struct {
	Data  []byte
	Begin int32
	End   int32
}

func Wrap(bytes []byte) *Octets {
	octets := Octets{}
	octets.Data = bytes
	return &octets
}

func (b *Octets) length() int32 {
	return int32(len(b.Data))
}

func (b *Octets) toInt32(i int32, an byte, lmv uint32) int32 {
	ret := int32(b.Data[i]) & int32(an)
	if lmv > 0 {
		ret <<= lmv
	}
	return ret
}

func byteMvInt32(b byte, an byte, lmv uint32) int32 {
	ret := int32(b)
	if an > 0 {
		ret &= int32(an)
	}
	if lmv > 0 {
		ret <<= lmv
	}
	return ret
}

func byteMvUInt64(b byte, an byte, lmv uint64) uint64 {
	ret := uint64(b)
	if an > 0 {
		ret &= uint64(an)
	}
	if lmv > 0 {
		ret <<= lmv
	}
	return ret
}

func byteMvUInt32(b byte, an byte, lmv uint32) uint32 {
	ret := uint32(b)
	if an > 0 {
		ret &= uint32(an)
	}
	if lmv > 0 {
		ret <<= lmv
	}
	return ret
}

func byteMvInt64(b byte, an byte, lmv uint32) int64 {
	ret := int64(b)
	if an > 0 {
		ret &= int64(an)
	}
	if lmv > 0 {
		ret <<= lmv
	}
	return ret
}

func (b *Octets) WriteMsgSize(x int32) {
	if x > 0 {
		if x < 0x80 {
			b.WriteByte(byte(x))
		} else if x < 0x4000 {
			b.WriteInt(x | 0x8000)
		} else if x < 0x200000 {
			b.WriteInt(x | 0xc00000)
		} else if x < 0x10000000 {
			b.WriteInt(int32(int(x) | 0xe0000000))
		}
	}

	/**
		if (x >= 0) {
	            if (x < 0x80) {
	                byteBuf.WriteByte(x);
	            } else if (x < 0x4000) {
	                byteBuf.WriteShort(x | 0x8000);
	            } else if (x < 0x200000) {
	                byteBuf.WriteMedium(x | 0xc00000);
	            } else if (x < 0x10000000) {
	                byteBuf.WriteInt(x | 0xe0000000);
	            } else {
	                throw new RuntimeException("exceed max unit");
	            }
	        }
	*/
}

func (b *Octets) WriteLong(x int64) {
	if x > 0 {
		if x < 0x80 {
			log("11 Write ", x)
			b.Data = append(b.Data, byte(x))
			return
		} else if x < 0x4000 {
			log("22 Write ", x)
			b.Data = append(b.Data, byte(x>>8|0x80), byte(x))
			return
		} else if x < 0x200000 {
			log("33 Write ", x)
			b.Data = append(b.Data, byte(x>>16|0xc0), byte(x>>8), byte(x))
			return
		} else if x < 0x10000000 {
			log("44 Write ", x)
			b.Data = append(b.Data, byte((x>>24)|0xe0), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x800000000 {
			b.Data = append(b.Data, byte((x>>32)|0xf0), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x40000000000 {
			b.Data = append(b.Data, byte((x>>40)|0xf8), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x200000000000 {
			b.Data = append(b.Data, byte((x>>48)|0xfc), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x100000000000000 {
			b.Data = append(b.Data, byte(0xfe), byte(x>>48), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		}
	}
	b.Data = append(b.Data, byte(0xff), byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
}

func (b *Octets) ReadLong() int64 {

	//fmt.Println("0xff ", reflect.TypeOf(0xff))

	h := byte(b.Data[b.Begin] & 0xff)

	//fmt.Println(" h  ", reflect.TypeOf(h))

	if h < 0x80 {
		b.Begin++
		log(" 111 Begin ", b.Begin)
		return int64(h)
	} else if h < 0xc0 {
		x := byteMvInt64(h, 0x3f, 8) | byteMvInt64(b.Data[b.Begin+1], 0xff, 0)
		b.Begin += 2
		log(" 222 Begin ", b.Begin)
		return x
	} else if h < 0xe0 {
		x := byteMvInt64(h, 0x1f, 16) | byteMvInt64(b.Data[b.Begin+1], 0xff, 8) | byteMvInt64(b.Data[b.Begin+2], 0xff, 0)
		b.Begin += 3
		log(" 333 Begin ", b.Begin)
		return x
	} else if h < 0xf0 {
		x := byteMvInt64(h, 0x0f, 24) | byteMvInt64(b.Data[b.Begin+1], 0xff, 16) | byteMvInt64(b.Data[b.Begin+2], 0xff, 8) | byteMvInt64(b.Data[b.Begin+3], 0xff, 0)
		b.Begin += 4
		log(" 444 Begin ", b.Begin)
		return x
	} else if h < 0xf8 {
		xl := byteMvInt64(b.Data[b.Begin+1], 0, 24) | byteMvInt64(b.Data[b.Begin+2], 0xff, 16) | byteMvInt64(b.Data[b.Begin+3], 0xff, 8) | byteMvInt64(b.Data[b.Begin+4], 0xff, 0)
		xh := byteMvInt64(h, 0x07, 0)
		b.Begin += 5
		return xh<<32 | xl&0xffffffff
	} else if h < 0xfc {
		xl := byteMvInt64(b.Data[b.Begin+2], 0, 24) | byteMvInt64(b.Data[b.Begin+3], 0xff, 16) | byteMvInt64(b.Data[b.Begin+4], 0xff, 8) | byteMvInt64(b.Data[b.Begin+5], 0xff, 0)
		xh := byteMvInt64(h, 0x03, 8) | byteMvInt64(b.Data[b.Begin+1], 0xff, 0)
		b.Begin += 6
		return xh<<32 | xl&0xffffffff
	} else if h < 0xfe {
		xl := byteMvInt64(b.Data[b.Begin+3], 0, 24) | byteMvInt64(b.Data[b.Begin+4], 0xff, 16) | byteMvInt64(b.Data[b.Begin+5], 0xff, 8) | byteMvInt64(b.Data[b.Begin+6], 0xff, 0)
		xh := byteMvInt64(h, 0x01, 16) | byteMvInt64(b.Data[b.Begin+1], 0xff, 8) | byteMvInt64(b.Data[b.Begin+2], 0xff, 0)
		b.Begin += 7
		return xh<<32 | xl&0xffffffff
	} else if h < 0xff {
		xl := byteMvInt64(b.Data[b.Begin+4], 0, 24) | byteMvInt64(b.Data[b.Begin+5], 0xff, 16) | byteMvInt64(b.Data[b.Begin+6], 0xff, 8) | byteMvInt64(b.Data[b.Begin+7], 0xff, 0)
		xh := byteMvInt64(b.Data[b.Begin+1], 0xff, 16) | byteMvInt64(b.Data[b.Begin+2], 0xff, 8) | byteMvInt64(b.Data[b.Begin+3], 0xff, 0)
		b.Begin += 8
		return xh<<32 | xl&0xffffffff
	} else {
		xl := byteMvInt64(b.Data[b.Begin+5], 0, 24) | byteMvInt64(b.Data[b.Begin+6], 0xff, 16) | byteMvInt64(b.Data[b.Begin+7], 0xff, 8) | byteMvInt64(b.Data[b.Begin+8], 0xff, 0)
		xh := byteMvInt64(b.Data[b.Begin+1], 0, 24) | byteMvInt64(b.Data[b.Begin+2], 0xff, 16) | byteMvInt64(b.Data[b.Begin+3], 0xff, 8) | byteMvInt64(b.Data[b.Begin+4], 0xff, 0)
		b.Begin += 9
		return xh<<32 | xl&0xffffffff
	}
}

func (b *Octets) WriteInt(x int32) {
	if x >= 0 {
		if x < 0x80 {
			log("11 Write ", x)
			b.Data = append(b.Data, byte(x))
			return
		} else if x < 0x4000 {
			log("22 Write ", x)
			b.Data = append(b.Data, byte(x>>8|0x80), byte(x))
			return
		} else if x < 0x200000 {
			log("33 Write ", x)
			b.Data = append(b.Data, byte(x>>16|0xc0), byte(x>>8), byte(x))
			return
		} else if x < 0x10000000 {
			log("44 Write ", x)
			b.Data = append(b.Data, byte((x>>24)|0xe0), byte(x>>16), byte(x>>8), byte(x))
			return
		}
	}
	log("55 Write ", x)

	b.Data = append(b.Data, 0xf0, byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
}

func (b *Octets) ReadInt() int32 {
	h := byte(b.Data[b.Begin] & 0xff)

	log(" Read h = ", h, reflect.TypeOf(h))

	//t := int32(h) << 24

	//log(" Read t = ", t, reflect.TypeOf(t))

	if h < 0x80 {
		b.Begin++
		log(" 111 Begin ", b.Begin)
		return int32(h)
	} else if h < 0xc0 {
		//x := int32(((h & 0x3f) << 8) | (b.Data[b.Begin+1] & 0xff))
		x := byteMvInt32(h, 0x3f, 8) | byteMvInt32(b.Data[b.Begin+1], 0xff, 0)
		b.Begin += 2
		log(" 222 Begin ", b.Begin)
		return x
	} else if h < 0xe0 {
		//x := int32(((h & 0x1f) << 16) | ((b.Data[b.Begin+1] & 0xff) << 8) | (b.Data[b.Begin+2] & 0xff))
		x := byteMvInt32(h, 0x1f, 16) | byteMvInt32(b.Data[b.Begin+1], 0xff, 8) | byteMvInt32(b.Data[b.Begin+2], 0xff, 0)
		b.Begin += 3
		log(" 333 Begin ", b.Begin)
		return x
	} else if h < 0xf0 {
		//x := int32(((h & 0x0f) << 24) | ((b.Data[b.Begin+1] & 0xff) << 16) | ((b.Data[b.Begin+2] & 0xff) << 8) | (b.Data[b.Begin+3] & 0xff))
		x := byteMvInt32(h, 0x0f, 24) | byteMvInt32(b.Data[b.Begin+1], 0xff, 16) | byteMvInt32(b.Data[b.Begin+2], 0xff, 8) | byteMvInt32(b.Data[b.Begin+3], 0xff, 0)
		b.Begin += 4
		log(" 444 Begin ", b.Begin)
		return x
	} else {
		//x := int32(((b.Data[b.Begin+1] & 0xff) << 24) | ((b.Data[b.Begin+2] & 0xff) << 16) | ((b.Data[b.Begin+3] & 0xff) << 8) | (b.Data[b.Begin+4] & 0xff))
		x := byteMvInt32(b.Data[b.Begin+1], 0xff, 24) | byteMvInt32(b.Data[b.Begin+2], 0xff, 16) | byteMvInt32(b.Data[b.Begin+3], 0xff, 8) | byteMvInt32(b.Data[b.Begin+4], 0xff, 0)
		b.Begin += 5
		log(" 555 Begin ", b.Begin)
		return x
	}
}

func (b *Octets) WriteBool(x bool) {
	if x {
		b.Data = append(b.Data, byte(1))
	} else {
		b.Data = append(b.Data, byte(0))
	}
}

func (b *Octets) ReadBool() bool {
	x := b.Data[b.Begin]
	b.Begin++
	if x > 0 {
		return true
	} else {
		return false
	}
}

func (b *Octets) WriteByte(x byte) {
	b.Data = append(b.Data, x)
}

func (b *Octets) ReadByte() byte {
	x := b.Data[b.Begin]
	return x
}

func (b *Octets) WriteBytes(x []byte) {
	b.WriteInt(int32(len(x)))
	b.Data = append(b.Data, x...)
}

func (b *Octets) ReadBytes() []byte {
	length := b.ReadInt()
	bytes := b.Data[b.Begin:length]
	b.Begin += length
	return bytes
}

func (b *Octets) WriteFloat(x float32) {
	bits := math.Float32bits(x)
	b.Data = append(b.Data, byte(bits&0xff))
	b.Data = append(b.Data, byte((bits>>8)&0xff))
	b.Data = append(b.Data, byte((bits>>16)&0xff))
	b.Data = append(b.Data, byte((bits>>24)&0xff))
}

func (b *Octets) ReadFloat() float32 {
	i1 := byteMvUInt32(b.Data[b.Begin], 0xff, 0)
	i2 := byteMvUInt32(b.Data[b.Begin+1], 0xff, 8)
	i3 := byteMvUInt32(b.Data[b.Begin+2], 0xff, 16)
	i4 := byteMvUInt32(b.Data[b.Begin+3], 0xff, 24)
	b.Begin += 4
	return math.Float32frombits(i1 + i2 + i3 + i4)
}

func (b *Octets) WriteDouble(x float64) {
	bits := math.Float64bits(x)
	b.Data = append(b.Data, byte(bits))
	b.Data = append(b.Data, byte(bits>>8))
	b.Data = append(b.Data, byte(bits>>16))
	b.Data = append(b.Data, byte(bits>>24))
	b.Data = append(b.Data, byte(bits>>32))
	b.Data = append(b.Data, byte(bits>>40))
	b.Data = append(b.Data, byte(bits>>48))
	b.Data = append(b.Data, byte(bits>>56))
}

func (b *Octets) ReadDouble() float64 {

	ui1 := byteMvUInt64(b.Data[b.Begin], 0xff, 0)
	ui2 := byteMvUInt64(b.Data[b.Begin+1], 0xff, 8)
	ui3 := byteMvUInt64(b.Data[b.Begin+2], 0xff, 16)
	ui4 := byteMvUInt64(b.Data[b.Begin+3], 0xff, 24)
	ui5 := byteMvUInt64(b.Data[b.Begin+4], 0xff, 32)
	ui6 := byteMvUInt64(b.Data[b.Begin+5], 0xff, 40)
	ui7 := byteMvUInt64(b.Data[b.Begin+6], 0xff, 48)
	ui8 := byteMvUInt64(b.Data[b.Begin+7], 0xff, 56)

	b.Begin += 8
	return math.Float64frombits(ui1 + ui2 + ui3 + ui4 + ui5 + ui6 + ui7 + ui8)
}

func (b *Octets) WriteString(x string) {
	length := len(x)
	b.WriteInt(int32(length))
	strBytes := []byte(x)
	b.Data = append(b.Data, strBytes...)

}

func (b *Octets) ReadString() string {
	length := b.ReadInt()
	bytes := b.Data[b.Begin : b.Begin+length]
	b.Begin += length
	return string(bytes[:])
}

func (b *Octets) WriteOctets(src *Octets) {
	srcLength := len(src.Data)
	b.WriteInt(int32(srcLength))
	b.WriteBytes(src.Data)
}

func (b *Octets) ReadOctets() *Octets {
	length := b.ReadInt()
	ret := Octets{}
	ret.Data = append(ret.Data, b.Data[b.Begin:b.Begin+length]...)
	//copy(ret.Data, )
	b.Begin += length
	return &ret
}

func (b *Octets) WriteShort(x int16) {

}

func (b *Octets) ReadShort() int16 {
	return -1
}

func (b *Octets) WriteSInt(x int32) {
	b.WriteInt(x)
}

func (b *Octets) ReadSInt() int32 {
	return b.ReadInt()
}

func (b *Octets) WriteSLong(x int64) {
	b.WriteSLong(x)
}
func (b *Octets) ReadSLong() int64 {
	return b.ReadLong()
}

func (b *Octets) WriteCompactUint(x int32) {
	b.WriteInt(x)
}

func (b *Octets) ReadCompactUint() int32 {
	return b.ReadInt()
}

func runTest() {

	myBuffer := Octets{}

	//myBuffer.Data = append(myBuffer.Data, 12312)

	myBuffer.WriteInt(7)
	myBuffer.WriteInt(222)
	myBuffer.WriteInt(33333)
	myBuffer.WriteInt(4444444)
	myBuffer.WriteInt(555555555)

	myBuffer.WriteInt(-7)
	myBuffer.WriteInt(-222)
	myBuffer.WriteInt(-33333)
	myBuffer.WriteInt(-4444444)
	myBuffer.WriteInt(-555555555)

	fmt.Println(myBuffer, len(myBuffer.Data), cap(myBuffer.Data))

	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	/////

	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())
	fmt.Println(myBuffer.ReadInt())

	fmt.Println(myBuffer, len(myBuffer.Data), cap(myBuffer.Data))

	longBuff := Octets{}

	longBuff.WriteLong(7)
	longBuff.WriteLong(222)
	longBuff.WriteLong(33333)
	longBuff.WriteLong(4444444)
	longBuff.WriteLong(555555555)
	longBuff.WriteLong(6666666666)
	longBuff.WriteLong(77777777777)
	longBuff.WriteLong(888888888888)
	longBuff.WriteLong(9999999999999)
	longBuff.WriteLong(19999999999999)
	longBuff.WriteLong(129999999999999)
	longBuff.WriteLong(1239999999999999)
	longBuff.WriteLong(12349999999999999)
	longBuff.WriteLong(123459999999999999)
	longBuff.WriteLong(1234569999999999999)
	//longBuff.WriteLong(12345679999999999999)
	//longBuff.WriteLong(123456789999999999999)
	//longBuff.WriteLong(1234567899999999999999)
	//longBuff.WriteLong(12345678909999999999999)

	longBuff.WriteLong(-7)
	longBuff.WriteLong(-222)
	longBuff.WriteLong(-33333)
	longBuff.WriteLong(-4444444)
	longBuff.WriteLong(-555555555)
	longBuff.WriteLong(-6666666666)
	longBuff.WriteLong(-77777777777)
	longBuff.WriteLong(-888888888888)
	longBuff.WriteLong(-9999999999999)
	longBuff.WriteLong(-19999999999999)
	longBuff.WriteLong(-129999999999999)
	longBuff.WriteLong(-1239999999999999)
	longBuff.WriteLong(-12349999999999999)
	longBuff.WriteLong(-123459999999999999)
	longBuff.WriteLong(-1234569999999999999)

	fmt.Println(longBuff, len(longBuff.Data), cap(longBuff.Data))

	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())

	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())
	//fmt.Println(longBuff.ReadLong())

	fmt.Println(longBuff, len(longBuff.Data), cap(longBuff.Data))

	mi := 0x800
	mi >>= 7
	fmt.Println(mi)

}

//func main() {
//	runTest()
//}

func log(a ...interface{}) {
	//fmt.Println(a)
}
