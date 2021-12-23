package populator

// Option defines options for Populator.
type Option func(*Populator)

func WithEngine(db Engine) Option {
	return func(p *Populator) {
		p.engine = db
	}
}

func WithParser(parser Parser) Option {
	return func(p *Populator) {
		p.parser = parser
	}
}
