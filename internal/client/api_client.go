package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"marketplace-notifications/internal/config"
	"marketplace-notifications/internal/marketplaces"
	"marketplace-notifications/internal/marketplaces/wb"
	"net/http"
	"net/url"
	"strconv"

	"golang.org/x/time/rate"
)

type APIClient struct {
	config     *config.APIConfig
	httpClient *http.Client
	wbLimiter  *rate.Limiter
}

func NewAPIClient(config *config.APIConfig) *APIClient {
	wbLimiter := rate.NewLimiter(rate.Limit(config.WB.RPS), config.WB.Burst)

	return &APIClient{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		wbLimiter: wbLimiter,
	}
}

func (client *APIClient) FetchWBQuestions() ([]wb.Question, error) {
	jsonData, err := client.FetchWBData(marketplaces.Question)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch questions: %w", err)
	}

	var questionsResponse wb.QuestionsResponse
	if err := json.Unmarshal(jsonData, &questionsResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal questions response: %w", err)
	}

	return questionsResponse.Data.Questions, nil
}

func (client *APIClient) FetchWBFeedbacks() ([]wb.Feedback, error) {
	jsonData, err := client.FetchWBData(marketplaces.Feedback)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feedbacks: %w", err)
	}

	var feedbacksResponse wb.FeedbacksResponse
	if err := json.Unmarshal(jsonData, &feedbacksResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal feedbacks response: %w", err)
	}

	return feedbacksResponse.Data.Feedbacks, nil
}

func (client *APIClient) FetchWBData(reactionType marketplaces.UserReactionType) ([]byte, error) {
	if err := client.wbLimiter.Wait(context.Background()); err != nil {
		return nil, fmt.Errorf("WB rate limiter error: %w", err)
	}

	var urlToParse string
	var maxNewReactions int
	if reactionType == marketplaces.Question {
		urlToParse = client.config.WB.QuestionsURL()
		maxNewReactions = client.config.WB.MaxNewQuestions
	} else {
		urlToParse = client.config.WB.FeedbacksURL()
		maxNewReactions = client.config.WB.MaxNewFeedbacks
	}

	baseURL, err := url.Parse(urlToParse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse URL: %w", err)
	}

	query := url.Values{}

	query.Set("isAnswered", strconv.FormatBool(true))
	query.Set("take", strconv.Itoa(maxNewReactions))
	query.Set("skip", strconv.Itoa(0))

	baseURL.RawQuery = query.Encode()

	fullURL := baseURL.String()

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("content-type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", client.config.WB.JWT))

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body")
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d instead of 200: %s", resp.StatusCode, body)
	}

	return body, nil
}
