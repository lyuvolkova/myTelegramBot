package handlers

type waterService interface {
	WriteColdData(text string, parTable [2]string, i int) string
	GetPrevData(i int) string
}
