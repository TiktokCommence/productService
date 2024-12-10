package service

var defaultPageSize uint32 = 20
var defaultListOptions = ListOptions{
	PageSize: defaultPageSize,
}

type ListOptions struct {
	PageSize uint32
	Category *string
}
type ListOption func(options *ListOptions)

func NewListOptions(opt ...ListOption) ListOptions {
	for _, o := range opt {
		o(&defaultListOptions)
	}
	return defaultListOptions
}

func WithCategory(category *string) ListOption {
	return func(options *ListOptions) {
		options.Category = category
	}
}
