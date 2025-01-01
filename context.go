package fortress

import "context"

type contextKey struct {
	string
}

func NewContextWithSessionToken(token string) context.Context {
	return context.WithValue(context.Background(), contextKey{"session_token"}, token)
}

func SessionTokenFromContext(c context.Context) string {
	token := c.Value(&contextKey{"session_token"})

	if token == nil {
		return ""
	}
	return token.(string)
}
