package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"golang.org/x/term"
)

const (
	DefaultPiHoleUrl = "http://pi.hole"
	DefaultDuration  = 60 * time.Second
)

const JsonContentType = "application/json"

type PiHoleAuthorizeResponse struct {
	Session struct {
		Valid    bool   `json:"valid"`
		Validity int    `json:"validity"`
		Sid      string `json:"sid"`
	} `json:"session"`
}

type PiHoleDNSBlockingRequest struct {
	Blocking bool `json:"blocking"`
	Timer    int  `json:"timer"`
}

func main() {

	/* Read duration */
	disabledDuirationSeconds := int(DefaultDuration.Seconds())

	if len(os.Args) > 2 {
		slog.Error("takes only 1 argument")
		os.Exit(1)
	}

	if len(os.Args) > 1 {
		userInputDuration, err := time.ParseDuration(os.Args[1])
		if err != nil {
			slog.Error("couldn't parse user input", "err", err)
			os.Exit(1)
		}

		disabledDuirationSeconds = int(userInputDuration.Seconds())
	}

	/* Read Passowrd */
	fmt.Printf("Enter pi-hole password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		slog.Error("couldn't read password", "err", err)
		os.Exit(1)
	}
	fmt.Println()

	/* Authorize */
	piHoleUrl := DefaultPiHoleUrl

	authUrl := fmt.Sprintf("%s/api/auth", piHoleUrl)
	body := map[string]string{
		"password": string(bytePassword),
	}

	b, err := json.Marshal(body)
	if err != nil {
		slog.Error("couldn't marshall authorize request body", "err", err)
		os.Exit(1)
	}

	resp, err := http.Post(authUrl, JsonContentType, bytes.NewBuffer(b))
	if err != nil {
		slog.Error("couldn't authorize with pihole", "piHoleUrl", piHoleUrl, "err", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if !isSuccessfulStatusCode(resp.StatusCode) {
		slog.Error("unexpected status code", "statusCode", resp.StatusCode)
		os.Exit(1)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("couldn't read response", "err", err)
		os.Exit(1)
	}

	var piHoleAuthResp PiHoleAuthorizeResponse

	if err := json.Unmarshal(respBody, &piHoleAuthResp); err != nil {
		slog.Error("couldn't unmarshal response", "err", err)
		os.Exit(1)
	}

	timeout := time.Duration(piHoleAuthResp.Session.Validity) * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	/* Block DNS for some time  */
	if err := disableDnsBlock(ctx, piHoleAuthResp.Session.Sid, disabledDuirationSeconds); err != nil {
		slog.Error("failed to disable DNS blocking", "err", err)
		cancel()
		os.Exit(1)
	}

	slog.Info("done")

}

func disableDnsBlock(ctx context.Context, sid string, duration int) error {
	req := PiHoleDNSBlockingRequest{
		Timer:    duration,
		Blocking: false,
	}

	reqByte, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("couldn't marshal DNS blocking request: %w", err)
	}

	dnsBlockurl := fmt.Sprintf("%s/api/dns/blocking?sid=%s", DefaultPiHoleUrl, sid)
	reqWithCtx, err := http.NewRequestWithContext(ctx, http.MethodPost, dnsBlockurl, bytes.NewReader(reqByte))
	if err != nil {
		return fmt.Errorf("couldn't build DNS blocking request: %w", err)
	}
	reqWithCtx.Header.Set("Content-Type", JsonContentType)

	resp, err := http.DefaultClient.Do(reqWithCtx)
	if err != nil {
		return fmt.Errorf("couldn't send request to disable DNS blocking: %w", err)
	}
	defer resp.Body.Close()
	if !isSuccessfulStatusCode(resp.StatusCode) {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("couldn't read DNS response: %w", err)
	}

	slog.Info("done", "response", string(b))

	return nil
}

func isSuccessfulStatusCode(statusCode int) bool {
	return statusCode >= 200 && statusCode < 400
}
