package ndb

func InitDB() {
	prs := NewDBparams("root", "", "localhost", "3306", "transfusion")
	db := NewDBx(prs, DBType.MySQL)
	db.Connect()
}
