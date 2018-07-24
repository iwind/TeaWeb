package teaproxy

type RedirectError struct {
}

func (error *RedirectError) Error() string {
	return ""
}
