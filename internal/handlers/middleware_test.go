package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"storeSystem/internal/auth"
	"storeSystem/internal/models"
	"testing"
)

func TestAuthMiddleware_Success(t *testing.T) {
	h := &Handlers{}

	// создаём валидный токен
	token, _ := auth.GenerateToken(1, 1, "test")

	called := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		called = true

		claims, ok := GetUserClaimsFromContext(r.Context())
		if !ok {
			t.Fatal("no claims in context")
		}
		if claims.UserID != 1 {
			t.Fatal("wrong user id")
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: "token", Value: token})

	rr := httptest.NewRecorder()

	h.AuthMiddleware(next).ServeHTTP(rr, req)

	if !called {
		t.Fatal("next handler not called")
	}

	if rr.Code != http.StatusOK && rr.Code != 0 {
		t.Fatalf("unexpected status %d", rr.Code)
	}
}

func TestAuthMiddleware_NoCookie(t *testing.T) {
	h := &Handlers{}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("should not be called")
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	h.AuthMiddleware(next).ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestRequireAdmin_NoClaims(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	RequireAdmin()(next).ServeHTTP(rr, req)

	if nextCalled {
		t.Fatal("next should not be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 got %d", rr.Code)
	}
}

func TestRequireAdmin_Forbidden(t *testing.T) {
	nextCalled := false

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)

	req = req.WithContext(context.WithValue(
		req.Context(),
		UserClaimsKey,
		&models.Claims{
			UserID:   1,
			RoleID:   2, // not admin
			Username: "test",
		},
	))

	rr := httptest.NewRecorder()

	RequireAdmin()(next).ServeHTTP(rr, req)

	if nextCalled {
		t.Fatal("next should not be called")
	}

	if rr.Code != http.StatusForbidden {
		t.Fatalf("expected 403 got %d", rr.Code)
	}
}
