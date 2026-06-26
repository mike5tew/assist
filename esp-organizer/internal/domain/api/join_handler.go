package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---------- Models ----------

// Subscriber represents a mailing-list signup captured from /join.
type Subscriber struct {
	Email     string    `json:"email" bson:"email"`
	Source    string    `json:"source,omitempty" bson:"source,omitempty"` // which post / page brought them
	Role      string    `json:"role,omitempty" bson:"role,omitempty"`     // parent / teacher / clinician / other
	IP        string    `json:"-" bson:"ip,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
}

// ---------- MongoDB connection (lazy singleton, same pattern as analytics) ----------

var subscriberDB *mongo.Database

func getSubscriberDB() (*mongo.Database, error) {
	if subscriberDB != nil {
		return subscriberDB, nil
	}

	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		return nil, fmt.Errorf("MONGODB_URI not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("mongo connect: %w", err)
	}
	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("mongo ping: %w", err)
	}

	subscriberDB = client.Database("esp_subscribers")

	coll := subscriberDB.Collection("subscribers")

	// Unique index on email — prevents duplicates
	emailIdx := mongo.IndexModel{
		Keys:    bson.D{{Key: "email", Value: 1}},
		Options: options.Index().SetUnique(true),
	}
	if _, err := coll.Indexes().CreateOne(ctx, emailIdx); err != nil {
		log.Printf("subscribers: email unique index: %v", err)
	}

	// Source index for analytics
	sourceIdx := mongo.IndexModel{Keys: bson.D{{Key: "source", Value: 1}}}
	coll.Indexes().CreateOne(ctx, sourceIdx)

	log.Println("✅ Subscriber MongoDB connection established (database: esp_subscribers)")
	return subscriberDB, nil
}

// ---------- Validation ----------

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

// ---------- Handler ----------

// JoinHandler accepts an email signup from the /join landing page.
// POST /api/join
func JoinHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Email  string `json:"email"`
		Source string `json:"source"`
		Role   string `json:"role"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid request body"}`, http.StatusBadRequest)
		return
	}

	req.Email = strings.TrimSpace(strings.ToLower(req.Email))
	req.Source = strings.TrimSpace(req.Source)
	req.Role = strings.TrimSpace(req.Role)

	if req.Email == "" || !isValidEmail(req.Email) {
		http.Error(w, `{"error":"valid email is required"}`, http.StatusBadRequest)
		return
	}

	// Limit field lengths
	if len(req.Source) > 200 {
		req.Source = req.Source[:200]
	}
	if len(req.Role) > 50 {
		req.Role = req.Role[:50]
	}

	db, err := getSubscriberDB()
	if err != nil {
		log.Printf("join: DB error: %v", err)
		http.Error(w, `{"error":"service unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	sub := Subscriber{
		Email:     req.Email,
		Source:    req.Source,
		Role:      req.Role,
		IP:        realIP(r),
		CreatedAt: time.Now().UTC(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = db.Collection("subscribers").InsertOne(ctx, sub)
	if err != nil {
		// Duplicate key = already subscribed — that's fine
		if mongo.IsDuplicateKeyError(err) {
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]string{
				"status":  "already_subscribed",
				"message": "You're already on the list.",
			})
			return
		}
		log.Printf("join: insert error: %v", err)
		http.Error(w, `{"error":"failed to save subscription"}`, http.StatusInternalServerError)
		return
	}

	log.Printf("✅ New subscriber: %s (source: %s, role: %s)", req.Email, req.Source, req.Role)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "subscribed",
		"message": "Welcome. You're in.",
	})
}
