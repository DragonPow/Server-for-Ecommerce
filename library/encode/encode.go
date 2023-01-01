package encode

type HttpEncoder interface {
	Marshal(v any) ([]byte, error)
}
