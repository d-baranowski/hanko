package persistence

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gobuffalo/pop/v6"
	"github.com/teamhanko/hanko/backend/persistence/models"
)

type OauthStatePersister interface {
	Create(state string) error
	Get(state string) (*models.OauthState, error)
	Delete(state string) error
}

type oauthStatePersister struct {
	db *pop.Connection
}

func NewOauthStatePersister(db *pop.Connection) OauthStatePersister {
	return &oauthStatePersister{db: db}
}

func (p *oauthStatePersister) Get(state string) (*models.OauthState, error) {
	m := []*models.OauthState{}
	err := p.db.Store.Select(&m, "SELECT * FROM oauth_state WHERE state = $1", state)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth state: %w", err)
	}

	return m[0], nil
}

func (p *oauthStatePersister) Create(state string) error {
	_, err := p.db.Store.Exec("INSERT INTO oauth_state (state, created_at, updated_at) VALUES ($1, now(), now())", state)
	if err != nil {
		return fmt.Errorf("failed to store oauth state: %w", err)
	}

	return nil
}

func (p *oauthStatePersister) Delete(state string) error {
	result, err := p.db.Store.Exec("DELETE FROM oauth_state WHERE state = ?", state)
	if err != nil {
		return fmt.Errorf("failed to verify oauth_state: %w", err)
	}

	n, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to verify oauth_state: %w", err)
	}

	if n < 1 {
		return fmt.Errorf("failed to verify oauth_state as this state was not present in store")
	}

	return nil
}
