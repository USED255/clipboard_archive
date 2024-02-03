package fts

func Open() error {
	err = connect()
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	sqlDB, err := orm.DB()
	if err != nil {
		return err
	}
	sqlDB.Close()
	orm = nil
	return nil
}
