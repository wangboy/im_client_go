package main

import (
	"fmt"
	"math"
	"reflect"
	"time"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		time.Sleep(100 * time.Millisecond)
		log(s)
	}
}

func sum(s []int, c chan int) {
	sum := 0
	for _, v := range s {
		sum += v
	}
	c <- sum // 把 sum 发送到通道 c
}

func fibonacci(n int, c chan int) {
	x, y := 0, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y
	}
	close(c)
}

type ImBuffer struct {
	data  []byte
	begin int32
	end   int32
}

func (b *ImBuffer) length() int32 {
	return int32(len(b.data))
}

func (b *ImBuffer) toInt32(i int32, an byte, lmv uint32) int32 {
	ret := int32(b.data[i]) & int32(an)
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

func (b *ImBuffer) writeMsgSize(x int32) {
	if x > 0 {
		if x < 0x80 {
			b.writeByte(byte(x))
		} else if x < 0x4000 {
			b.writeInt(x | 0x8000)
		} else if x < 0x200000 {
			b.writeInt(x | 0xc00000)
		} else if x < 0x10000000 {
			b.writeInt(int32(int(x) | 0xe0000000))
		}
	}

	/**
	if (x >= 0) {
            if (x < 0x80) {
                byteBuf.writeByte(x);
            } else if (x < 0x4000) {
                byteBuf.writeShort(x | 0x8000);
            } else if (x < 0x200000) {
                byteBuf.writeMedium(x | 0xc00000);
            } else if (x < 0x10000000) {
                byteBuf.writeInt(x | 0xe0000000);
            } else {
                throw new RuntimeException("exceed max unit");
            }
        }
	 */
}

func (b *ImBuffer) writeLong(x int64) {
	if x > 0 {
		if x < 0x80 {
			log("11 write ", x)
			b.data = append(b.data, byte(x))
			return
		} else if x < 0x4000 {
			log("22 write ", x)
			b.data = append(b.data, byte(x>>8|0x80), byte(x))
			return
		} else if x < 0x200000 {
			log("33 write ", x)
			b.data = append(b.data, byte(x>>16|0xc0), byte(x>>8), byte(x))
			return
		} else if x < 0x10000000 {
			log("44 write ", x)
			b.data = append(b.data, byte((x>>24)|0xe0), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x800000000 {
			b.data = append(b.data, byte((x>>32)|0xf0), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x40000000000 {
			b.data = append(b.data, byte((x>>40)|0xf8), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x200000000000 {
			b.data = append(b.data, byte((x>>48)|0xfc), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		} else if x < 0x100000000000000 {
			b.data = append(b.data, byte(0xfe), byte(x>>48), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
			return
		}
	}
	b.data = append(b.data, byte(0xff), byte(x>>56), byte(x>>48), byte(x>>40), byte(x>>32), byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
}

func (b *ImBuffer) readLong() int64 {

	//fmt.Println("0xff ", reflect.TypeOf(0xff))

	h := byte(b.data[b.begin] & 0xff)

	//fmt.Println(" h  ", reflect.TypeOf(h))

	if h < 0x80 {
		b.begin ++
		log(" 111 begin ", b.begin)
		return int64(h)
	} else if h < 0xc0 {
		x := byteMvInt64(h, 0x3f, 8) | byteMvInt64(b.data[b.begin+1], 0xff, 0)
		b.begin += 2
		log(" 222 begin ", b.begin)
		return x
	} else if h < 0xe0 {
		x := byteMvInt64(h, 0x1f, 16) | byteMvInt64(b.data[b.begin+1], 0xff, 8) | byteMvInt64(b.data[b.begin+2], 0xff, 0)
		b.begin += 3
		log(" 333 begin ", b.begin)
		return x
	} else if h < 0xf0 {
		x := byteMvInt64(h, 0x0f, 24) | byteMvInt64(b.data[b.begin+1], 0xff, 16) | byteMvInt64(b.data[b.begin+2], 0xff, 8) | byteMvInt64(b.data[b.begin+3], 0xff, 0)
		b.begin += 4
		log(" 444 begin ", b.begin)
		return x
	} else if h < 0xf8 {
		xl := byteMvInt64(b.data[b.begin+1], 0, 24) | byteMvInt64(b.data[b.begin+2], 0xff, 16) | byteMvInt64(b.data[b.begin+3], 0xff, 8) | byteMvInt64(b.data[b.begin+4], 0xff, 0)
		xh := byteMvInt64(h, 0x07, 0)
		b.begin += 5
		return xh<<32 | xl&0xffffffff
	} else if h < 0xfc {
		xl := byteMvInt64(b.data[b.begin+2], 0, 24) | byteMvInt64(b.data[b.begin+3], 0xff, 16) | byteMvInt64(b.data[b.begin+4], 0xff, 8) | byteMvInt64(b.data[b.begin+5], 0xff, 0)
		xh := byteMvInt64(h, 0x03, 8) | byteMvInt64(b.data[b.begin+1], 0xff, 0)
		b.begin += 6
		return xh<<32 | xl&0xffffffff
	} else if h < 0xfe {
		xl := byteMvInt64(b.data[b.begin+3], 0, 24) | byteMvInt64(b.data[b.begin+4], 0xff, 16) | byteMvInt64(b.data[b.begin+5], 0xff, 8) | byteMvInt64(b.data[b.begin+6], 0xff, 0)
		xh := byteMvInt64(h, 0x01, 16) | byteMvInt64(b.data[b.begin+1], 0xff, 8) | byteMvInt64(b.data[b.begin+2], 0xff, 0)
		b.begin += 7
		return xh<<32 | xl&0xffffffff
	} else if h < 0xff {
		xl := byteMvInt64(b.data[b.begin+4], 0, 24) | byteMvInt64(b.data[b.begin+5], 0xff, 16) | byteMvInt64(b.data[b.begin+6], 0xff, 8) | byteMvInt64(b.data[b.begin+7], 0xff, 0)
		xh := byteMvInt64(b.data[b.begin+1], 0xff, 16) | byteMvInt64(b.data[b.begin+2], 0xff, 8) | byteMvInt64(b.data[b.begin+3], 0xff, 0)
		b.begin += 8
		return xh<<32 | xl&0xffffffff
	} else {
		xl := byteMvInt64(b.data[b.begin+5], 0, 24) | byteMvInt64(b.data[b.begin+6], 0xff, 16) | byteMvInt64(b.data[b.begin+7], 0xff, 8) | byteMvInt64(b.data[b.begin+8], 0xff, 0)
		xh := byteMvInt64(b.data[b.begin+1], 0, 24) | byteMvInt64(b.data[b.begin+2], 0xff, 16) | byteMvInt64(b.data[b.begin+3], 0xff, 8) | byteMvInt64(b.data[b.begin+4], 0xff, 0)
		b.begin += 9
		return xh<<32 | xl&0xffffffff
	}
}

func (b *ImBuffer) writeInt(x int32) {
	if x >= 0 {
		if x < 0x80 {
			log("11 write ", x)
			b.data = append(b.data, byte(x))
			return
		} else if x < 0x4000 {
			log("22 write ", x)
			b.data = append(b.data, byte(x>>8|0x80), byte(x))
			return
		} else if x < 0x200000 {
			log("33 write ", x)
			b.data = append(b.data, byte(x>>16|0xc0), byte(x>>8), byte(x))
			return
		} else if x < 0x10000000 {
			log("44 write ", x)
			b.data = append(b.data, byte((x>>24)|0xe0), byte(x>>16), byte(x>>8), byte(x))
			return
		}
	}
	log("55 write ", x)

	b.data = append(b.data, 0xf0, byte(x>>24), byte(x>>16), byte(x>>8), byte(x))
}

func (b *ImBuffer) readInt() int32 {
	h := byte(b.data[b.begin] & 0xff)

	log(" read h = ", h, reflect.TypeOf(h))

	t := int32(h) << 24

	log(" read t = ", t, reflect.TypeOf(t))

	if h < 0x80 {
		b.begin ++
		log(" 111 begin ", b.begin)
		return int32(h)
	} else if h < 0xc0 {
		//x := int32(((h & 0x3f) << 8) | (b.data[b.begin+1] & 0xff))
		x := byteMvInt32(h, 0x3f, 8) | byteMvInt32(b.data[b.begin+1], 0xff, 0)
		b.begin += 2
		log(" 222 begin ", b.begin)
		return x
	} else if h < 0xe0 {
		//x := int32(((h & 0x1f) << 16) | ((b.data[b.begin+1] & 0xff) << 8) | (b.data[b.begin+2] & 0xff))
		x := byteMvInt32(h, 0x1f, 16) | byteMvInt32(b.data[b.begin+1], 0xff, 8) | byteMvInt32(b.data[b.begin+2], 0xff, 0)
		b.begin += 3
		log(" 333 begin ", b.begin)
		return x
	} else if h < 0xf0 {
		//x := int32(((h & 0x0f) << 24) | ((b.data[b.begin+1] & 0xff) << 16) | ((b.data[b.begin+2] & 0xff) << 8) | (b.data[b.begin+3] & 0xff))
		x := byteMvInt32(h, 0x0f, 24) | byteMvInt32(b.data[b.begin+1], 0xff, 16) | byteMvInt32(b.data[b.begin+2], 0xff, 8) | byteMvInt32(b.data[b.begin+3], 0xff, 0)
		b.begin += 4
		log(" 444 begin ", b.begin)
		return x
	} else {
		//x := int32(((b.data[b.begin+1] & 0xff) << 24) | ((b.data[b.begin+2] & 0xff) << 16) | ((b.data[b.begin+3] & 0xff) << 8) | (b.data[b.begin+4] & 0xff))
		x := byteMvInt32(b.data[b.begin+1], 0xff, 24) | byteMvInt32(b.data[b.begin+2], 0xff, 16) | byteMvInt32(b.data[b.begin+3], 0xff, 8) | byteMvInt32(b.data[b.begin+4], 0xff, 0)
		b.begin += 5
		log(" 555 begin ", b.begin)
		return x
	}
}

func (b *ImBuffer) writeBool(x bool) {
	if x {
		b.data = append(b.data, byte(1))
	} else {
		b.data = append(b.data, byte(0))
	}
}

func (b *ImBuffer) readBool() bool {
	x := b.data[b.begin]
	b.begin ++
	if x > 0 {
		return true
	} else {
		return false
	}
}

func (b *ImBuffer) writeByte(x byte) {
	b.data = append(b.data, x)
}

func (b *ImBuffer) readByte() byte {
	x := b.data[b.begin]
	return x
}

func (b *ImBuffer) writeBytes(x []byte) {
	b.writeInt(int32(len(x)))
	b.data = append(b.data, x...)
}

func (b *ImBuffer) readBytes() []byte {
	length := b.readInt()
	bytes := b.data[b.begin:length]
	b.begin += length
	return bytes
}

func (b *ImBuffer) writeFloat(x float32) {
	bits := math.Float32bits(x)
	b.data = append(b.data, byte(bits&0xff))
	b.data = append(b.data, byte((bits>>8)&0xff))
	b.data = append(b.data, byte((bits>>16)&0xff))
	b.data = append(b.data, byte((bits>>24)&0xff))
}

func (b *ImBuffer) readFloat() float32 {
	i1 := byteMvUInt32(b.data[b.begin], 0xff, 0)
	i2 := byteMvUInt32(b.data[b.begin+1], 0xff, 8)
	i3 := byteMvUInt32(b.data[b.begin+2], 0xff, 16)
	i4 := byteMvUInt32(b.data[b.begin+3], 0xff, 24)
	b.begin += 4
	return math.Float32frombits(i1 + i2 + i3 + i4)
}

func (b *ImBuffer) writeDouble(x float64) {
	bits := math.Float64bits(x)
	b.data = append(b.data, byte(bits))
	b.data = append(b.data, byte(bits>>8))
	b.data = append(b.data, byte(bits>>16))
	b.data = append(b.data, byte(bits>>24))
	b.data = append(b.data, byte(bits>>32))
	b.data = append(b.data, byte(bits>>40))
	b.data = append(b.data, byte(bits>>48))
	b.data = append(b.data, byte(bits>>56))
}

func (b *ImBuffer) readDouble() float64 {

	ui1 := byteMvUInt64(b.data[b.begin], 0xff, 0)
	ui2 := byteMvUInt64(b.data[b.begin+1], 0xff, 8)
	ui3 := byteMvUInt64(b.data[b.begin+2], 0xff, 16)
	ui4 := byteMvUInt64(b.data[b.begin+3], 0xff, 24)
	ui5 := byteMvUInt64(b.data[b.begin+4], 0xff, 32)
	ui6 := byteMvUInt64(b.data[b.begin+5], 0xff, 40)
	ui7 := byteMvUInt64(b.data[b.begin+6], 0xff, 48)
	ui8 := byteMvUInt64(b.data[b.begin+7], 0xff, 56)

	b.begin += 8
	return math.Float64frombits(ui1 + ui2 + ui3 + ui4 + ui5 + ui6 + ui7 + ui8)
}

func (b *ImBuffer) writeString(x string) {
	length := len(x)
	b.writeInt(int32(length))
	strBytes := []byte(x)
	b.data = append(b.data, strBytes...)

}

func (b *ImBuffer) readString() string {
	length := b.readInt()
	bytes := b.data[b.begin : b.begin+length]
	b.begin += length
	return string(bytes[:])
}

func (b *ImBuffer) writeBuffer(src *ImBuffer) {
	srcLength := len(src.data)
	b.writeInt(int32(srcLength))
	b.writeBytes(src.data)
}

func (b *ImBuffer) readBuffer() *ImBuffer {
	length := b.readInt()
	ret := ImBuffer{}
	ret.data = append(ret.data, b.data[b.begin:b.begin+length]...)
	//copy(ret.data, )
	b.begin += length
	return &ret
}

func (b *ImBuffer) writeShort(x int16) {

}

func (b *ImBuffer) readShort() int16 {
	return -1
}

func (b *ImBuffer) writeSInt(x int32) {

}

func (b *ImBuffer) readSInt() int32 {
	return -1
}

func (b *ImBuffer) writeSLong(x int64) {

}
func (b *ImBuffer) readSLong() int64 {
	return -1
}

func runTest() {

	myBuffer := ImBuffer{}

	//myBuffer.data = append(myBuffer.data, 12312)

	myBuffer.writeInt(7)
	myBuffer.writeInt(222)
	myBuffer.writeInt(33333)
	myBuffer.writeInt(4444444)
	myBuffer.writeInt(555555555)

	myBuffer.writeInt(-7)
	myBuffer.writeInt(-222)
	myBuffer.writeInt(-33333)
	myBuffer.writeInt(-4444444)
	myBuffer.writeInt(-555555555)

	fmt.Println(myBuffer, len(myBuffer.data), cap(myBuffer.data))

	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	/////

	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())
	fmt.Println(myBuffer.readInt())

	fmt.Println(myBuffer, len(myBuffer.data), cap(myBuffer.data))

	longBuff := ImBuffer{}

	longBuff.writeLong(7)
	longBuff.writeLong(222)
	longBuff.writeLong(33333)
	longBuff.writeLong(4444444)
	longBuff.writeLong(555555555)
	longBuff.writeLong(6666666666)
	longBuff.writeLong(77777777777)
	longBuff.writeLong(888888888888)
	longBuff.writeLong(9999999999999)
	longBuff.writeLong(19999999999999)
	longBuff.writeLong(129999999999999)
	longBuff.writeLong(1239999999999999)
	longBuff.writeLong(12349999999999999)
	longBuff.writeLong(123459999999999999)
	longBuff.writeLong(1234569999999999999)
	//longBuff.writeLong(12345679999999999999)
	//longBuff.writeLong(123456789999999999999)
	//longBuff.writeLong(1234567899999999999999)
	//longBuff.writeLong(12345678909999999999999)

	longBuff.writeLong(-7)
	longBuff.writeLong(-222)
	longBuff.writeLong(-33333)
	longBuff.writeLong(-4444444)
	longBuff.writeLong(-555555555)
	longBuff.writeLong(-6666666666)
	longBuff.writeLong(-77777777777)
	longBuff.writeLong(-888888888888)
	longBuff.writeLong(-9999999999999)
	longBuff.writeLong(-19999999999999)
	longBuff.writeLong(-129999999999999)
	longBuff.writeLong(-1239999999999999)
	longBuff.writeLong(-12349999999999999)
	longBuff.writeLong(-123459999999999999)
	longBuff.writeLong(-1234569999999999999)

	fmt.Println(longBuff, len(longBuff.data), cap(longBuff.data))

	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())

	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())
	//fmt.Println(longBuff.readLong())

	fmt.Println(longBuff, len(longBuff.data), cap(longBuff.data))

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
