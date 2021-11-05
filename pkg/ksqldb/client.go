package ksqldb

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

const (

	// Default response content type for pull & push queries
	// In the case of a successful query, if the content type is application/vnd.ksqlapi.delimited.v1,
	// the results are returned as a header JSON object followed by zero or more JSON arrays that are delimited by newlines.
	ContentTypeDelim = "application/vnd.ksqlapi.delimited.v1; charset=utf-8"

	// Default serialization format for requests and responses.
	ContentTypeDefault = "application/vnd.ksql.v1+json; charset=utf-8"

	// EndpointRunStreamQuery is used to run push and pull queries.
	// These endpoints are only available when using HTTP 2.
	EndpointRunStreamQuery string = "/query-stream"
	// EndpointCloseQuery used to terminates a push query.
	EndpointCloseQuery string = "/close-query"
)

type BasicAuth struct {
	Username string
	Password string
}

type QueryRequest struct {
	SQL        string            `json:"sql"`
	Properties map[string]string `json:"streamsProperties,omitempty"`
}

type QueryResponse struct {
	QueryID     string   `json:"queryID"`
	ColumnNames []string `json:"columnNames"`
	ColumnTypes []string `json:"columnTypes"`
}

// NewClient creates a new ksqlDB client.
func NewClient(url string, options ...func(*Client)) *Client {
	client := &Client{
		url: url,
		HTTPClient: &http.Client{
			// In go, the standard http.Client is used for HTTP/2 requests as well.
			// The only difference is the usage of http2.Transport instead of http.Transport in the client’s Transport field
			Transport: &http2.Transport{
				AllowHTTP: true,
				// Pretend we are dialing a TLS endpoint.
				// Note, we ignore the passed tls.Config
				DialTLS: func(network string, addr string, cfg *tls.Config) (net.Conn, error) {
					return net.Dial(network, addr)
				},
			},
		},
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

type Client struct {
	url        string
	HTTPClient *http.Client
	BasicAuth  *BasicAuth
}

// postRequest makes POST request.
func (c *Client) req(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	content, err := json.Marshal(payload)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("%s/%s", c.url, endpoint), bytes.NewBuffer(content))

	if err != nil {
		return nil, err
	}

	switch endpoint {
	case EndpointRunStreamQuery:
		req.Header.Set("Content-Type", ContentTypeDelim)
		req.Header.Set("Accept-Encoding", "identity")
		req.Header.Set("Accept", ContentTypeDelim)
	default:
		req.Header.Set("Content-Type", ContentTypeDefault)
		req.Header.Set("Accept-Encoding", "identity")
		req.Header.Set("Accept", ContentTypeDefault)
	}

	if c.BasicAuth != nil {
		req.SetBasicAuth(c.BasicAuth.Username, c.BasicAuth.Password)
	}

	return c.HTTPClient.Do(req)
}

func (c *Client) Pull(ctx context.Context, q *QueryRequest) (*QueryResponse, []Row, error) {
	res, err := c.req(ctx, EndpointRunStreamQuery, q)
	qr := new(QueryResponse)
	rows := []Row{}

	if err != nil {
		return qr, rows, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return qr, rows, err
		}

		return qr, rows, fmt.Errorf("%s:%s", http.StatusText(res.StatusCode), string(data))
	}

	scn := bufio.NewScanner(bufio.NewReader(res.Body))

	for scn.Scan() {
		if len(qr.ColumnNames) <= 0 {
			if err := json.Unmarshal([]byte(scn.Text()), qr); err != nil {
				log.Println(err)
			}
			continue
		}

		row := Row{}

		if err := json.Unmarshal([]byte(scn.Text()), &row); err != nil {
			log.Println(err)
		} else {
			rows = append(rows, row)
		}
	}

	if err := scn.Err(); err != io.EOF && err != nil {
		return qr, rows, err
	}

	return qr, rows, nil
}

func (c *Client) Push(ctx context.Context, q *QueryRequest, cb func(qr *QueryResponse, row Row)) error {
	res, err := c.req(ctx, EndpointRunStreamQuery, q)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			return err
		}

		return fmt.Errorf("%s:%s", http.StatusText(res.StatusCode), string(data))
	}

	qr := new(QueryResponse)
	scn := bufio.NewScanner(bufio.NewReader(res.Body))

	for scn.Scan() {
		if len(qr.ColumnNames) <= 0 {
			if err := json.Unmarshal([]byte(scn.Text()), qr); err != nil {
				log.Println(err)
			}
			continue
		}

		row := Row{}

		if err := json.Unmarshal([]byte(scn.Text()), &row); err != nil {
			log.Println(err)
		} else {
			cb(qr, row)
		}
	}

	if len(qr.QueryID) > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := c.CloseQuery(ctx, qr.QueryID); err != nil {
			return err
		}
	}

	if err := scn.Err(); err != io.EOF && err != nil {
		return err
	}

	return nil
}

// CloseQuery terminates a query.
func (c *Client) CloseQuery(ctx context.Context, queryId string) error {
	buf := &bytes.Buffer{}

	if err := json.NewEncoder(buf).Encode(map[string]string{"queryId": queryId}); err != nil {
		return err
	}

	_, err := c.req(ctx, EndpointCloseQuery, buf)
	return err
}
