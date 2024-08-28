package interfaces

type IDbHandler interface {
	Execute(statement string)
	Query(statement string)
}
