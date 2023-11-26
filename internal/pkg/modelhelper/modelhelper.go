package modelhelper

type ApiModel[V any] interface {
	ToAPI() V
}

func APISlice[V any, T any, PT interface {
	*T
	ApiModel[V]
}](s []T) []V {
	res := make([]V, len(s))
	for i := range s {
		res[i] = PT(&s[i]).ToAPI()
	}
	return res
}

func APISlicePtr[V any, T ApiModel[V]](s []T) []V {
	res := make([]V, len(s))
	for i := range s {
		res[i] = s[i].ToAPI()
	}
	return res
}
