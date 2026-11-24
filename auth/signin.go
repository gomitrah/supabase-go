package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
)

// SignInWithPasswordByEmail signs a user in with their email address and
// password. The result carries their session and, when the server sent one,
// its advisory that the password falls short of the project's current
// requirements.
//
// The server's rejections are returned as [*Error], with a Code such as:
//   - invalid_credentials, alike for a wrong password, an unknown address and
//     an account that can only sign in another way
//   - email_not_confirmed, for an address not yet confirmed
func (c *Client) SignInWithPasswordByEmail(ctx context.Context, email, password string) (*SignInWithPasswordResult, error) {
	return c.signInWithPassword(ctx, passwordSignInPayload{Email: email, Password: password})
}

// SignInWithPasswordByPhone signs a user in with their phone number and
// password. The result carries their session and, when the server sent one,
// its advisory that the password falls short of the project's current
// requirements.
//
// The server's rejections are returned as [*Error], with a Code such as:
//   - invalid_credentials, alike for a wrong password, an unknown number and
//     an account that can only sign in another way
//   - phone_not_confirmed, for a number not yet confirmed
func (c *Client) SignInWithPasswordByPhone(ctx context.Context, phone, password string) (*SignInWithPasswordResult, error) {
	return c.signInWithPassword(ctx, passwordSignInPayload{Phone: phone, Password: password})
}

// signInWithPassword sends body to the token endpoint as a password grant and
// shapes the answer into a result with the session and, when the server sent
// one, its weak-password advisory.
func (c *Client) signInWithPassword(ctx context.Context, body passwordSignInPayload) (*SignInWithPasswordResult, error) {
	endpoint := c.baseURL.JoinPath("token")
	endpoint.RawQuery = url.Values{"grant_type": {"password"}}.Encode()
	response, err := send(ctx, c.httpClient, http.MethodPost, endpoint, body, "")
	if err != nil {
		return nil, err
	}
	state, err := parseTokenResponse(response)
	if err != nil {
		return nil, err
	}
	var advisory struct {
		WeakPassword *struct {
			Message string   `json:"message"`
			Reasons []string `json:"reasons"`
		} `json:"weak_password"`
	}
	if err := json.Unmarshal(response, &advisory); err != nil {
		return nil, fmt.Errorf("auth: decoding token response: %w", err)
	}
	result := &SignInWithPasswordResult{session: c.newSession(state, nil)}
	if advisory.WeakPassword != nil {
		result.weakPassword = &WeakPassword{message: advisory.WeakPassword.Message, reasons: advisory.WeakPassword.Reasons}
	}
	return result, nil
}

// passwordSignInPayload is the wire body of a password sign-in. password is
// always stated, since the server requires it.
type passwordSignInPayload struct {
	Email    string `json:"email,omitzero"`
	Phone    string `json:"phone,omitzero"`
	Password string `json:"password"`
}

// SignInWithPasswordResult is the answer to
// [Client.SignInWithPasswordByEmail] and [Client.SignInWithPasswordByPhone]:
// the new [Session] and, when the server sent one, the advisory that the
// accepted password falls short of the project's current requirements.
type SignInWithPasswordResult struct {
	session      *Session
	weakPassword *WeakPassword
}

// Session returns the signed-in user's session. It is never nil.
func (r *SignInWithPasswordResult) Session() *Session { return r.session }

// WeakPassword returns the server's advisory that the accepted password falls
// short of the project's current password requirements, or nil when the
// server sent none.
func (r *SignInWithPasswordResult) WeakPassword() *WeakPassword { return r.weakPassword }

// WeakPassword is the server's advisory, sent beside a successful sign-in,
// that the accepted password falls short of the project's current password
// requirements. A WeakPassword is immutable and safe for concurrent use by
// multiple goroutines.
type WeakPassword struct {
	message string
	reasons []string
}

// Message returns the server's description of the shortfall.
func (w *WeakPassword) Message() string { return w.message }

// Reasons returns the server's codes for each requirement the password fails,
// such as "length" or "characters". The result is a copy the caller may
// retain and modify freely.
func (w *WeakPassword) Reasons() []string { return slices.Clone(w.reasons) }
