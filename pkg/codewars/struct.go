package codewars

import "time"

// User curl https://www.codewars.com/api/v1/users/some_user
type User struct {
	Username            string   `json:"username"`
	Name                string   `json:"name"`
	Honor               int      `json:"honor"`
	Clan                string   `json:"clan"`
	LeaderboardPosition int      `json:"leaderboardPosition"`
	Skills              []string `json:"skills"`
	Ranks               struct {
		Overall   *UserRank            `json:"overall"`
		Languages map[string]*UserRank `json:"languages"`
	} `json:"ranks"`
	CodeChallenges struct {
		TotalAuthored        int `json:"totalAuthored"`
		TotalCompletedUnique int `json:"totalCompleted"`
		TotalCompletedAll    int `json:"totalCompletedAll"`
	} `json:"codeChallenges"`
	Reason string `json:"reason"`
}

type UserRank struct {
	Rank  int    `json:"rank"`
	Name  string `json:"name"`
	Color string `json:"color"`
	Score int    `json:"score"`
}

// Kata curl https://www.codewars.com/api/v1/code-challenges/valid-braces
type Kata struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Slg         string     `json:"slug"`
	Category    string     `json:"category"`
	PublishedAt *time.Time `json:"publishedAt"`
	ApprovedAt  *time.Time `json:"approvedAt"`
	Languages   []string   `json:"languages"`
	URL         string     `json:"url"`
	Rank        struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Color string `json:"color"`
	} `json:"rank"`
	CreatedAt          *time.Time    `json:"createdAt"`
	CreatedBy          *KataUserData `json:"createdBy"`
	ApprovedBy         *KataUserData `json:"approvedBy"`
	Description        string        `json:"description"`
	TotalAttempts      int           `json:"totalAttempts"`
	TotalCompleted     int           `json:"totalCompleted"`
	TotalStars         int           `json:"totalStars"`
	VoteScore          int           `json:"voteScore"`
	Tags               []string      `json:"tags"`
	ContributorsWanted bool          `json:"contributorsWanted"`
	Unresolved         struct {
		Issues      int `json:"issues"`
		Suggestions int `json:"suggestions"`
	} `json:"unresolved"`
	Reason string `json:"reason"`
}

type KataUserData struct {
	Username string `json:"username"`
	URL      string `json:"url"`
}

// CodeChallengesCompleted curl http://www.codewars.com/api/v1/users/some_user/code-challenges/completed?page=0
type CodeChallengesCompleted struct {
	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
	Data       []struct {
		Id                 string    `json:"id"`
		Name               string    `json:"name"`
		Slug               string    `json:"slug"`
		CompletedAt        time.Time `json:"completedAt"`
		CompletedLanguages []string  `json:"completedLanguages"`
	} `json:"data"`
	Reason string `json:"reason"`
}

// CodeChallengesAuthored curl http://www.codewars.com/api/v1/users/some_user/code-challenges/completed?page=0
type CodeChallengesAuthored struct {
	Data []struct {
		Id          string   `json:"id"`
		Name        string   `json:"name"`
		Description string   `json:"description"`
		Rank        int      `json:"rank"`
		RankName    string   `json:"rankName"`
		Tags        []string `json:"tags"`
		Languages   []string `json:"languages"`
	} `json:"data"`
	Reason string `json:"reason"`
}
