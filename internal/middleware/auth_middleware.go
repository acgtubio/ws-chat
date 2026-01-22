package middleware

type AuthenticationMiddleware interface {
	NewAuthMiddlware() Middleware
}
