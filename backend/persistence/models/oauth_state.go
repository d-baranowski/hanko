package models

import "time"

/*
Used to store the state during oauth authentication with third party providers.
When the code is returned by the third party provider state is fetched from the database to verify it
*/
type OauthState struct {
	State     string    `db:"state" json:"state"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
	UpdatedAt time.Time `db:"updated_at" json:"updated_at"`
}
