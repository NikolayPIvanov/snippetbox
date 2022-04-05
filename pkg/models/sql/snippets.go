package sql

import (
	"database/sql"
	"errors"
	"fmt"

	"NikolayPIvanov/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := fmt.Sprintf(`INSERT INTO snippets (title, content, created, expires)
	VALUES('%s', '%s', GETUTCDATE(), DATEADD(day, %s, GETUTCDATE())); select ID = convert(bigint, SCOPE_IDENTITY())`,
		title, content, expires)

	id := 0
	m.DB.QueryRow(stmt, title, content, expires).Scan(&id)
	if id == 0 {
		return id, errors.New("could not insert snippet")
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := fmt.Sprintf(`SELECT id, title, content, created, expires FROM snippets
	WHERE expires > GETUTCDATE() AND id = '%d'`, id)

	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT TOP 10 id, title, content, created, expires FROM snippets
			 WHERE expires > GETUTCDATE() ORDER BY created DESC`
	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	snippets := []*models.Snippet{}
	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
