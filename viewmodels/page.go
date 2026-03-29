package viewmodels

import "main/api"

type PageData struct {
	Artists      []api.Artist
	Title        string
	DisplayMode  string
	CurrentPage  int
	TotalPages   int
	HasPrevious  bool
	HasNext      bool
	PreviousPage int
	NextPage     int
}
