package handler

import (
	"encoding/json"
	"gorm.io/datatypes"
)

/*
CREATE TABLE user_sessions (

	session_id VARCHAR(255) PRIMARY KEY,
	user_id INT NOT NULL,
	device_fingerprint VARCHAR(255),
	platform VARCHAR(50),
	ip_address VARCHAR(45),
	user_agent TEXT,
	login_time TIMESTAMP NOT NULL,
	expiration_time TIMESTAMP,
	session_token TEXT NOT NULL,
	last_activity_time TIMESTAMP,
	session_status VARCHAR(50),
	geo_location TEXT,
	metadata JSON

);
*/
type SessionInfo struct {
	SessionStatus string              `json:"sessionStatus"`
	RequestHeader map[string][]string `json:"requestHeader"`
	Fingerprint   string              `json:"fingerprint"`
	Platform      string              `json:"platform"`
}

func (si *SessionInfo) Serialize() []byte {
	rawData, _ := json.Marshal(si)
	return rawData
}

type UserSession struct {
	RecordId    int64          `json:"recordId" gorm:"primary_key;auto_increment;not_null"`
	UserId      int            `json:"userId"`
	SessionTime int64          `json:"sessionTime"` //epoch with millisecons
	SessionId   string         `json:"sessionId"`
	SessionInfo datatypes.JSON `json:"sessionInfo"`
}

func (us *UserSession) GetSessionInfo() *SessionInfo {
	sessionInfo := SessionInfo{}
	json.Unmarshal(us.SessionInfo, &sessionInfo)
	return &sessionInfo
}
