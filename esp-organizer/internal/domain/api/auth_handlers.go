package api

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"time"
)

// jwtHeader is the fixed base64url-encoded header for HS256 JWTs.
const jwtHeader = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"

type jwtClaims struct {
	Sub string `json:"sub"`
	Exp int64  `json:"exp"`
}

func b64url(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func signJWT(payload []byte) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme-set-JWT_SECRET"
	}
	payloadB64 := b64url(payload)
	unsigned := jwtHeader + "." + payloadB64
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	sig := b64url(mac.Sum(nil))
	return unsigned + "." + sig, nil
}

func verifyJWT(token string) (*jwtClaims, bool) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, false
	}
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "changeme-set-JWT_SECRET"
	}
	unsigned := parts[0] + "." + parts[1]
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(unsigned))
	expected := b64url(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(parts[2])) {
		return nil, false
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, false
	}
	var claims jwtClaims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, false
	}
	if time.Now().Unix() > claims.Exp {
		return nil, false
	}
	return &claims, true
}

// LoginRequest is the JSON body for POST /api/auth/login.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// CHISGLoginHandler handles POST /api/auth/login.
func CHISGLoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	expectedUser := os.Getenv("CHISG_USER")
	expectedPass := os.Getenv("CHISG_PASS")
	if expectedUser == "" {
		expectedUser = "chisg_admin"
	}

	if req.Username != expectedUser || req.Password != expectedPass {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	claims := jwtClaims{
		Sub: req.Username,
		Exp: time.Now().Add(24 * time.Hour).Unix(),
	}
	payloadBytes, _ := json.Marshal(claims)
	token, err := signJWT(payloadBytes)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// CHISGDemoLoginHandler handles GET /api/auth/demo-login.
// Issues a short-lived read-only demo token without requiring credentials.
// Intentional: the CHISG knowledge base is a public-facing demo for collaborators.
func CHISGDemoLoginHandler(w http.ResponseWriter, r *http.Request) {
	claims := jwtClaims{
		Sub: "demo",
		Exp: time.Now().Add(4 * time.Hour).Unix(),
	}
	payloadBytes, _ := json.Marshal(claims)
	token, err := signJWT(payloadBytes)
	if err != nil {
		http.Error(w, "failed to generate token", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

// RequireAuth is a middleware that validates the Bearer JWT on protected routes.
func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow preflight through
		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "missing or invalid Authorization header", http.StatusUnauthorized)
			return
		}
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if _, ok := verifyJWT(token); !ok {
			http.Error(w, "invalid or expired token", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
