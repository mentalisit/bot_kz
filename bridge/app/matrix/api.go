package matrix

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

// apiCall: Outgoing helper using raw HTTP (No SDK)
// userID: The Matrix ID to act as (e.g. @bridge_bot:...)
func (m *Matrix) apiCall(method, path, userID string, body interface{}) ([]byte, error) {
	url := fmt.Sprintf("%s%s?access_token=%s&user_id=%s",
		m.Config.Matrix.HomeserverURL, path, m.Config.Matrix.ASToken, userID)

	var bodyReader io.Reader
	if body != nil {
		b, _ := json.Marshal(body)
		bodyReader = bytes.NewBuffer(b)
	}

	req, _ := http.NewRequest(method, url, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 300 {
		return respBody, fmt.Errorf("matrix api error: %d %s", resp.StatusCode, string(respBody))
	}

	return respBody, nil
}

func (m *Matrix) GetOrCreateGhost(localpart, displayName, avatarURL string) string {
	mxid := fmt.Sprintf("@%s:%s", localpart, m.getHomeserverDomain())

	// Register the ghost
	m.Register(mxid)

	// Check cached profile to avoid unnecessary API calls
	m.mu.RLock()
	cached, exists := m.ProfileCache[mxid]
	m.mu.RUnlock()

	needUpdate := false

	// Update displayname only if changed
	if displayName != "" && displayName != cached.DisplayName {
		m.SetDisplayname(mxid, displayName)
		cached.DisplayName = displayName
		needUpdate = true
	}

	// Update avatar only if changed
	if avatarURL != "" && avatarURL != cached.AvatarURL {
		m.SetAvatar(mxid, avatarURL)
		cached.AvatarURL = avatarURL
		needUpdate = true
	}

	// Update cache if needed
	if needUpdate || !exists {
		m.mu.Lock()
		m.ProfileCache[mxid] = cached
		m.mu.Unlock()
	}

	return mxid
}

func (m *Matrix) Register(userID string) {
	// Extract localpart from @name:domain
	localpart := userID
	if strings.HasPrefix(userID, "@") {
		parts := strings.Split(strings.TrimPrefix(userID, "@"), ":")
		localpart = parts[0]
	}

	payload := map[string]interface{}{
		"username": localpart,
		"type":     "m.login.application_service",
	}

	// Use an empty user_id for registration itself or the bot ID
	_, err := m.apiCall("POST", "/_matrix/client/v3/register", userID, payload)
	if err != nil {
		if !strings.Contains(err.Error(), "M_USER_IN_USE") {
			log.Printf("Registration for %s failed: %v", userID, err)
		}
	} else {
		log.Printf("Successfully registered/activated user %s", userID)
	}
}

func (m *Matrix) getHomeserverDomain() string {
	domain := m.Config.Matrix.Username
	if strings.Contains(domain, ":") {
		return strings.Split(domain, ":")[1]
	}
	return "matrix.org"
}

func (m *Matrix) SetDisplayname(userID, displayname string) {
	payload := map[string]string{"displayname": displayname}
	_, err := m.apiCall("PUT", "/_matrix/client/v3/profile/"+userID+"/displayname", userID, payload)
	if err != nil {
		log.Printf("Failed to set displayname for %s: %v", userID, err)
	}
}

func (m *Matrix) SetAvatar(userID, avatarURL string) {
	if avatarURL == "" {
		return
	}

	mxcURL := avatarURL
	// If it's an HTTP URL, upload it to Matrix first
	if strings.HasPrefix(avatarURL, "http") {
		mxc, _, err := m.UploadMediaFromURL(avatarURL, "avatar")
		if err != nil {
			log.Printf("Failed to upload avatar from %s: %v", avatarURL, err)
			return
		}
		mxcURL = mxc
	}

	payload := map[string]string{"avatar_url": mxcURL}
	_, err := m.apiCall("PUT", "/_matrix/client/v3/profile/"+userID+"/avatar_url", userID, payload)
	if err != nil {
		log.Printf("Failed to set avatar for %s: %v", userID, err)
	}
}
