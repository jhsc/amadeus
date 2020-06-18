package sqlite

import (
	"database/sql"
	"time"

	"gitlab.com/jhsc/amadeus/store"
)

type projectStore struct {
	db *sql.DB
}

// New creates a new user.
func (s *projectStore) New(projectName string) (int64, error) {
	res, err := s.db.Exec(
		`insert into projects(created_at, name) values(?, ?)`,
		time.Now(), projectName,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

const selectFromProjects = `
	select
		id,
		coalesce(name, '') as name,
		created_at,
		error,
		notes
	from project
`

func (s *projectStore) scanProject(scanner scanner) (*store.Project, error) {
	p := new(store.Project)
	err := scanner.Scan(
		&p.ID,
		&p.Name,
		&p.CreatedAt,
		&p.Error,
		&p.Notes,
	)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return p, nil
}

// Get finds a project by ID.
func (s *projectStore) Get(id int64) (*store.Project, error) {
	row := s.db.QueryRow(selectFromProjects+` where id=?`, id)
	return s.scanProject(row)
}

// GetMany finds projects by IDs.
func (s *projectStore) GetMany(ids []int64) (map[int64]*store.Project, error) {
	if len(ids) == 0 {
		return make(map[int64]*store.Project), nil
	}

	projects := make(map[int64]*store.Project)
	for _, id := range ids {
		projects[id] = nil
	}

	var params []interface{}
	for id := range projects {
		params = append(params, id)
	}

	rows, err := s.db.Query(
		selectFromProjects+` where id in (`+placeholders(len(params))+`)`,
		params...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		user, err := s.scanProject(rows)
		if err != nil {
			return nil, err
		}
		projects[user.ID] = user
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, user := range projects {
		if user == nil {
			return nil, store.ErrNotFound
		}
	}
	return projects, nil
}
