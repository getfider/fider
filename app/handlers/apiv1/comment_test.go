package apiv1_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/getfider/fider/app/handlers/apiv1"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/query"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/mock"
)

// Shared response structs for testing
type paginationInfo struct {
	HasNext    bool   `json:"hasNext"`
	NextCursor string `json:"nextCursor,omitempty"`
}

type commentsResponse struct {
	Data       []*entity.CommentRef `json:"data"`
	Pagination *paginationInfo      `json:"pagination"`
}

func TestAllCommentsHandler(t *testing.T) {
	RegisterT(t)

	// Create 8 test comment refs (less than limit to ensure hasNext = false)
	testComments := make([]*entity.CommentRef, 8)
	for i := 0; i < 8; i++ {
		testComments[i] = &entity.CommentRef{
			ID:        i + 1,
			CreatedAt: time.Now().Add(time.Duration(i) * time.Minute),
			UserID:    1,
			PostID:    1,
		}
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentRefs) error {
		q.Result = testComments
		return nil
	})

	code, body := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("/api/v1/comments?limit=10").
		Execute(apiv1.AllComments())

	Expect(code).Equals(http.StatusOK)

	// Properly unmarshal the JSON response
	var response commentsResponse
	err := json.Unmarshal(body.Body.Bytes(), &response)
	Expect(err).IsNil()

	// Test that the response has the correct structure
	Expect(response.Data).IsNotNil()
	Expect(len(response.Data)).Equals(8) // Should have exactly 8 items
	Expect(response.Pagination).IsNotNil()

	// Since we got fewer results than the limit (8 < 10), hasNext should be false
	Expect(response.Pagination.HasNext).IsFalse()
	Expect(response.Pagination.NextCursor).Equals("") // No cursor when no next page
}

func TestAllCommentsHandler_WithLimit(t *testing.T) {
	RegisterT(t)

	// Create 5 test comment refs (less than the requested limit)
	testComments := make([]*entity.CommentRef, 5)
	for i := 0; i < 5; i++ {
		testComments[i] = &entity.CommentRef{
			ID:        i + 1,
			CreatedAt: time.Now().Add(time.Duration(i) * time.Minute),
			UserID:    1,
			PostID:    1,
		}
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentRefs) error {
		// Should receive the requested limit of 20
		Expect(q.Limit).Equals(20)
		q.Result = testComments
		return nil
	})

	code, body := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("/api/v1/comments?limit=20").
		Execute(apiv1.AllComments())

	Expect(code).Equals(http.StatusOK)

	// Properly unmarshal the JSON response
	var response commentsResponse
	err := json.Unmarshal(body.Body.Bytes(), &response)
	Expect(err).IsNil()

	// Test that the response has the correct structure
	Expect(response.Data).IsNotNil()
	Expect(len(response.Data)).Equals(5) // Should have exactly 5 items
	Expect(response.Pagination).IsNotNil()

	// Since we got fewer results than the limit (5 < 20), hasNext should be false
	Expect(response.Pagination.HasNext).IsFalse()
	Expect(response.Pagination.NextCursor).Equals("") // No cursor when no next page
}

func TestAllCommentsHandler_WithPagination(t *testing.T) {
	RegisterT(t)

	// Create 20 test comment refs (equal to limit to simulate full page)
	testComments := make([]*entity.CommentRef, 20)
	baseTime := time.Now()
	for i := 0; i < 20; i++ {
		testComments[i] = &entity.CommentRef{
			ID:        i + 1,
			CreatedAt: baseTime.Add(time.Duration(i) * time.Minute),
			UserID:    1,
			PostID:    1,
		}
	}

	bus.AddHandler(func(ctx context.Context, q *query.GetCommentRefs) error {
		// Should receive the requested limit of 20
		Expect(q.Limit).Equals(20)
		q.Result = testComments
		return nil
	})

	code, body := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("/api/v1/comments?limit=20").
		Execute(apiv1.AllComments())

	Expect(code).Equals(http.StatusOK)

	// Properly unmarshal the JSON response
	var response commentsResponse
	err := json.Unmarshal(body.Body.Bytes(), &response)
	Expect(err).IsNil()

	// Test that the response has the correct structure
	Expect(response.Data).IsNotNil()
	Expect(len(response.Data)).Equals(20) // Should have exactly 20 items
	Expect(response.Pagination).IsNotNil()

	// Since we got exactly the limit (20 == 20), hasNext should be true (might be more pages)
	Expect(response.Pagination.HasNext).IsTrue()
	Expect(response.Pagination.NextCursor).IsNotEmpty() // Should have cursor for next page

	// The cursor should be the timestamp of the last comment
	expectedCursor := testComments[19].CreatedAt.Format(time.RFC3339)
	Expect(response.Pagination.NextCursor).Equals(expectedCursor)
}

func TestAllCommentsHandler_FullPaginationFlow(t *testing.T) {
	RegisterT(t)

	// Create 15 test comment refs with incrementing timestamps
	testComments := make([]*entity.CommentRef, 15)
	baseTime := time.Date(2023, 10, 1, 10, 0, 0, 0, time.UTC) // Use fixed time for predictability
	for i := 0; i < 15; i++ {
		testComments[i] = &entity.CommentRef{
			ID:        i + 1,
			CreatedAt: baseTime.Add(time.Duration(i) * time.Minute),
			UserID:    1,
			PostID:    1,
		}
	}

	// Set up mock handler that will be called twice
	callCount := 0
	bus.AddHandler(func(ctx context.Context, q *query.GetCommentRefs) error {
		callCount++

		if callCount == 1 {
			// First call: return first 10 comments (limit=10, no since filter)
			Expect(q.Limit).Equals(10)
			Expect(q.Since.IsZero()).IsTrue() // No since parameter on first call
			q.Result = testComments[:10]      // Return first 10 comments
		} else if callCount == 2 {
			// Second call: return remaining 5 comments (limit=10, with since filter)
			Expect(q.Limit).Equals(10)
			Expect(q.Since.IsZero()).IsFalse() // Should have since parameter
			// The since parameter should match the timestamp of the 10th comment
			expectedSince := testComments[9].CreatedAt // 10th comment (index 9)
			// Compare timestamps with some tolerance for RFC3339 parsing
			timeDiff := q.Since.Sub(expectedSince)
			if timeDiff < 0 {
				timeDiff = -timeDiff
			}
			Expect(timeDiff < time.Second).IsTrue() // Should be within 1 second

			// Return comments that would come after the since timestamp
			// In real implementation, this would filter by COALESCE(edited_at, created_at) >= since
			q.Result = testComments[10:] // Return remaining 5 comments (indexes 10-14)
		}

		return nil
	})

	// FIRST API CALL: Get first page (10 comments, limit 10)
	code1, body1 := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL("/api/v1/comments?limit=10").
		Execute(apiv1.AllComments())

	Expect(code1).Equals(http.StatusOK)

	var response1 commentsResponse
	err1 := json.Unmarshal(body1.Body.Bytes(), &response1)
	Expect(err1).IsNil()

	// Validate first page response
	Expect(response1.Data).IsNotNil()
	Expect(len(response1.Data)).Equals(10) // Should have exactly 10 items
	Expect(response1.Pagination).IsNotNil()

	// Since we got exactly the limit (10 == 10), hasNext should be true
	Expect(response1.Pagination.HasNext).IsTrue()
	Expect(response1.Pagination.NextCursor).IsNotEmpty()

	// The cursor should be the timestamp of the last comment from first page
	expectedCursor1 := testComments[9].CreatedAt.Format(time.RFC3339) // 10th comment (index 9)
	Expect(response1.Pagination.NextCursor).Equals(expectedCursor1)

	// SECOND API CALL: Get second page using the cursor from first page
	secondPageURL := "/api/v1/comments?limit=10&since=" + url.QueryEscape(response1.Pagination.NextCursor)
	code2, body2 := mock.NewServer().
		OnTenant(mock.DemoTenant).
		AsUser(mock.JonSnow).
		WithURL(secondPageURL).
		Execute(apiv1.AllComments())

	Expect(code2).Equals(http.StatusOK)

	var response2 commentsResponse
	err2 := json.Unmarshal(body2.Body.Bytes(), &response2)
	Expect(err2).IsNil()

	// Validate second page response
	Expect(response2.Data).IsNotNil()
	Expect(len(response2.Data)).Equals(5) // Should have exactly 5 items (remaining comments)
	Expect(response2.Pagination).IsNotNil()

	// Since we got fewer than the limit (5 < 10), hasNext should be false
	Expect(response2.Pagination.HasNext).IsFalse()
	Expect(response2.Pagination.NextCursor).Equals("") // No cursor when no more pages

	// Validate that we got all 15 comments across both pages
	// First page should have comments 1-10, second page should have comments 11-15
	Expect(response1.Data[0].ID).Equals(1)  // First comment from first page
	Expect(response1.Data[9].ID).Equals(10) // Last comment from first page
	Expect(response2.Data[0].ID).Equals(11) // First comment from second page
	Expect(response2.Data[4].ID).Equals(15) // Last comment from second page (index 4 = 5th item)

	// Verify that we made exactly 2 calls to the handler
	Expect(callCount).Equals(2)
}
