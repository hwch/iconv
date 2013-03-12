package better

import (
        "encoding/binary"
        "encoding/gob"
        "errors"
        "fmt"
        "os"
        "path"
        "runtime"
        "unsafe"
)

type CODING_IDX int

const (
        GBK2312_UNICODE_IDX  = iota
        GBK18030_UNICODE_IDX = iota
        GBK_UNICODE_IDX      = iota
        UNICODE_GBK_IDX      = iota
        UNICODE_GBK2312_IDX  = iota
        UNICODE_GBK18030_IDX = iota
        UTF8_UTF16_LE_IDX    = iota
        UTF16_LE_UTF8_IDX    = iota
        UTF8_UTF16_BE_IDX    = iota
        UTF16_BE_UTF8_IDX    = iota
        UTF8_GBK_IDX         = iota
        UTF8_GBK2312_IDX     = iota
        UTF8_GBK18030_IDX    = iota
        GBK_UTF8_IDX         = iota
        GBK2312_UTF8_IDX     = iota
        GBK18030_UTF8_IDX    = iota
)

type eleMent struct {
        filename string
        fn       func(map[uint64]uint64, []byte, []byte) (int, error)
}

var g_CodeMap = map[CODING_IDX]*eleMent{
        GBK18030_UNICODE_IDX: &eleMent{"Gbk2Unicode.db", convertGBKToUNICODE},
        GBK2312_UNICODE_IDX:  &eleMent{"Gbk2Unicode.db", convertGBKToUNICODE},
        GBK_UNICODE_IDX:      &eleMent{"Gbk2Unicode.db", convertGBKToUNICODE},
        UNICODE_GBK_IDX:      &eleMent{"Unicode2Gbk.db", convertUNICODEToGBK},
        UNICODE_GBK2312_IDX:  &eleMent{"Unicode2Gbk.db", convertUNICODEToGBK},
        UNICODE_GBK18030_IDX: &eleMent{"Unicode2Gbk.db", convertUNICODEToGBK},
        GBK18030_UTF8_IDX:    &eleMent{"Gbk2Unicode.db", convertGBKToUTF8},
        GBK2312_UTF8_IDX:     &eleMent{"Gbk2Unicode.db", convertGBKToUTF8},
        GBK_UTF8_IDX:         &eleMent{"Gbk2Unicode.db", convertGBKToUTF8},
        UTF8_GBK_IDX:         &eleMent{"Unicode2Gbk.db", convertUTF8ToGBK},
        UTF8_GBK18030_IDX:    &eleMent{"Unicode2Gbk.db", convertUTF8ToGBK},
        UTF8_GBK2312_IDX:     &eleMent{"Unicode2Gbk.db", convertUTF8ToGBK},
        UTF8_UTF16_LE_IDX:    &eleMent{"nil", convertUTF8ToUTF16LE},
        UTF16_LE_UTF8_IDX:    &eleMent{"nil", convertUTF16LEToUTF8},
        UTF8_UTF16_BE_IDX:    &eleMent{"nil", convertUTF8ToUTF16BE},
        UTF16_BE_UTF8_IDX:    &eleMent{"nil", convertUTF16BEToUTF8},
}

type Converter struct {
        isOpen  bool
        codeMap map[uint64]uint64
        // The first argument is input argument, that is to be translated.
        // The second argument is output argument, this has been translated. 
        // The memory of the second argument should be enough big,and it should be allocated by the caller of the function
        CodeConvertFunc func([]byte, []byte) (int, error)
}

func NewCoder(idx CODING_IDX) (*Converter, error) {
        var ele *eleMent
        ret := new(Converter)
        if v, ok := g_CodeMap[idx]; ok {
                ele = v
        } else {
                return nil, fmt.Errorf("Error: 未知编码格式\n")
        }
        if ele.filename != "nil" {
                _, filename, _, _ := runtime.Caller(0)
                ele.filename = path.Join(path.Dir(filename), ele.filename)
                ret.CodeConvertFunc = func(in []byte, out []byte) (int, error) {
                        return ele.fn(ret.codeMap, in, out)
                }
                ret.isOpen = true
        } else {

                ret.CodeConvertFunc = func(in []byte, out []byte) (int, error) {
                        return ele.fn(ret.codeMap, in, out)
                }
                return ret, nil
        }

        rf, err := os.Open(ele.filename)
        if err != nil {
                return nil, err
        }
        defer rf.Close()
        ret.codeMap = make(map[uint64]uint64)
        gb := gob.NewDecoder(rf)
        if err := gb.Decode(&ret.codeMap); err != nil {
                return nil, err
        }

        return ret, nil
}

func isLittleEndian() bool {
        var x int32 = 0x12345678
        p := unsafe.Pointer(&x)
        p1 := (*[4]byte)(p)
        return p1[0] == 0x78
}

func convertGBKToUNICODE(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        j := 0
        return j, nil
}

func convertUNICODEToGBK(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        j := 0
        return j, nil
}

func convertUTF16LEToUTF8(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        i := 0
        j := 0
        var tmpUnicode uint64
        fromLen := len(from)

        if fromLen%2 != 0 {
                return 0, fmt.Errorf("无效Unicode字符串")
        }
        if from[0] == 0xff {
                i += 2
        }
        for i < fromLen {
                tmpUnicode = uint64(binary.LittleEndian.Uint16(from[i:]))
                switch {
                case tmpUnicode < 0x00000080:
                        to[j] = byte(tmpUnicode)
                        j++
                case tmpUnicode < 0x00000800:
                        to[j] = 0xc0 | byte((tmpUnicode>>6)&0x1f)
                        to[j+1] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 2
                case tmpUnicode < 0x00010000:
                        to[j] = 0xe0 | byte((tmpUnicode>>12)&0x0f)
                        to[j+1] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+2] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 3
                case tmpUnicode < 0x00200000:
                        to[j] = 0xf0 | byte((tmpUnicode>>18)&0x07)
                        to[j+1] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+3] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 4
                case tmpUnicode < 0x04000000:
                        to[j] = 0xf8 | byte((tmpUnicode>>24)&0x03)
                        to[j+1] = 0x80 | byte((tmpUnicode>>18)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+3] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+4] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 5
                case tmpUnicode < 0x80000000:
                        to[j] = 0xfc | byte((tmpUnicode>>30)&0x01)
                        to[j+1] = 0x80 | byte((tmpUnicode>>24)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>18)&0x3f)
                        to[j+3] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+4] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+5] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 6
                default:
                        return 0, fmt.Errorf("非法字符[0x%x]", tmpUnicode)
                }
                i += 2
        }

        return j, nil
}

func convertUTF16BEToUTF8(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        i := 0
        j := 0
        var tmpUnicode uint64
        fromLen := len(from)

        if fromLen%2 != 0 {
                return 0, fmt.Errorf("无效Unicode字符串")
        }
        if from[0] == 0xfe {
                i += 2
        }
        for i < fromLen {
                tmpUnicode = uint64(binary.BigEndian.Uint16(from[i:]))
                switch {
                case tmpUnicode < 0x00000080:
                        to[j] = byte(tmpUnicode)
                        j++
                case tmpUnicode < 0x00000800:
                        to[j] = 0xc0 | byte((tmpUnicode>>6)&0x1f)
                        to[j+1] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 2
                case tmpUnicode < 0x00010000:
                        to[j] = 0xe0 | byte((tmpUnicode>>12)&0x0f)
                        to[j+1] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+2] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 3
                case tmpUnicode < 0x00200000:
                        to[j] = 0xf0 | byte((tmpUnicode>>18)&0x07)
                        to[j+1] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+3] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 4
                case tmpUnicode < 0x04000000:
                        to[j] = 0xf8 | byte((tmpUnicode>>24)&0x03)
                        to[j+1] = 0x80 | byte((tmpUnicode>>18)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+3] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+4] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 5
                case tmpUnicode < 0x80000000:
                        to[j] = 0xfc | byte((tmpUnicode>>30)&0x01)
                        to[j+1] = 0x80 | byte((tmpUnicode>>24)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>18)&0x3f)
                        to[j+3] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+4] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+5] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 6
                default:
                        return 0, fmt.Errorf("非法字符[0x%x]", tmpUnicode)
                }
                i += 2
        }

        return j, nil
}

//将GBK编码转换为UTF-8编码
func convertGBKToUTF8(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        i := 0
        var tmpGbk uint64
        var tmpUnicode uint64

        fromLen := len(from)
        for i < fromLen {
                if from[i]&0x80 == 0 { // ascii
                        i++
                } else {
                        if from[i+1] > 0x39 || from[i+1] < 0x30 {
                                i += 2
                        } else {
                                i += 4
                        }
                }
        }
        if i != fromLen {
                return 0, errors.New("非法GBK编码")
        }
        i = 0
        j := 0
        for i < fromLen {
                tmpGbk = 0x0
                if from[i]&0x80 == 0 { // ascii
                        tmpGbk = uint64(from[i])
                        i++
                } else {
                        if from[i+1] > 0x39 || from[i+1] < 0x30 {
                                tmpGbk = uint64(from[i+1])
                                tmpGbk |= uint64(from[i]) << 8
                                i += 2
                        } else {
                                tmpGbk = uint64(from[i+3])
                                tmpGbk |= uint64(from[i+2]) << 8
                                tmpGbk |= uint64(from[i+1]) << 16
                                tmpGbk |= uint64(from[i]) << 24
                                i += 4
                        }
                }
                if v, ok := tbl_map[tmpGbk]; ok {
                        tmpUnicode = v
                } else {
                        return 0, fmt.Errorf("未找到对应字符[0x%x]", tmpGbk)
                }
                switch {
                case tmpUnicode < 0x00000080:
                        to[j] = byte(tmpUnicode)
                        j++
                case tmpUnicode < 0x00000800:
                        to[j] = 0xc0 | byte((tmpUnicode>>6)&0x1f)
                        to[j+1] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 2
                case tmpUnicode < 0x00010000:
                        to[j] = 0xe0 | byte((tmpUnicode>>12)&0x0f)
                        to[j+1] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+2] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 3
                case tmpUnicode < 0x00200000:
                        to[j] = 0xf0 | byte((tmpUnicode>>18)&0x07)
                        to[j+1] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+3] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 4
                case tmpUnicode < 0x04000000:
                        to[j] = 0xf8 | byte((tmpUnicode>>24)&0x03)
                        to[j+1] = 0x80 | byte((tmpUnicode>>18)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+3] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+4] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 5
                case tmpUnicode < 0x80000000:
                        to[j] = 0xfc | byte((tmpUnicode>>30)&0x01)
                        to[j+1] = 0x80 | byte((tmpUnicode>>24)&0x3f)
                        to[j+2] = 0x80 | byte((tmpUnicode>>18)&0x3f)
                        to[j+3] = 0x80 | byte((tmpUnicode>>12)&0x3f)
                        to[j+4] = 0x80 | byte((tmpUnicode>>6)&0x3f)
                        to[j+5] = 0x80 | byte(tmpUnicode&0x3f)
                        j += 6
                default:
                        return 0, fmt.Errorf("非法字符[0x%x]", tmpUnicode)
                }
        }
        return j, nil
}

func convertUTF8ToUTF16LE(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        i := 0
        j := 0
        var tmpUnicode uint64
        fromLen := len(from)
        for i < fromLen {
                switch {
                case 0x80&from[i] == 0x0:
                        i++
                case 0xe0&from[i] == 0xc0:
                        i += 2
                case 0xf0&from[i] == 0xe0:
                        i += 3
                case 0xf8&from[i] == 0xf0:
                        i += 4
                case 0xfc&from[i] == 0xf8:
                        i += 5
                case 0xfe&from[i] == 0xfc:
                        i += 6
                default:
                        return 0, errors.New("无效UTF-8字符")
                }
        }
        if i != fromLen {
                return 0, errors.New("无效UTF-8字符串")
        }
        i = 0
        // -------------------------------------
        to[j] = 0xff
        to[j+1] = 0xfe
        j += 2
        // -------------------------------------
        for i < fromLen {
                tmpUnicode = 0x00
                switch {
                case 0x80&from[i] == 0x00:
                        tmpUnicode = uint64(from[i])
                        i++
                case 0xe0&from[i] == 0xc0:
                        tmpUnicode |= (uint64(from[i]&0x1f) << 6) | uint64(from[i+1]&0x3f)
                        i += 2
                case 0xf0&from[i] == 0xe0:
                        tmpUnicode |= (uint64(from[i]&0x0f) << 12) | (uint64(from[i+1]&0x3f) << 6)
                        tmpUnicode |= uint64(from[i+2] & 0x3f)
                        i += 3
                case 0xf8&from[i] == 0xf0:
                        tmpUnicode |= (uint64(from[i]&0x07) << 18) | (uint64(from[i+1]&0x3f) << 12)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 6) | uint64(from[i+3]&0x3f)
                        i += 4
                case 0xfc&from[i] == 0xf8:
                        tmpUnicode |= (uint64(from[i]&0x03) << 24) | (uint64(from[i+1]&0x3f) << 18)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 12) | (uint64(from[i+3]&0x3f) << 6)
                        tmpUnicode |= uint64(from[i+4] & 0x3f)
                        i += 5
                case 0xfe&from[i] == 0xfc:
                        tmpUnicode |= (uint64(from[i]&0x01) << 30) | (uint64(from[i+1]&0x3f) << 24)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 18) | (uint64(from[i+3]&0x3f) << 12)
                        tmpUnicode |= (uint64(from[i+4]&0x3f) << 6) | uint64(from[i+5]&0x3f)
                        i += 6
                default:
                        return 0, fmt.Errorf("无效UTF-8字符[0x%x]", from[i])
                }
                binary.LittleEndian.PutUint16(to[j:], uint16(tmpUnicode))
                j += 2
        }

        return j, nil
}

func convertUTF8ToUTF16BE(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        i := 0
        j := 0
        var tmpUnicode uint64
        fromLen := len(from)
        for i < fromLen {
                switch {
                case 0x80&from[i] == 0x0:
                        i++
                case 0xe0&from[i] == 0xc0:
                        i += 2
                case 0xf0&from[i] == 0xe0:
                        i += 3
                case 0xf8&from[i] == 0xf0:
                        i += 4
                case 0xfc&from[i] == 0xf8:
                        i += 5
                case 0xfe&from[i] == 0xfc:
                        i += 6
                default:
                        return 0, errors.New("无效UTF-8字符")
                }
        }
        if i != fromLen {
                return 0, errors.New("无效UTF-8字符串")
        }
        i = 0
        // -------------------------------------
        to[j] = 0xfe
        to[j+1] = 0xff
        j += 2
        // -------------------------------------
        for i < fromLen {
                tmpUnicode = 0x00
                switch {
                case 0x80&from[i] == 0x00:
                        tmpUnicode = uint64(from[i])
                        i++
                case 0xe0&from[i] == 0xc0:
                        tmpUnicode |= (uint64(from[i]&0x1f) << 6) | uint64(from[i+1]&0x3f)
                        i += 2
                case 0xf0&from[i] == 0xe0:
                        tmpUnicode |= (uint64(from[i]&0x0f) << 12) | (uint64(from[i+1]&0x3f) << 6)
                        tmpUnicode |= uint64(from[i+2] & 0x3f)
                        i += 3
                case 0xf8&from[i] == 0xf0:
                        tmpUnicode |= (uint64(from[i]&0x07) << 18) | (uint64(from[i+1]&0x3f) << 12)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 6) | uint64(from[i+3]&0x3f)
                        i += 4
                case 0xfc&from[i] == 0xf8:
                        tmpUnicode |= (uint64(from[i]&0x03) << 24) | (uint64(from[i+1]&0x3f) << 18)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 12) | (uint64(from[i+3]&0x3f) << 6)
                        tmpUnicode |= uint64(from[i+4] & 0x3f)
                        i += 5
                case 0xfe&from[i] == 0xfc:
                        tmpUnicode |= (uint64(from[i]&0x01) << 30) | (uint64(from[i+1]&0x3f) << 24)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 18) | (uint64(from[i+3]&0x3f) << 12)
                        tmpUnicode |= (uint64(from[i+4]&0x3f) << 6) | uint64(from[i+5]&0x3f)
                        i += 6
                default:
                        return 0, fmt.Errorf("无效UTF-8字符[0x%x]", from[i])
                }
                binary.BigEndian.PutUint16(to[j:], uint16(tmpUnicode))
                j += 2
        }

        return j, nil
}

//将UTF-8编码转换为GBK编码
func convertUTF8ToGBK(tbl_map map[uint64]uint64, from []byte, to []byte) (int, error) {
        var tmpUnicode uint64
        var tmpGbk uint64
        i := 0
        fromLen := len(from)
        for i < fromLen {
                switch {
                case 0x80&from[i] == 0x0:
                        i++
                case 0xe0&from[i] == 0xc0:
                        i += 2
                case 0xf0&from[i] == 0xe0:
                        i += 3
                case 0xf8&from[i] == 0xf0:
                        i += 4
                case 0xfc&from[i] == 0xf8:
                        i += 5
                case 0xfe&from[i] == 0xfc:
                        i += 6
                default:
                        return 0, errors.New("无效UTF-8字符")
                }
        }
        if i != fromLen {
                return 0, errors.New("无效长度")
        }
        i = 0
        j := 0
        for i < fromLen {
                tmpUnicode = 0x00
                switch {
                case 0x80&from[i] == 0x00:
                        tmpUnicode = uint64(from[i])
                        i++
                case 0xe0&from[i] == 0xc0:
                        tmpUnicode |= (uint64(from[i]&0x1f) << 6) | uint64(from[i+1]&0x3f)
                        i += 2
                case 0xf0&from[i] == 0xe0:
                        tmpUnicode |= (uint64(from[i]&0x0f) << 12) | (uint64(from[i+1]&0x3f) << 6)
                        tmpUnicode |= uint64(from[i+2] & 0x3f)
                        i += 3
                case 0xf8&from[i] == 0xf0:
                        tmpUnicode |= (uint64(from[i]&0x07) << 18) | (uint64(from[i+1]&0x3f) << 12)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 6) | uint64(from[i+3]&0x3f)
                        i += 4
                case 0xfc&from[i] == 0xf8:
                        tmpUnicode |= (uint64(from[i]&0x03) << 24) | (uint64(from[i+1]&0x3f) << 18)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 12) | (uint64(from[i+3]&0x3f) << 6)
                        tmpUnicode |= uint64(from[i+4] & 0x3f)
                        i += 5
                case 0xfe&from[i] == 0xfc:
                        tmpUnicode |= (uint64(from[i]&0x01) << 30) | (uint64(from[i+1]&0x3f) << 24)
                        tmpUnicode |= (uint64(from[i+2]&0x3f) << 18) | (uint64(from[i+3]&0x3f) << 12)
                        tmpUnicode |= (uint64(from[i+4]&0x3f) << 6) | uint64(from[i+5]&0x3f)
                        i += 6
                default:
                        return 0, fmt.Errorf("无效UTF-8字符[0x%x]", from[i])
                }
                if v, ok := tbl_map[tmpUnicode]; ok {
                        tmpGbk = v
                } else {
                        return 0, fmt.Errorf("未找到对应字符[0x%x]", tmpUnicode)
                }
                switch {
                case tmpGbk < 0x80:
                        to[j] = byte(tmpGbk)
                        j++
                case tmpGbk < 0x10000:
                        to[j] = byte(tmpGbk >> 8)
                        to[j+1] = byte(tmpGbk)
                        j += 2
                case tmpGbk < 0x100000000:
                        to[j] = byte(tmpGbk >> 24)
                        to[j+1] = byte(tmpGbk >> 16)
                        to[j+2] = byte(tmpGbk >> 8)
                        to[j+3] = byte(tmpGbk)
                        j += 4
                default:
                        return 0, fmt.Errorf("非法对应字符[0x%x]", tmpGbk)
                }
        }
        return j, nil
}
