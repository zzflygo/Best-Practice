package week2

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

var ErrNoFound = errors.New("Dao:数据不存在")

const (
	ErrNoFoundCode = 40001 //sql.NoR
	ErrSystem      = 50001
)

type Data struct {
	Msg  string
	Code int
	Data interface{}
}

func Logic1() (err error) {
	//在logic层处理err
	err = Dao1("")
	if errors.Is(err, ErrNoFound) {
		//看业务是不是需要返回给前端.
		//不需要就用nil替代
		return nil
	}
	if err != nil {
		//出现的数据库查询错误,就让上层决定处理方式.
		return err
	}
	return nil
}

func Dao1(sqlstr string) (err error) {
	//dao层无论什么错误都往上抛
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNoFound
		}
		//其他错误返回
		return err
	}
	return nil
}

func IsNoRows(err error) bool {
	return strings.HasPrefix(err.Error(), fmt.Sprintf("%d", ErrNoFoundCode))
}

func Logic2() (err error) {
	err = Dao2("")
	if IsNoRows(err) {
		return err
	} else if err != nil {
		return err
	}
	return nil
}

func Dao2(sqlstr string) (err error) {
	if err == sql.ErrNoRows {
		return fmt.Errorf("%d , not found", ErrNoFoundCode)
	} else if err != nil {
		return fmt.Errorf("%d , not found", ErrSystem)
	}
	return nil

}

func main() {

}
