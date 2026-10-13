package server

type Module interface {
	Init(
		app *App,
		ctr *Container,
	)
}
