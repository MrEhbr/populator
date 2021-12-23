package postgres

// Option defines options for Engine.
type Option func(*Engine)

func DisableForeignKeyCheck() Option {
	return func(e *Engine) {
		e.disableForeignKeyCheck = true
	}
}
