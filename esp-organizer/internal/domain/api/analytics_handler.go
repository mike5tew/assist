package api

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ---------- Models ----------

// PageView represents a single page-view event.
type PageView struct {
	Path      string    `json:"path" bson:"path"`
	Referrer  string    `json:"referrer,omitempty" bson:"referrer,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	UserAgent string    `json:"-" bson:"user_agent"`
	IP        string    `json:"-" bson:"ip"`
	VisitorID string    `json:"-" bson:"visitor_id"` // hashed fingerprint
	Timestamp time.Time `json:"ts" bson:"ts"`
	Country   string    `json:"-" bson:"country,omitempty"`
}

// AnalyticsEvent represents a CTA click or custom event.
type AnalyticsEvent struct {
	Name      string            `json:"name" bson:"name"`                     // e.g. "contact_open", "cta_click"
	Path      string            `json:"path" bson:"path"`                     // page where the event occurred
	Meta      map[string]string `json:"meta,omitempty" bson:"meta,omitempty"` // extra data
	UserAgent string            `json:"-" bson:"user_agent"`
	IP        string            `json:"-" bson:"ip"`
	VisitorID string            `json:"-" bson:"visitor_id"`
	Timestamp time.Time         `json:"ts" bson:"ts"`
}

// StatsResponse is the JSON returned by the /api/analytics/stats endpoint.
type StatsResponse struct {
	Period         string         `json:"period"`
	TotalViews     int64          `json:"total_views"`
	UniqueVisitors int64          `json:"unique_visitors"`
	TopPages       []PageStat     `json:"top_pages"`
	TopReferrers   []ReferrerStat `json:"top_referrers"`
	TopEvents      []EventStat    `json:"top_events"`
	ViewsByDay     []DayStat      `json:"views_by_day"`
}

type PageStat struct {
	Path  string `json:"path" bson:"_id"`
	Views int64  `json:"views" bson:"count"`
}

type ReferrerStat struct {
	Referrer string `json:"referrer" bson:"_id"`
	Count    int64  `json:"count" bson:"count"`
}

type EventStat struct {
	Name  string `json:"name" bson:"_id"`
	Count int64  `json:"count" bson:"count"`
}

type DayStat struct {
	Date  string `json:"date" bson:"_id"`
	Views int64  `json:"views" bson:"count"`
}

// ---------- MongoDB connection (lazy singleton) ----------

var analyticsDB *mongo.Database

func getAnalyticsDB() (*mongo.Database, error) {
	if analyticsDB != nil {
		return analyticsDB, nil
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

	analyticsDB = client.Database("esp_analytics")

	// Create TTL index: auto-delete documents older than 365 days
	viewsColl := analyticsDB.Collection("pageviews")
	eventsColl := analyticsDB.Collection("events")
	idxModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "ts", Value: 1}},
		Options: options.Index().SetExpireAfterSeconds(365 * 24 * 60 * 60),
	}
	if _, err := viewsColl.Indexes().CreateOne(ctx, idxModel); err != nil {
		log.Printf("analytics: TTL index (pageviews): %v", err)
	}
	if _, err := eventsColl.Indexes().CreateOne(ctx, idxModel); err != nil {
		log.Printf("analytics: TTL index (events): %v", err)
	}

	// Path index for aggregation speed
	pathIdx := mongo.IndexModel{Keys: bson.D{{Key: "path", Value: 1}}}
	viewsColl.Indexes().CreateOne(ctx, pathIdx)

	log.Println("✅ Analytics MongoDB connection established (database: esp_analytics)")
	return analyticsDB, nil
}

// ---------- Helpers ----------

// visitorFingerprint creates a privacy-friendly hash from IP + UA + daily salt.
// Changes each day so we can count daily uniques without persisting personal data.
func visitorFingerprint(ip, ua string) string {
	daySalt := time.Now().UTC().Format("2006-01-02")
	h := sha256.Sum256([]byte(ip + "|" + ua + "|" + daySalt))
	return hex.EncodeToString(h[:8]) // 16 hex chars, plenty for counting
}

func realIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

// ---------- Handlers ----------

// AnalyticsPageViewHandler records a page view.
// POST /api/analytics/pageview
func AnalyticsPageViewHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Path     string `json:"path"`
		Referrer string `json:"referrer"`
		Title    string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	if req.Path == "" {
		http.Error(w, `{"error":"path required"}`, http.StatusBadRequest)
		return
	}

	ip := realIP(r)
	ua := r.UserAgent()

	pv := PageView{
		Path:      req.Path,
		Referrer:  req.Referrer,
		Title:     req.Title,
		UserAgent: ua,
		IP:        ip,
		VisitorID: visitorFingerprint(ip, ua),
		Timestamp: time.Now().UTC(),
	}

	db, err := getAnalyticsDB()
	if err != nil {
		log.Printf("analytics: db error: %v", err)
		// Still return 200 — tracking should never block the user
		json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := db.Collection("pageviews").InsertOne(ctx, pv); err != nil {
		log.Printf("analytics: insert pageview: %v", err)
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// AnalyticsEventHandler records a custom event (CTA click, contact form open, etc.).
// POST /api/analytics/event
func AnalyticsEventHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var req struct {
		Name string            `json:"name"`
		Path string            `json:"path"`
		Meta map[string]string `json:"meta,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"invalid body"}`, http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, `{"error":"name required"}`, http.StatusBadRequest)
		return
	}

	ip := realIP(r)
	ua := r.UserAgent()

	evt := AnalyticsEvent{
		Name:      req.Name,
		Path:      req.Path,
		Meta:      req.Meta,
		UserAgent: ua,
		IP:        ip,
		VisitorID: visitorFingerprint(ip, ua),
		Timestamp: time.Now().UTC(),
	}

	db, err := getAnalyticsDB()
	if err != nil {
		log.Printf("analytics: db error: %v", err)
		json.NewEncoder(w).Encode(map[string]string{"status": "logged"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if _, err := db.Collection("events").InsertOne(ctx, evt); err != nil {
		log.Printf("analytics: insert event: %v", err)
	}

	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// AnalyticsStatsHandler returns aggregate stats for the dashboard.
// GET /api/analytics/stats?period=7d (default 7d, options: 24h, 7d, 30d, 90d)
// Protected by a simple bearer token (ANALYTICS_TOKEN env var).
func AnalyticsStatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Check auth token
	token := os.Getenv("ANALYTICS_TOKEN")
	if token != "" {
		auth := r.Header.Get("Authorization")
		if auth != "Bearer "+token {
			http.Error(w, `{"error":"unauthorized"}`, http.StatusUnauthorized)
			return
		}
	}

	period := r.URL.Query().Get("period")
	var since time.Time
	now := time.Now().UTC()
	switch period {
	case "24h":
		since = now.Add(-24 * time.Hour)
	case "30d":
		since = now.Add(-30 * 24 * time.Hour)
	case "90d":
		since = now.Add(-90 * 24 * time.Hour)
	default:
		period = "7d"
		since = now.Add(-7 * 24 * time.Hour)
	}

	db, err := getAnalyticsDB()
	if err != nil {
		http.Error(w, `{"error":"db unavailable"}`, http.StatusServiceUnavailable)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	viewsColl := db.Collection("pageviews")
	eventsColl := db.Collection("events")

	filter := bson.M{"ts": bson.M{"$gte": since}}

	// Total views
	totalViews, _ := viewsColl.CountDocuments(ctx, filter)

	// Unique visitors (distinct visitor_id)
	distinctVIDs, _ := viewsColl.Distinct(ctx, "visitor_id", filter)
	uniqueVisitors := int64(len(distinctVIDs))

	// Top pages
	topPages := aggregateTopField(ctx, viewsColl, "path", filter, 10)

	// Top referrers (exclude empty)
	refFilter := bson.M{"ts": bson.M{"$gte": since}, "referrer": bson.M{"$ne": ""}}
	topReferrers := aggregateTopField(ctx, viewsColl, "referrer", refFilter, 10)

	// Top events
	topEvents := aggregateTopField(ctx, eventsColl, "name", filter, 10)

	// Views by day
	viewsByDay := aggregateByDay(ctx, viewsColl, filter)

	resp := StatsResponse{
		Period:         period,
		TotalViews:     totalViews,
		UniqueVisitors: uniqueVisitors,
		TopPages:       toPageStats(topPages),
		TopReferrers:   toReferrerStats(topReferrers),
		TopEvents:      toEventStats(topEvents),
		ViewsByDay:     viewsByDay,
	}

	json.NewEncoder(w).Encode(resp)
}

// ---------- Aggregation helpers ----------

type fieldCount struct {
	ID    string `bson:"_id"`
	Count int64  `bson:"count"`
}

func aggregateTopField(ctx context.Context, coll *mongo.Collection, field string, filter bson.M, limit int) []fieldCount {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$" + field}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
		{{Key: "$sort", Value: bson.D{{Key: "count", Value: -1}}}},
		{{Key: "$limit", Value: limit}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("analytics: aggregate %s: %v", field, err)
		return nil
	}
	defer cursor.Close(ctx)

	var results []fieldCount
	cursor.All(ctx, &results)
	return results
}

func aggregateByDay(ctx context.Context, coll *mongo.Collection, filter bson.M) []DayStat {
	pipeline := mongo.Pipeline{
		{{Key: "$match", Value: filter}},
		{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{{Key: "$dateToString", Value: bson.D{
				{Key: "format", Value: "%Y-%m-%d"},
				{Key: "date", Value: "$ts"},
			}}}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
		}}},
		{{Key: "$sort", Value: bson.D{{Key: "_id", Value: 1}}}},
	}

	cursor, err := coll.Aggregate(ctx, pipeline)
	if err != nil {
		log.Printf("analytics: aggregate by day: %v", err)
		return nil
	}
	defer cursor.Close(ctx)

	var results []DayStat
	cursor.All(ctx, &results)
	return results
}

func toPageStats(fcs []fieldCount) []PageStat {
	out := make([]PageStat, len(fcs))
	for i, fc := range fcs {
		out[i] = PageStat{Path: fc.ID, Views: fc.Count}
	}
	return out
}

func toReferrerStats(fcs []fieldCount) []ReferrerStat {
	out := make([]ReferrerStat, len(fcs))
	for i, fc := range fcs {
		out[i] = ReferrerStat{Referrer: fc.ID, Count: fc.Count}
	}
	return out
}

func toEventStats(fcs []fieldCount) []EventStat {
	out := make([]EventStat, len(fcs))
	for i, fc := range fcs {
		out[i] = EventStat{Name: fc.ID, Count: fc.Count}
	}
	return out
}
