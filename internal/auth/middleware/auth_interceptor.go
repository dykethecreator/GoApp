package middleware

import (
	"context"
	"strings"

	appjwt "github.com/dykethecreator/GoApp/pkg/jwt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// userIDKeyType is an unexported type for context keys defined in this package.
// This prevents collisions with context keys defined in other packages.
type userIDKeyType struct{}

var userIDKey = userIDKeyType{}

// UserIDFromContext returns the authenticated user ID from context, if present.
func UserIDFromContext(ctx context.Context) (string, bool) {
	v := ctx.Value(userIDKey)
	if v == nil {
		return "", false
	}
	id, ok := v.(string)
	return id, ok
}

// UnaryAuthInterceptor returns a grpc.UnaryServerInterceptor that validates
// incoming requests using the provided TokenManager. On success, it injects
// the user ID into the context for downstream handlers.
func UnaryAuthInterceptor(tm *appjwt.TokenManager) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// Exempt AuthService methods by default (public endpoints)
		if strings.HasPrefix(info.FullMethod, "/auth.AuthService/") {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.Unauthenticated, "missing metadata")
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, status.Error(codes.Unauthenticated, "missing authorization header")
		}

		auth := authHeaders[0]
		if !strings.HasPrefix(strings.ToLower(auth), "bearer ") {
			return nil, status.Error(codes.Unauthenticated, "invalid authorization header")
		}
		token := strings.TrimSpace(auth[len("bearer "):])

		claims, err := tm.ValidateToken(token)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "invalid or expired token")
		}
		if claims.Type != appjwt.TokenTypeAccess {
			return nil, status.Error(codes.Unauthenticated, "token must be an access token")
		}

		// Inject user ID into context
		ctx = context.WithValue(ctx, userIDKey, claims.Subject)
		return handler(ctx, req)
	}
}
