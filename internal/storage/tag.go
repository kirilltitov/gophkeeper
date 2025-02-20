package storage

import (
	"context"

	"github.com/google/uuid"
)

type Tags []string

func (s *PgSQL) AddTag(ctx context.Context, secretID uuid.UUID, tag string) error {
	query := `
		insert into public.tag (secret_id, text)
		values ($1, $2)
		on conflict (secret_id, text) do update set text = excluded.text
	`
	_, err := s.Conn.Exec(ctx, query, secretID, tag)
	if err != nil {
		return err
	}

	return nil
}

func (s *PgSQL) DeleteTag(ctx context.Context, secretID uuid.UUID, tag string) error {
	query := `delete from public.tag where secret_id = $1 and text = $2`
	_, err := s.Conn.Exec(ctx, query, secretID, tag)
	if err != nil {
		return err
	}

	return nil
}
