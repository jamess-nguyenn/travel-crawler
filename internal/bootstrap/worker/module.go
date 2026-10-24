package worker

type Module interface {
	Init(app *App, ctr *Container) error
}
