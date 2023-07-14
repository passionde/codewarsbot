package codewars

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "https://www.codewars.com/api/v1/"

func GetUser(username string) (User, error) {
	var response User

	body, err := get(fmt.Sprintf("%susers/%s", baseURL, username))
	if err != nil {
		return response, err
	}

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, errors.New(string(body))
	}

	if response.Reason != "" {
		return response, errors.New(response.Reason)
	}

	return response, nil
}

func GetCompletedChallenges(username string, page int) (CodeChallengesCompleted, error) {
	var challenges CodeChallengesCompleted

	body, err := get(fmt.Sprintf("%susers/%s/code-challenges/completed?page=%d", baseURL, username, page))
	if err != nil {
		return challenges, err
	}

	err = json.Unmarshal(body, &challenges)
	if err != nil {
		return challenges, errors.New(string(body))
	}

	if challenges.Reason != "" {
		return challenges, errors.New(challenges.Reason)
	}

	return challenges, nil
}

func GetKata(id string) (Kata, error) {
	var challenge Kata

	body, err := get(fmt.Sprintf("%scode-challenges/%s", baseURL, id))
	if err != nil {
		return challenge, err
	}

	err = json.Unmarshal(body, &challenge)
	if err != nil {
		return challenge, errors.New(string(body))
	}

	if challenge.Reason != "" {
		return challenge, errors.New(challenge.Reason)
	}

	return challenge, nil

}

func GetAuthoredChallenges(username string) (CodeChallengesAuthored, error) {
	var challenges CodeChallengesAuthored

	body, err := get(fmt.Sprintf("%susers/%s/code-challenges/authored", baseURL, username))
	if err != nil {
		return challenges, err
	}

	err = json.Unmarshal(body, &challenges)
	if err != nil {
		return challenges, errors.New(string(body))
	}

	if challenges.Reason != "" {
		return challenges, errors.New(challenges.Reason)
	}

	return challenges, nil
}

func get(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 && resp.StatusCode > 299 {
		return nil, errors.New(fmt.Sprintf("%d", resp.StatusCode))
	}

	return body, nil
}
