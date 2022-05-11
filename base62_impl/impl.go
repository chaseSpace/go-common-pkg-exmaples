package base62_impl

import "strconv"

const CODE62 = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var mapTable = map[byte]uint8{'0': 0, '1': 1, '2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7, '8': 8, '9': 9, 'a': 10,
	'b': 11, 'c': 12, 'd': 13, 'e': 14, 'f': 15, 'g': 16, 'h': 17, 'i': 18, 'j': 19, 'k': 20, 'l': 21, 'm': 22,
	'n': 23, 'o': 24, 'p': 25, 'q': 26, 'r': 27, 's': 28, 't': 29, 'u': 30, 'v': 31, 'w': 32, 'x': 33, 'y': 34,
	'z': 35, 'A': 36, 'B': 37, 'C': 38, 'D': 39, 'E': 40, 'F': 41, 'G': 42, 'H': 43, 'I': 44, 'J': 45, 'K': 46,
	'L': 47, 'M': 48, 'N': 49, 'O': 50, 'P': 51, 'Q': 52, 'R': 53, 'S': 54, 'T': 55, 'U': 56, 'V': 57, 'W': 58,
	'X': 59, 'Y': 60, 'Z': 61}

type InvalidInputErr int64

func (e InvalidInputErr) Error() string {
	return "illegal base62 data at input byte " + strconv.Itoa(int(e))
}

type InvalidInputLenErr int64

func (e InvalidInputLenErr) Error() string {
	return "invalid base62 data length " + strconv.Itoa(int(e))
}

func b62encode(input []byte) (output []byte) {
	if len(input) == 0 {
		return
	}

	_5bit := uint(0)
	_tmp8bytes := make([]byte, 8)

	idx := 0
	n := len(input) / 5 * 5 // 得到5的倍数，方便下面for循环操作

	for idx < n {
		// 将这个5byte合并到一个8byte的数中
		mergedVal := uint(input[idx])<<32 |
			uint(input[idx+1])<<24 |
			uint(input[idx+2])<<16 |
			uint(input[idx+3])<<8 |
			uint(input[idx+4])
		// 然后将8组5bit数据转换为base62字符
		for i := 0; i < 8; i++ {
			shift := 40 - 5*(i+1)
			if shift < 0 {
				_5bit = mergedVal << uint(^(shift - 1)) & 0x1f // ^(shift - 1) 是取绝对值
			} else {
				_5bit = mergedVal >> uint(shift) & 0x1f
			}
			_tmp8bytes[i] = toBase62Char(_5bit)
		}
		output = append(output, _tmp8bytes...)
		/* 下面是上述for循环的展开形式 */
		//_5bit = mergedVal >> (32 + 3) // 取出 0-4 bit (第一组5b)
		//_tmp8bytes[0] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> (24 + 6) & 0x1f // 取出 5-9 bit（第二组5b）
		//_tmp8bytes[1] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> (24 + 1) & 0x1f // 取出 10-14 bit（第三组5b）
		//_tmp8bytes[2] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> (16 + 4) & 0x1f // 第四组 15-19
		//_tmp8bytes[3] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> (8 + 7) & 0x1f // 第五组 20-24
		//_tmp8bytes[4] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> (8 + 2) & 0x1f // 第六组 25-29
		//_tmp8bytes[5] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> (0 + 5) & 0x1f // 第七组 30-34
		//_tmp8bytes[6] = toBase62Char(_5bit)
		//
		//_5bit = mergedVal >> 0 & 0x1f // 第八组 35-39
		//_tmp8bytes[7] = toBase62Char(_5bit)
		//output = append(output, _tmp8bytes...)

		idx += 5
	}
	// 处理非5倍数的情况
	remain := len(input) - n // remain 可能的值为 0|1|2|3|4
	if remain == 0 {
		return
	}
	groupNum := (remain*8 + 4) / 5 // 计算分为5bit一组的分组数，注意+4是为了不是5倍数的时候凑整
	mergedVal := uint(0)
	for i := 0; i < remain; i++ {
		mergedVal |= uint(input[idx+i]) << uint((remain-i-1)*8)
	}
	for i := 0; i < groupNum; i++ {
		shift := remain*8 - 5*(i+1)
		if shift < 0 {
			_5bit = mergedVal << uint(^(shift - 1)) & 0x1f // ^(shift - 1) 是取绝对值
		} else {
			_5bit = mergedVal >> uint(shift) & 0x1f
		}
		_tmp8bytes[i] = toBase62Char(_5bit)
	}
	output = append(output, _tmp8bytes[:groupNum]...)
	return
}

// 5bit转8bit，再从映射表中找到对应字符
func toBase62Char(_5bit uint) byte {
	_first7 := _5bit & 0xfe //  再取前7位
	_last1 := _5bit & 0x1   // 再取最后1位
	return CODE62[_first7<<1|_last1]
}

func b62decode(encodes []byte) (origin []byte, err error) {
	if len(encodes) == 0 {
		return
	}

	_tmp5bytes := make([]byte, 5)

	idx := 0
	n := len(encodes) / 8 * 8 // 得到8的倍数，方便for循环操作
	for idx < n {
		mergedVal := uint(0)
		// 每次操作8byte，从每个byte中取出5bit，存到一个8byte数中
		for i := 0; i < 8; i++ {
			num, ok := mapTable[encodes[idx+i]]
			if !ok {
				return nil, InvalidInputErr(idx + i)
			}
			mergedVal |= uint(parseOrigin5bits(num)) << uint((8-i-1)*5)
		}
		// mergedVal(64bit) = 24*0 + 40*originBit
		// 接下来把后面40bit 拆分为5byte
		for i := 0; i < 5; i++ {
			_tmp5bytes[i] = byte(mergedVal >> uint((5-i-1)*8) & 0xff)
		}
		origin = append(origin, _tmp5bytes...)
		idx += 8
	}

	// 处理非8倍数的情况
	remain := len(encodes) - n // remain 可能的值为 2,4,5,7
	if remain == 0 {
		return
	}
	switch remain {
	case 2, 4, 5, 7:
	default:
		return nil, InvalidInputLenErr(len(encodes))
	}
	bytesLen := (remain * 5) / 8 // 计算这部分原数据的长度，可能的值是 1,2,3,4
	mergedVal := uint(0)
	for i := 0; i < remain; i++ {
		num, ok := mapTable[encodes[idx+i]]
		if !ok {
			return nil, InvalidInputErr(idx + i)
		}
		mergedVal |= uint(parseOrigin5bits(num)) << uint((remain-i-1)*5)
	}
	for i := 0; i < bytesLen; i++ {
		shift := remain*5 - 8*(i+1)
		_tmp5bytes[i] = byte(mergedVal >> uint(shift) & 0xff)
	}
	origin = append(origin, _tmp5bytes[:bytesLen]...)
	return
}

// 0011 1101 => 0001 1111
func parseOrigin5bits(_8bit byte) byte {
	_first4b := _8bit >> 1
	_last1b := _8bit & 1
	return _first4b | _last1b
}
