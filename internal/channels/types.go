package channels

type (
	Quote struct {
		Symbol string
		Price  int64
	}

	Params struct {
		Quotes []Quote
		Buffer int
	}

	Result struct {
		Received []Quote
	}
)
