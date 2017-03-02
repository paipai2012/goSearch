/*****************************************************************************
 *  file name : NumberProfile.go
 *  author : paipai2012
 *  email  : paipai2012@gmail.com
 *
 *  file description : 文本类倒排索引类
 *
******************************************************************************/

package segment

import (
	"encoding/binary"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"utils"
)

//profile 正排索引，detail也是保存在这里
type profile struct {
	startDocId uint32
	curDocId   uint32
	isMomery   bool
	fieldType  uint64
	pflOffset  int64
	docLen     uint64
	fieldName  string `json:"fullname"` //完整的名字，用来进行文件操作的
	shift      uint8
	fake       bool
	pflNumber  []int64       `json:"-"`
	pflString  []string      `json:"-"`
	pflMmap    *utils.Mmap   `json:"-"`
	dtlMmap    *utils.Mmap   `json:"-"`
	Logger     *utils.Log4FE `json:"-"` //logger
}

func newEmptyFakeProfile(fieldType uint64, shift uint8, fieldName string, start uint32, docLen uint64, logger *utils.Log4FE) *profile {
	this := &profile{docLen: docLen, pflOffset: 0, shift: shift, isMomery: true, fieldType: fieldType, fieldName: fieldName, startDocId: start, curDocId: start, Logger: logger, pflNumber: nil, pflString: nil}
	this.pflString = make([]string, 0)
	this.pflNumber = make([]int64, 0)
	this.fake = true
	return this
}

// newEmptyProfile function description : 新建空的字符型正排索引
// params :
// return :
func newEmptyProfile(fieldType uint64, shift uint8, fieldName string, start uint32, logger *utils.Log4FE) *profile {
	this := &profile{fake: false, pflOffset: 0, shift: shift, isMomery: true, fieldType: fieldType, fieldName: fieldName, startDocId: start, curDocId: start, Logger: logger, pflNumber: nil, pflString: nil}
	this.pflString = make([]string, 0)
	this.pflNumber = make([]int64, 0)

	return this
}

// newProfileWithLocalFile function description : 新建空的字符型正排索引
// params :
// return :
func newProfileWithLocalFile(fieldType uint64, shift uint8, fullsegmentname string, pflMmap, dtlMmap *utils.Mmap, offset int64, docLen uint64, isMomery bool, logger *utils.Log4FE) *profile {

	this := &profile{fake: false, docLen: docLen, shift: shift, pflOffset: offset, isMomery: isMomery, fieldType: fieldType, pflMmap: pflMmap, dtlMmap: dtlMmap, Logger: logger}

	/*
	   	//打开正排文件
	   	pflFileName := fmt.Sprintf("%v.pfl", fullsegmentname)
	   	this.Logger.Info("[INFO] NumberProfile --> NewNumberProfileWithLocalFile :: Load NumberProfile pflFileName: %v", pflFileName)
	   	pflFd, err := os.OpenFile(pflFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	   	if err != nil {
	   		return &NumberProfile{isMomery: false, pflType: idxType, Logger: logger, pflContent: make([]int64, 0)}
	   	}
	   	defer pflFd.Close()

	   	os, offseterr := pflFd.Seek(offset, 0)
	   	if offseterr != nil || os != offset {
	   		this.Logger.Error("[ERROR] NumberProfile --> NewNumberProfileWithLocalFile :: Seek Offset Error %v", offseterr)
	   		return &NumberProfile{isMomery: false, pflType: idxType, Logger: logger, pflContent: make([]int64, 0)}
	   	}

	   	for index := 0; index < docLen; index++ {
	           var value int64
	   		var pfl utils.DetailInfo
	   		pfl.Len = 8//int(lens)
	   		pfl.Offset = os
	           err := binary.Read(pflFd, binary.LittleEndian, &value)
	           if err != nil {
	               this.Logger.Error("[ERROR] NumberProfile --> NewNumberProfileWithLocalFile :: Read PosFile error %v", err)
	               this.pflPostion = append(this.pflPostion, utils.DetailInfo{0, 0})
	               this.pflContent= append(this.pflContent,0xFFFFFFFF)
	               continue
	           }
	           this.pflContent=append(this.pflContent,value)
	   		this.pflPostion = append(this.pflPostion, pfl)

	   		offset := os + 8
	   		os, offseterr = pflFd.Seek(offset, 0)
	   		if offseterr != nil || os != offset {
	   			this.Logger.Error("[ERROR] NumberProfile --> NewNumberProfileWithLocalFile :: Seek Offset Error %v", offseterr)
	   			this.pflPostion = append(this.pflPostion, utils.DetailInfo{0, 0})
	               this.pflContent=append(this.pflContent,0xFFFFFFFF)
	   			continue
	   		}
	   	}
	*/
	//this.Logger.Info("[INFO] Load  Profile : %v.pfl", fullsegmentname)
	return this

}

// addDocument function description : 增加一个doc文档
// params : docid docid的编号
//			contentstr string  文档内容
// return : error 成功返回Nil，否则返回相应的错误信息
func (this *profile) addDocument(docid uint32, content interface{}) error {

	if docid != this.curDocId || this.isMomery == false {
		return errors.New("profile --> AddDocument :: Wrong DocId Number")
	}
	this.Logger.Trace("[TRACE] docid %v content %v", docid, content)

	vtype := reflect.TypeOf(content)
	var value int64 = 0xFFFFFFFF
	var ok error
	switch vtype.Name() {
	case "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":

		value, ok = strconv.ParseInt(fmt.Sprintf("%v", content), 0, 0)
		if ok != nil {
			value = 0xFFFFFFFF
		}
		this.pflNumber = append(this.pflNumber, value)
	case "float32":
		v, _ := content.(float32)
		value = int64(v * 100)
		this.pflNumber = append(this.pflNumber, value)
	case "float64":
		v, _ := content.(float64)
		value = int64(v * 100)
		this.pflNumber = append(this.pflNumber, value)
	case "string":
		if this.fieldType == utils.IDX_TYPE_NUMBER {
			value, ok = strconv.ParseInt(fmt.Sprintf("%v", content), 0, 0)
			if ok != nil {
				value = 0xFFFFFFFF
			}
			//this.Logger.Info("[INFO] value %v", value)
			this.pflNumber = append(this.pflNumber, value)
			//this.pflString = append(this.pflString, fmt.Sprintf("%v", content))
		} else if this.fieldType == utils.IDX_TYPE_DATE {

			value, _ = utils.IsDateTime(fmt.Sprintf("%v", content))
			this.pflNumber = append(this.pflNumber, value)

		} else {
			this.pflString = append(this.pflString, fmt.Sprintf("%v", content))
		}
	default:
		this.pflString = append(this.pflString, fmt.Sprintf("%v", content))
	}
	this.curDocId++
	return nil
}

// serialization function description : 序列化正排索引（标准操作）
// params :
// return : error 正确返回Nil，否则返回错误类型
func (this *profile) serialization(fullsegmentname string) (int64, int, error) {

	//打开正排文件
	pflFileName := fmt.Sprintf("%v.pfl", fullsegmentname)
	var pflFd *os.File
	var err error
	//this.Logger.Info("[INFO] NumberProfile --> Serialization :: Load NumberProfile pflFileName: %v", pflFileName)
	pflFd, err = os.OpenFile(pflFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer pflFd.Close()
	fi, _ := pflFd.Stat()
	offset := fi.Size()
	this.pflOffset = offset
	var lens int
	if this.fieldType == utils.IDX_TYPE_NUMBER || this.fieldType == utils.IDX_TYPE_DATE {
		valueBufer := make([]byte, 8)
		for _, info := range this.pflNumber {
			binary.LittleEndian.PutUint64(valueBufer, uint64(info))
			_, err = pflFd.Write(valueBufer)
			if err != nil {
				this.Logger.Error("[ERROR] NumberProfile --> Serialization :: Write Error %v", err)
			}
		}

		lens = len(this.pflNumber)
	} else {

		//打开dtl文件
		dtlFileName := fmt.Sprintf("%v.dtl", fullsegmentname)
		//this.Logger.Info("[INFO] StringProfile --> Serialization :: Load StringProfile dtlFileName: %v", dtlFileName)
		dtlFd, err := os.OpenFile(dtlFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return 0, 0, err
		}
		defer dtlFd.Close()
		fi, _ = dtlFd.Stat()
		dtloffset := fi.Size()

		lenBufer := make([]byte, 8)
		for _, info := range this.pflString {
			infolen := len(info)
			binary.LittleEndian.PutUint64(lenBufer, uint64(infolen))
			_, err := dtlFd.Write(lenBufer)
			cnt, err := dtlFd.WriteString(info)
			if err != nil || cnt != infolen {
				this.Logger.Error("[ERROR] StringProfile --> Serialization :: Write Error %v", err)
			}
			//存储offset
			//this.Logger.Info("[INFO] dtloffset %v,%v",dtloffset,infolen)
			binary.LittleEndian.PutUint64(lenBufer, uint64(dtloffset))
			_, err = pflFd.Write(lenBufer)
			if err != nil {
				this.Logger.Error("[ERROR] StringProfile --> Serialization :: Write Error %v", err)
			}
			dtloffset = dtloffset + int64(infolen) + 8

		}
		lens = len(this.pflString)

	}

	this.isMomery = false
	this.pflString = nil
	this.pflNumber = nil
	return offset, lens, err

}

// GetValue function description : 获取值
// params :
// return :
func (this *profile) getValue(pos uint32) (string, bool) {

	if this.fake {
		return "", true
	}

	if this.isMomery && pos < uint32(len(this.pflNumber)) {
		if this.fieldType == utils.IDX_TYPE_NUMBER {
			return fmt.Sprintf("%v", this.pflNumber[pos]), true
		} else if this.fieldType == utils.IDX_TYPE_DATE {
			return utils.FormatDateTime(this.pflNumber[pos])
		}
		return this.pflString[pos], true

	}
	if this.pflMmap == nil {
		return "", false
	}

	offset := this.pflOffset + int64(pos*8)
	if this.fieldType == utils.IDX_TYPE_NUMBER {
		//ov:=this.pflMmap.ReadInt64(offset)
		//if this.shift>0{
		//    return fmt.Sprintf("%v",float64(ov)/(math.Pow10(int(this.shift))) ), true
		//}
		return fmt.Sprintf("%v", this.pflMmap.ReadInt64(offset)), true
	} else if this.fieldType == utils.IDX_TYPE_DATE {
		return utils.FormatDateTime(this.pflMmap.ReadInt64(offset))

	}

	if this.dtlMmap == nil {
		return "", false
	}
	/*
	   if pos > 800000 {
	        this.Logger.Info("[INFO] Offset %v Pos %v  ",offset,pos)
	   }
	*/
	dtloffset := this.pflMmap.ReadInt64(offset)
	lens := this.dtlMmap.ReadInt64(dtloffset)
	return this.dtlMmap.ReadString(dtloffset+8, lens), true

}

func (this *profile) getIntValue(pos uint32) (int64, bool) {

	if this.fake {
		return 0xFFFFFFFF, true
	}

	if this.isMomery {
		if (this.fieldType == utils.IDX_TYPE_NUMBER || this.fieldType == utils.IDX_TYPE_DATE) &&
			pos < uint32(len(this.pflNumber)) {
			return this.pflNumber[pos], true
		}
		return 0xFFFFFFFF, false

	}
	if this.pflMmap == nil {
		return 0xFFFFFFFF, true
	}

	offset := this.pflOffset + int64(pos*8)
	if this.fieldType == utils.IDX_TYPE_NUMBER || this.fieldType == utils.IDX_TYPE_DATE {
		//ov:=this.pflMmap.ReadInt64(offset)
		//if this.shift>0{
		//    return fmt.Sprintf("%v",float64(ov)/(math.Pow10(int(this.shift))) ), true
		//}
		return this.pflMmap.ReadInt64(offset), true
	}

	return 0xFFFFFFFF, false
}

func (this *profile) filterNums(pos uint32, filtertype uint64, rangenum []int64) bool {
	var value int64
	if this.fake {
		return false
	}

	if this.fieldType == utils.IDX_TYPE_NUMBER {
		if this.isMomery {
			value = this.pflNumber[pos]
		} else if this.pflMmap == nil {
			return false
		}

		offset := this.pflOffset + int64(pos*8)
		value = this.pflMmap.ReadInt64(offset)

		switch filtertype {
		case utils.FILT_EQ:
			for _, start := range rangenum {
				if ok := (0xFFFFFFFF&value != 0xFFFFFFFF) && (value == start); ok {
					return true
				}
			}
			return false
		case utils.FILT_NOT:
			for _, start := range rangenum {
				if ok := (0xFFFFFFFF&value != 0xFFFFFFFF) && (value == start); ok {
					return false
				}
			}
			return true

		default:
			return false
		}

	}

	return false
}

// Filter function description : 过滤
// params :
// return :
func (this *profile) filter(pos uint32, filtertype uint64, start, end int64, str string) bool {

	var value int64
	if /*(this.fieldType != utils.IDX_TYPE_NUMBER && this.fieldType != utils.IDX_TYPE_DATE) ||*/ this.fake {
		return false
	}

	if this.fieldType == utils.IDX_TYPE_NUMBER {
		if this.isMomery {
			value = this.pflNumber[pos]
		} else if this.pflMmap == nil {
			return false
		}

		offset := this.pflOffset + int64(pos*8)
		value = this.pflMmap.ReadInt64(offset)

		switch filtertype {
		case utils.FILT_EQ:

			return (0xFFFFFFFF&value != 0xFFFFFFFF) && (value == start)
		case utils.FILT_OVER:
			return (0xFFFFFFFF&value != 0xFFFFFFFF) && (value >= start)
		case utils.FILT_LESS:
			return (0xFFFFFFFF&value != 0xFFFFFFFF) && (value <= start)
		case utils.FILT_RANGE:
			return (0xFFFFFFFF&value != 0xFFFFFFFF) && (value >= start && value <= end)
		case utils.FILT_NOT:
			return (0xFFFFFFFF&value != 0xFFFFFFFF) && (value != start)
		default:
			return false
		}
	}

	if this.fieldType == utils.IDX_TYPE_STRING_SINGLE {
		vl := strings.Split(str, ",")
		switch filtertype {

		case utils.FILT_STR_PREFIX:
			if vstr, ok := this.getValue(pos); ok {

				for _, v := range vl {
					if strings.HasPrefix(vstr, v) {
						return true
					}
				}

			}
			return false
		case utils.FILT_STR_SUFFIX:
			if vstr, ok := this.getValue(pos); ok {
				for _, v := range vl {
					if strings.HasSuffix(vstr, v) {
						return true
					}
				}
			}
			return false
		case utils.FILT_STR_RANGE:
			if vstr, ok := this.getValue(pos); ok {
				for _, v := range vl {
					if !strings.Contains(vstr, v) {
						return false
					}
				}
				return true
			}
			return false
		case utils.FILT_STR_ALL:

			if vstr, ok := this.getValue(pos); ok {
				for _, v := range vl {
					if vstr == v {
						return true
					}
				}
			}
			return false
		default:
			return false
		}

	}

	return false

}

// destroy function description : 销毁
// params :
// return :
func (this *profile) destroy() error {
	this.pflNumber = nil
	this.pflString = nil
	return nil
}

func (this *profile) setPflMmap(mmap *utils.Mmap) {
	this.pflMmap = mmap
}

func (this *profile) setDtlMmap(mmap *utils.Mmap) {
	this.dtlMmap = mmap
}

func (this *profile) updateDocument(docid uint32, content interface{}) error {

	if this.fieldType != utils.IDX_TYPE_NUMBER || this.fieldType != utils.IDX_TYPE_DATE {
		return errors.New("not support")
	}

	vtype := reflect.TypeOf(content)
	var value int64 = 0xFFFFFFFF
	switch vtype.Name() {
	case "string", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64":
		var ok error
		if this.fieldType == utils.IDX_TYPE_DATE {
			value, _ = utils.IsDateTime(fmt.Sprintf("%v", content))
		}
		value, ok = strconv.ParseInt(fmt.Sprintf("%v", content), 0, 0)
		if ok != nil {
			value = 0xFFFFFFFF
		}

	case "float32":
		v, _ := content.(float32)
		value = int64(v * 100)
	case "float64":
		v, _ := content.(float64)
		value = int64(v * 100)
	default:
		value = 0xFFFFFFFF
	}
	if this.isMomery == true {
		this.pflNumber[docid-this.startDocId] = value
	} else {
		offset := this.pflOffset + int64((docid-this.startDocId)*8)
		this.pflMmap.WriteInt64(offset, value)
	}
	return nil
}

func (this *profile) mergeProfiles(srclist []*profile, fullsegmentname string) (int64, int, error) {

	//this.Logger.Info("[INFO] mergeProfiles  %v",fullsegmentname )
	//if this.startDocId != 0 {
	//    this.Logger.Error("[ERROR] DocId Wrong %v",this.startDocId)
	//    return 0,0,errors.New("DocId Wrong")
	//}
	//打开正排文件
	pflFileName := fmt.Sprintf("%v.pfl", fullsegmentname)
	var pflFd *os.File
	var err error
	//this.Logger.Info("[INFO] NumberProfile --> Serialization :: Load NumberProfile pflFileName: %v", pflFileName)
	pflFd, err = os.OpenFile(pflFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return 0, 0, err
	}
	defer pflFd.Close()
	fi, _ := pflFd.Stat()
	offset := fi.Size()
	this.pflOffset = offset
	var lens int
	if this.fieldType == utils.IDX_TYPE_NUMBER || this.fieldType == utils.IDX_TYPE_DATE {
		valueBufer := make([]byte, 8)
		for _, src := range srclist {
			for i := uint32(0); i < uint32(src.docLen); i++ {
				val, _ := src.getIntValue(i)
				binary.LittleEndian.PutUint64(valueBufer, uint64(val))
				_, err = pflFd.Write(valueBufer)
				if err != nil {
					this.Logger.Error("[ERROR] NumberProfile --> Serialization :: Write Error %v", err)
				}
				this.curDocId++
			}
		}

		lens = int(this.curDocId - this.startDocId)
	} else {

		//打开dtl文件
		dtlFileName := fmt.Sprintf("%v.dtl", fullsegmentname)
		//this.Logger.Info("[INFO] StringProfile --> Serialization :: Load StringProfile dtlFileName: %v", dtlFileName)
		dtlFd, err := os.OpenFile(dtlFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return 0, 0, err
		}
		defer dtlFd.Close()
		fi, _ = dtlFd.Stat()
		dtloffset := fi.Size()

		lenBufer := make([]byte, 8)
		for _, src := range srclist {
			for i := uint32(0); i < uint32(src.docLen); i++ {
				info, _ := src.getValue(i)
				infolen := len(info)
				binary.LittleEndian.PutUint64(lenBufer, uint64(infolen))
				_, err := dtlFd.Write(lenBufer)
				cnt, err := dtlFd.WriteString(info)
				if err != nil || cnt != infolen {
					this.Logger.Error("[ERROR] StringProfile --> Serialization :: Write Error %v", err)
				}
				//存储offset
				//this.Logger.Info("[INFO] dtloffset %v,%v",dtloffset,infolen)
				binary.LittleEndian.PutUint64(lenBufer, uint64(dtloffset))
				_, err = pflFd.Write(lenBufer)
				if err != nil {
					this.Logger.Error("[ERROR] StringProfile --> Serialization :: Write Error %v", err)
				}
				dtloffset = dtloffset + int64(infolen) + 8
				this.curDocId++
			}
		}

		lens = int(this.curDocId - this.startDocId)

	}
	this.isMomery = false
	this.pflString = nil
	this.pflNumber = nil
	return offset, lens, nil

}
