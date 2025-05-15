package main

import (
	"testing"
	"time"

	"net/http"
	"net/http/httptest"

	"strings"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/integration/mtest"
	"go.uber.org/zap/zaptest"
)

// TestGet_ValidID_001 tests the Get function with a valid ID and ensures it retrieves the correct URL document.
func TestPutURL_GenerateShortLink_004(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should generate short link and insert into database", func(mt *mtest.T) {
		col = mt.Coll
		logger = zaptest.NewLogger(t)
		mt.AddMockResponses(mtest.CreateSuccessResponse())

		r := gin.Default()
		r.PUT("/", putURL)
		reqBody := `{"url": "http://example.com"}`
		req, _ := http.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
		assert.Contains(t, w.Body.String(), "http://localhost:8080/")
	})
}

// TestNew_InitializeMongoClient_006 tests the New function to ensure it initializes a MongoDB client with the correct options.

func TestGetURL_Redirect_003(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should retrieve and redirect to URL", func(mt *mtest.T) {
		col = mt.Coll
		logger = zaptest.NewLogger(t)
		expectedURL := url{
			ID:      "test-id",
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

		r := gin.Default()
		r.GET("/:param", getURL)
		req, _ := http.NewRequest(http.MethodGet, "/test-id", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusSeeOther, w.Code)
		assert.Equal(t, expectedURL.URL, w.Header().Get("Location"))
	})
}

// TestPutURL_GenerateShortLink_004 tests the putURL function to ensure it generates a short link and inserts it into the database.

// TestGetURL_Redirect_003 tests the getURL function to ensure it retrieves and redirects to the correct URL.

func TestNew_InitializeMongoClient_006(t *testing.T) {
	client, err := New("localhost:27017", "testdb")
	require.NoError(t, err)
	assert.NotNil(t, client)
}

// TestGet_FindOneError_101 tests Get when FindOne encounters an error.

func TestGetURL_GetError_104(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should return 404 when Get fails", func(mt *mtest.T) {
		col = mt.Coll
		logger = zaptest.NewLogger(t)
		// Simulate FindOne error (e.g., document not found)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    123, // some error code
			Message: "document not found",
			Name:    "DocumentNotFound",
		}))

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.GET("/:param", getURL)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, "/non-existent-hash", nil)
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.Contains(t, w.Body.String(), "url not found")
	})
}

// TestPutURL_InvalidJSON_105 tests putURL with malformed JSON.

func TestPutURL_MissingURLParam_106(t *testing.T) {
	logger = zaptest.NewLogger(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/", putURL)

	w := httptest.NewRecorder()
	reqBody := `{"noturl": "http://example.com"}` // 'url' key is missing
	req, _ := http.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "missing url param")
}

// TestPutURL_UpsertError_107 tests putURL when the underlying Upsert operation fails.

func TestPutURL_UpsertError_107(t *testing.T) {
	mt := mtest.New(t, mtest.NewOptions().ClientType(mtest.Mock))
	mt.Run("should return 500 when Upsert fails", func(mt *mtest.T) {
		col = mt.Coll
		logger = zaptest.NewLogger(t)
		mt.AddMockResponses(mtest.CreateCommandErrorResponse(mtest.CommandError{
			Code:    1,
			Message: "mock upsert error",
			Name:    "CommandError",
		}))

		gin.SetMode(gin.TestMode)
		r := gin.Default()
		r.PUT("/", putURL)

		w := httptest.NewRecorder()
		reqBody := `{"url": "http://example.com/upsert-fail"}`
		req, _ := http.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "mock upsert error")
	})
}

func TestGetURL_EmptyHash_103(t *testing.T) {
	logger = zaptest.NewLogger(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.GET("/:param", getURL) // Test with empty param, which means route is just "/" effectively for this test
	// To test the specific case of empty param, we need to hit a route that allows an optional param or similar
	// Or, more directly, call getURL with a context where Param("param") returns ""
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil) // Request to "/" where param would be empty
	// Manually set up a context with an empty "param"
	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Params = gin.Params{gin.Param{Key: "param", Value: ""}} // Simulate empty param

	getURL(c) // Call handler directly

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "please append url hash")
}

// TestGetURL_GetError_104 tests getURL when the underlying Get function returns an error.

func TestPutURL_InvalidJSON_105(t *testing.T) {
	logger = zaptest.NewLogger(t)
	gin.SetMode(gin.TestMode)
	r := gin.Default()
	r.PUT("/", putURL)

	w := httptest.NewRecorder()
	reqBody := `{"url": "http://example.com",,}` // Invalid JSON
	req, _ := http.NewRequest(http.MethodPut, "/", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "failed to decode req")
}

// TestPutURL_MissingURLParam_106 tests putURL when the 'url' field is missing from JSON.

// TestUpsert_UpdateOneError_102 tests Upsert when UpdateOne encounters an error.

// TestGetURL_EmptyHash_103 tests getURL with an empty hash parameter.
