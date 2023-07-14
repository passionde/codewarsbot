package store

import (
	"encoding/json"
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/models"
	"strings"
)

type ChallengeDatabase struct {
	store *Store
}

func (c *ChallengeDatabase) AddDB(challenge *models.Challenge) error {
	jsonInfo, err := json.Marshal(&challenge.Info)
	if err != nil {
		return err
	}

	_, err = c.FindByIdDB(challenge.ID)
	if err != nil {
		if strings.ToLower(err.Error()) == "sql: no rows in result set" {
			_, err = c.store.db.Exec(
				"INSERT INTO challenge (id, info, lastUpdate) VALUES ($1, $2, $3)",
				challenge.ID, jsonInfo, challenge.LastUpdate,
			)
		}
		return err
	}

	_, err = c.store.db.Exec(
		"UPDATE challenge SET info = $1, lastUpdate = $2 WHERE id = $3",
		jsonInfo, challenge.LastUpdate, challenge.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChallengeDatabase) FindByIdDB(id string) (*models.Challenge, error) {
	var res models.Challenge

	row := c.store.db.QueryRow("SELECT * FROM challenge WHERE id = $1", id)
	if row.Err() != nil {
		return &res, row.Err()
	}

	var infoStr []byte

	err := row.Scan(&res.ID, &infoStr, &res.LastUpdate)
	if err != nil {
		return &res, err
	}

	err = json.Unmarshal(infoStr, &res.Info)
	return &res, err
}

func (c *ChallengeDatabase) SelectAllRow() (map[string]models.Challenge, error) {
	challenges := make(map[string]models.Challenge)

	rows, err := c.store.db.Query("SELECT * FROM challenge")
	if err != nil {
		return challenges, err
	}
	defer rows.Close()

	for rows.Next() {
		c := models.Challenge{}
		var infoStr []byte

		err := rows.Scan(&c.ID, &infoStr, &c.LastUpdate)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(infoStr, &c.Info)
		if err != nil {
			return nil, err
		}

		challenges[c.ID] = c
	}

	return challenges, nil
}
