package domain

type DBTransaction interface {
	Commit()
	Rollback()
}
