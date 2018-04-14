package dbOperate

type sDB_Msg struct {
	ConnectDBErr  string
	InsertDataErr string
	DeleteDataErr string
	QueryDataErr  string
	UpdateDataErr string
}

var MsgDB sDB_Msg

func init() {
	MsgDB.ConnectDBErr = "错误,无法连接到指定数据库！"
	MsgDB.InsertDataErr = "错误,插入数据库操作失败！"
	MsgDB.DeleteDataErr = "错误,册除数据操作失败！"
	MsgDB.QueryDataErr = "错误,查询数据信息失败！"
	MsgDB.UpdateDataErr = "错误,更数据信息失败！"
}
