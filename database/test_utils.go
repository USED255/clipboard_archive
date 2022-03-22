package database

func createVersion0Database() {
	Orm.Exec(Query1)
	Orm.Exec(Query2)
}

func createVersion1Database() {
	Orm.Exec(Query1)
	Orm.Exec(Query2)
	Orm.Exec(Query3)

	Orm.Exec(`
		INSERT INTO "main"."configs" (
			"key", 
			"value"
		)
		 VALUES (
			"version", 
			"1.0.0"
		);
	`)
}

func createVersion2Database() {
	Orm.Exec(Query1)
	Orm.Exec(Query2)
	Orm.Exec(Query3)

	Orm.Exec(`
		INSERT INTO "main"."configs" (
			"key", 
			"value"
		)
		VALUES (
			"version", 
			"2.0.0"
		);
	`)

	Orm.Exec(Query4)
	Orm.Exec(Query5)
}
