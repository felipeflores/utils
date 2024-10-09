package rest

type Request interface {
	Validate() error
}

type Response interface {
}
