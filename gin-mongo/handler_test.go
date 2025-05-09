package main

import (
	"testing"
	"time"

	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap/zaptest"
)

// TestGet_ValidID_123 tests the Get function with a valid ID.
func TestPutURL_ValidAndInvalidJSON_321(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should handle valid and invalid JSON inputs", func(mt *mtest.T) {
		logger = zaptest.NewLogger(t)
		col = mt.Coll

		router := gin.Default()
		router.PUT("/", putURL)

		// Mock valid JSON
		validJSON := `{"url": "http://example.com"}`
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		req, _ := http.NewRequest("PUT", "/", strings.NewReader(validJSON))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		// Mock invalid JSON
		invalidJSON := `{"invalid": "data"}`
		req, _ = http.NewRequest("PUT", "/", strings.NewReader(invalidJSON))
		req.Header.Set("Content-Type", "application/json")
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
		}
	})
}

// TestGetURL_ValidAndInvalidHash_789 tests the `getURL` function with valid and invalid hash inputs.

func TestGetURL_ValidAndInvalidHash_789(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should handle valid and invalid hash inputs", func(mt *mtest.T) {
		logger = zaptest.NewLogger(t)
		col = mt.Coll

		router := gin.Default()
		router.GET("/:param", getURL)

		// Mock valid hash
		validHash := "test-id"
		expectedURL := url{
			ID:      validHash,
			Created: time.Now(),
			Updated: time.Now(),
			URL:     "http://example.com",
		}
		first := mtest.CreateCursorResponse(
			1,
			"testdb.testcoll",
			mtest.FirstBatch,
			bson.D{
				{Key: "_id", Value: expectedURL.ID},
				{Key: "created", Value: expectedURL.Created},
				{Key: "updated", Value: expectedURL.Updated},
				{Key: "url", Value: expectedURL.URL},
			},
		)
		second := mtest.CreateCursorResponse(0, "testdb.testcoll", mtest.NextBatch)
		mt.AddMockResponses(first, second)

		req, _ := http.NewRequest("GET", "/"+validHash, nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("expected status %d, got %d", http.StatusSeeOther, w.Code)
		}

		// Mock invalid hash
		invalidHash := "invalid-id"
		mt.AddMockResponses(mtest.CreateCursorResponse(0, "testdb.testcoll", mtest.FirstBatch))

		req, _ = http.NewRequest("GET", "/"+invalidHash, nil)
		w = httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, w.Code)
		}
	})
}

// TestGetURL_EmptyHashParam_555 tests getURL with an empty hash parameter.

// TestPutURL_ValidAndInvalidJSON_321 tests the putURL function with valid and invalid JSON inputs.

// TestUpsert_ValidURL_456 tests the Upsert function with a valid URL.

func TestGetURL_EmptyHashParam_555(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger = zaptest.NewLogger(t) // Initialize logger as getURL might use it

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	// Create a request, though its details might not be crucial if params are manually set
	c.Request, _ = http.NewRequest("GET", "/testroute", nil)
	// Manually set Param to simulate an empty hash, bypassing router complexities
	c.Params = gin.Params{gin.Param{Key: "param", Value: ""}}

	getURL(c) // Directly call handler with the crafted context

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.JSONEq(t, `{"error": "please append url hash"}`, w.Body.String())
}
