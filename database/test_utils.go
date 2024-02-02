package database

func createVersion0Database() {
	Orm.Exec(CreateClipboardItemsTableQuery)
	Orm.Exec(CreateConfigsTableQuery)
}
