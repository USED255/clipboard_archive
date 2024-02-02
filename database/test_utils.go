package database

func OpenNoDatabase() {
	connectDatabase("file::memory:?cache=shared")
}

func createVersion0Database() {
	Orm.Exec(CreateClipboardItemsTableQuery)
	Orm.Exec(CreateConfigsTableQuery)
}
