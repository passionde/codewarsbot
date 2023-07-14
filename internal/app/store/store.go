package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"time"
)

type Store struct {
	config    *Config
	db        *sql.DB
	challenge *ChallengeRepo
}

func New(config *Config) *Store {
	return &Store{
		config: config,
	}
}

func (s *Store) Open() error {
	db, err := sql.Open("postgres", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.db = db

	return nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Challenge() *ChallengeRepo {
	if s.challenge != nil {
		return s.challenge
	}
	//repo, err := NewChallengeRepo(3*24*time.Hour, 6*time.Hour, s)
	repo, err := NewChallengeRepo(3*24*time.Hour, 6*time.Hour, s)
	if err != nil {
		log.Fatal(err)
	}

	s.challenge = repo

	return s.challenge

}
