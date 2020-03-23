package app

import (
	"github.com/shohrukh56/bookFileService/pkg/core/token"
	"github.com/shohrukh56/bookFileService/pkg/middleware/authenticated"
	"github.com/shohrukh56/bookFileService/pkg/middleware/jwt"
	"github.com/shohrukh56/bookFileService/pkg/middleware/logger"
	"reflect"
)

func (s *Server) InitRoutes() {

	s.router.GET(
		"/",
		s.handleIndex(),
	)
	s.router.POST(
		"/save",
		s.handleSaveFiles(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("Save."),
	)
	s.router.GET(
		"/media/{id}",
		s.handleGetFile(),
		authenticated.Authenticated(jwt.IsContextNonEmpty),
		jwt.JWT(reflect.TypeOf((*token.Payload)(nil)).Elem(), s.secret),
		logger.Logger("Media."),
		)

}
