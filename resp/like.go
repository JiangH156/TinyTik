package resp

import (
	"TinyTik/model"
)

type FavoriteList struct {
	Res    Response
	Videos *[]model.VideoList
}
