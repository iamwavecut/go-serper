package serper

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type Logger interface {
	Debug(string, ...any)
	Info(string, ...any)
	Warn(string, ...any)
	Error(string, ...any)
}

type Client struct {
	httpClient      *http.Client
	apiKey          string
	baseURL         string
	logger          Logger
	retryCount      int
	retryBaseDelay  time.Duration
	requestTimeout  time.Duration
	totalTimeout    time.Duration
	lastRawResponse string
}

type SearchRequest struct {
	Query       string `json:"q"`
	Country     string `json:"gl,omitempty"`
	Location    string `json:"location,omitempty"`
	Language    string `json:"hl,omitempty"`
	AutoCorrect bool   `json:"autocorrect,omitempty"`
	MaxResults  int    `json:"num,omitempty"`
	Page        int    `json:"page,omitempty"`
}

type SearchResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	Results          []SearchResult    `json:"organic,omitempty"`
	KnowledgeGraph   *KnowledgeGraph   `json:"knowledgeGraph,omitempty"`
	AnswerBox        *AnswerBox        `json:"answerBox,omitempty"`
	PeopleAlsoAsk    []PeopleAlsoAsk   `json:"peopleAlsoAsk,omitempty"`
	RelatedSearches  []RelatedSearch   `json:"relatedSearches,omitempty"`
	TopStories       []TopStory        `json:"topStories,omitempty"`
}

type ImageResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	Images           []ImageResult     `json:"images,omitempty"`
}

type VideoResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	Videos           []VideoResult     `json:"videos,omitempty"`
}

type PlaceResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	Places           []PlaceResult     `json:"places,omitempty"`
}

type NewsResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	News             []NewsResult      `json:"news,omitempty"`
}

type ShoppingResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	Shopping         []ShoppingResult  `json:"shopping,omitempty"`
}

type ScholarResponse struct {
	SearchParameters *SearchParameters `json:"searchParameters,omitempty"`
	Results          []ScholarResult   `json:"organic,omitempty"`
}

type SearchParameters struct {
	Query  string `json:"q"`
	Type   string `json:"type,omitempty"`
	Engine string `json:"engine,omitempty"`
}

type SearchResult struct {
	Title      string            `json:"title"`
	URL        string            `json:"link"`
	Content    string            `json:"snippet"`
	Position   int               `json:"position"`
	Date       string            `json:"date,omitempty"`
	Sitelinks  []Sitelink        `json:"sitelinks,omitempty"`
	Attributes map[string]string `json:"attributes,omitempty"`
}

type KnowledgeGraph struct {
	Title             string            `json:"title,omitempty"`
	Type              string            `json:"type,omitempty"`
	Website           string            `json:"website,omitempty"`
	ImageURL          string            `json:"imageUrl,omitempty"`
	Description       string            `json:"description,omitempty"`
	DescriptionSource string            `json:"descriptionSource,omitempty"`
	DescriptionLink   string            `json:"descriptionLink,omitempty"`
	Attributes        map[string]string `json:"attributes,omitempty"`
}

type AnswerBox struct {
	Answer  string `json:"answer,omitempty"`
	Title   string `json:"title,omitempty"`
	URL     string `json:"link,omitempty"`
	Snippet string `json:"snippet,omitempty"`
}

type PeopleAlsoAsk struct {
	Question string `json:"question"`
	Snippet  string `json:"snippet,omitempty"`
	Title    string `json:"title,omitempty"`
	URL      string `json:"link,omitempty"`
}

type RelatedSearch struct {
	Query string `json:"query"`
}

type ImageResult struct {
	Title           string `json:"title,omitempty"`
	URL             string `json:"imageUrl"`
	ImageWidth      int    `json:"imageWidth,omitempty"`
	ImageHeight     int    `json:"imageHeight,omitempty"`
	ThumbnailURL    string `json:"thumbnailUrl,omitempty"`
	ThumbnailWidth  int    `json:"thumbnailWidth,omitempty"`
	ThumbnailHeight int    `json:"thumbnailHeight,omitempty"`
	Source          string `json:"source,omitempty"`
	Domain          string `json:"domain,omitempty"`
	Link            string `json:"link,omitempty"`
	GoogleURL       string `json:"googleUrl,omitempty"`
	Position        int    `json:"position,omitempty"`
}

type TopStory struct {
	Title    string `json:"title,omitempty"`
	URL      string `json:"link,omitempty"`
	Source   string `json:"source,omitempty"`
	Date     string `json:"date,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
}

type VideoResult struct {
	Title    string `json:"title,omitempty"`
	URL      string `json:"link,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	Duration string `json:"duration,omitempty"`
	Source   string `json:"source,omitempty"`
	Channel  string `json:"channel,omitempty"`
	Date     string `json:"date,omitempty"`
	Position int    `json:"position,omitempty"`
}

type PlaceResult struct {
	Position    int     `json:"position,omitempty"`
	Name        string  `json:"name,omitempty"`
	Address     string  `json:"address,omitempty"`
	Latitude    float64 `json:"latitude,omitempty"`
	Longitude   float64 `json:"longitude,omitempty"`
	Rating      float64 `json:"rating,omitempty"`
	RatingCount int     `json:"ratingCount,omitempty"`
	Category    string  `json:"category,omitempty"`
	Identifier  string  `json:"identifier,omitempty"`
}

type NewsResult struct {
	Title    string `json:"title,omitempty"`
	URL      string `json:"link,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
	Date     string `json:"date,omitempty"`
	Source   string `json:"source,omitempty"`
	ImageURL string `json:"imageUrl,omitempty"`
	Position int    `json:"position,omitempty"`
}

type ShoppingResult struct {
	Title       string  `json:"title,omitempty"`
	Source      string  `json:"source,omitempty"`
	URL         string  `json:"link,omitempty"`
	Price       string  `json:"price,omitempty"`
	Delivery    string  `json:"delivery,omitempty"`
	ImageURL    string  `json:"imageUrl,omitempty"`
	Rating      float64 `json:"rating,omitempty"`
	RatingCount int     `json:"ratingCount,omitempty"`
	Offers      string  `json:"offers,omitempty"`
	ProductID   string  `json:"productId,omitempty"`
	Position    int     `json:"position,omitempty"`
}

type ScholarResult struct {
	Title           string `json:"title,omitempty"`
	URL             string `json:"link,omitempty"`
	PublicationInfo string `json:"publicationInfo,omitempty"`
	Snippet         string `json:"snippet,omitempty"`
	Year            int    `json:"year,omitempty"`
	CitedBy         int    `json:"citedBy,omitempty"`
}

type Sitelink struct {
	Title string `json:"title"`
	URL   string `json:"link"`
}

type responseCapture struct {
	base         http.RoundTripper
	lastResponse *string
}

func (rc *responseCapture) RoundTrip(req *http.Request) (*http.Response, error) {
	resp, err := rc.base.RoundTrip(req)
	if err != nil {
		return resp, err
	}
	if resp.Body != nil && rc.lastResponse != nil {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr == nil {
			*rc.lastResponse = string(bodyBytes)
			resp.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		} else {
			resp.Body = io.NopCloser(bytes.NewReader([]byte{}))
		}
	}
	return resp, err
}

type Option func(*Client)

func WithHTTPClient(client *http.Client) Option {
	return func(c *Client) {
		c.httpClient = client
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

func WithRetryConfig(retryCount int, baseDelay time.Duration) Option {
	return func(c *Client) {
		c.retryCount = retryCount
		c.retryBaseDelay = baseDelay
	}
}

func WithTimeouts(requestTimeout, totalTimeout time.Duration) Option {
	return func(c *Client) {
		c.requestTimeout = requestTimeout
		c.totalTimeout = totalTimeout
		transport := &responseCapture{base: http.DefaultTransport, lastResponse: &c.lastRawResponse}
		c.httpClient = &http.Client{Timeout: requestTimeout, Transport: transport}
	}
}

func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

func NewClient(apiKey string, opts ...Option) *Client {
	client := &Client{
		apiKey:         apiKey,
		baseURL:        "https://google.serper.dev",
		retryCount:     3,
		retryBaseDelay: 1 * time.Second,
		requestTimeout: 10 * time.Second,
		totalTimeout:   30 * time.Second,
	}
	for _, opt := range opts {
		opt(client)
	}
	if client.httpClient == nil {
		transport := &responseCapture{base: http.DefaultTransport, lastResponse: &client.lastRawResponse}
		client.httpClient = &http.Client{Timeout: client.requestTimeout, Transport: transport}
	}
	return client
}

func (c *Client) LastRawResponse() string {
	return c.lastRawResponse
}

func (c *Client) Search(ctx context.Context, req SearchRequest) (*SearchResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result SearchResponse
	err := c.performRequestWithRetry(ctx, req, "/search", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchImages(ctx context.Context, req SearchRequest) (*ImageResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result ImageResponse
	err := c.performRequestWithRetry(ctx, req, "/images", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchVideos(ctx context.Context, req SearchRequest) (*VideoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result VideoResponse
	err := c.performRequestWithRetry(ctx, req, "/videos", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchPlaces(ctx context.Context, req SearchRequest) (*PlaceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result PlaceResponse
	err := c.performRequestWithRetry(ctx, req, "/places", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchNews(ctx context.Context, req SearchRequest) (*NewsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result NewsResponse
	err := c.performRequestWithRetry(ctx, req, "/news", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchShopping(ctx context.Context, req SearchRequest) (*ShoppingResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result ShoppingResponse
	err := c.performRequestWithRetry(ctx, req, "/shopping", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) SearchScholar(ctx context.Context, req SearchRequest) (*ScholarResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, c.totalTimeout)
	defer cancel()
	var result ScholarResponse
	err := c.performRequestWithRetry(ctx, req, "/scholar", &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) performRequestWithRetry(ctx context.Context, req SearchRequest, endpoint string, result any) error {
	var lastErr error
	for attempt := 0; attempt <= c.retryCount; attempt++ {
		if attempt > 0 {
			if c.logger != nil {
				c.logger.Warn("retrying request", "attempt", attempt)
			}
			delay := c.retryBaseDelay * time.Duration(2*attempt-1)
			select {
			case <-ctx.Done():
				return fmt.Errorf("request cancelled during retry delay: %w", ctx.Err())
			case <-time.After(delay):
			}
		}
		err := c.performRequest(ctx, req, endpoint, result)
		if err == nil {
			if c.logger != nil && attempt > 0 {
				c.logger.Info("request succeeded after retry", "attempt", attempt)
			}
			return nil
		}
		lastErr = err
		if !c.shouldRetry(err, attempt) {
			break
		}
	}
	return fmt.Errorf("request failed after %d attempts: %w", c.retryCount+1, lastErr)
}

func (c *Client) performRequest(ctx context.Context, req SearchRequest, endpoint string, result any) error {
	if c.logger != nil {
		c.logger.Debug("sending request", "endpoint", endpoint)
	}
	if len(req.Query) > 400 {
		if c.logger != nil {
			c.logger.Info("clamping query to 400 chars", "len", len(req.Query))
		}
		req.Query = req.Query[:400]
	}
	body, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("marshal request: %w", err)
	}
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("X-API-KEY", c.apiKey)
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		if c.logger != nil {
			c.logger.Error("request failed", "endpoint", endpoint, "error", err)
		}
		return fmt.Errorf("request error: %w", err)
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		if c.logger != nil {
			c.logger.Error("api error", "status", resp.StatusCode, "body", string(respBody))
		}
		return fmt.Errorf("api error: status %d, body: %s", resp.StatusCode, string(respBody))
	}
	if err := json.Unmarshal(respBody, result); err != nil {
		if c.logger != nil {
			c.logger.Error("unmarshal error", "error", err)
		}
		return fmt.Errorf("unmarshal response: %w", err)
	}
	if c.logger != nil {
		c.logger.Debug("received response", "endpoint", endpoint)
	}
	return nil
}

func (c *Client) shouldRetry(err error, attempt int) bool {
	if attempt >= c.retryCount {
		return false
	}
	msg := strings.ToLower(err.Error())
	networkErrors := []string{"timeout", "connection", "network", "dial", "dns", "500", "502", "503", "504", "internal server error", "bad gateway", "service unavailable", "gateway timeout"}
	for _, s := range networkErrors {
		if strings.Contains(msg, s) {
			return true
		}
	}
	if strings.Contains(msg, "rate limit") || strings.Contains(msg, "429") {
		return true
	}
	noRetry := []string{"401", "403", "400", "unauthorized", "forbidden", "bad request", "invalid api key", "authentication", "unmarshal", "json", "parse"}
	for _, s := range noRetry {
		if strings.Contains(msg, s) {
			return false
		}
	}
	return true
}
