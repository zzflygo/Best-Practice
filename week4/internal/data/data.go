package data

type Dao struct {
	Message string
}

func NewDao(sqlstr string) Dao {
	return Dao{Message: "hi~" + sqlstr}
}
