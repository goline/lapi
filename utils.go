package lapi

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func Must(errors ...error) {
	for _, err := range errors {
		PanicOnError(err)
	}
}
