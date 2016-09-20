package gdb

import (
	"bufio"
	"os"

	"github.com/louch2010/gocache/core"
	"github.com/louch2010/gocache/log"
	"github.com/louch2010/goutil"
)

//解析gdb文件
func parseGDB(file *os.File) error {
	bfRd := bufio.NewReader(file)
	//魔数
	content, err := read(bfRd, LEN_GOCACHE)
	if err != nil {
		return err
	}
	if content != GOCACHE {
		return GDB_FILE_INVALID
	}
	//版本
	content, _ = read(bfRd, LEN_VERSION)
	if content != VERSION {
		return GDB_FILE_VERSION_ERROR
	}
	//库标识
	content, _ = read(bfRd, LEN_DATABASE)
	if content != DATABASE {
		log.Error("gdb文件格式错误，需要'database',但值为:", content)
		return GDB_FILE_FORMAT_ERROR
	}
	//库长
	databaseLenStr, _ := read(bfRd, LEN_KEY)
	databaseLen, _ := goutil.StringUtil().StrToInt(databaseLenStr)
	//遍历库
	for j := 0; j < databaseLen; j++ {
		//表名
		tableNameLenStr, _ := read(bfRd, LEN_KEY)
		tableNameLen, _ := goutil.StringUtil().StrToInt(tableNameLenStr)
		tableName, _ := read(bfRd, tableNameLen)
		table, err := core.Cache(tableName)
		if err != nil {
			log.Error("获取表失败！")
			return err
		}
		//键值对数
		keySizeStr, _ := read(bfRd, LEN_KEY)
		keySize, _ := goutil.StringUtil().StrToInt(keySizeStr)
		for i := 0; i < keySize; i++ {
			//数据类型
			dataType, _ := read(bfRd, LEN_DATATYPE)
			//过期时间

			//键
			keyLenStr, _ := read(bfRd, LEN_KEY)
			keyLen, _ := goutil.StringUtil().StrToInt(keyLenStr)
			key, _ := read(bfRd, keyLen)
			//值
			valueLenStr, _ := read(bfRd, LEN_VALUE)
			valueLen, _ := goutil.StringUtil().StrToInt(valueLenStr)
			value, _ := read(bfRd, valueLen)

			switch dataType {
			case TYPE_STRING:
				table.Set(key, value, 0, DATA_TYPE_STRING)
				break
			case TYPE_NUMBER:

				break
			}
		}
	}

	return nil
}

func read(bfRd *bufio.Reader, length int) (string, error) {
	buf := make([]byte, length)
	n, err := bfRd.Read(buf)
	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}
