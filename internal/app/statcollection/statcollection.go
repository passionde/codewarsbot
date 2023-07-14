package statcollection

import (
	"github.com/Yarik-xxx/CodeWarsRestApi/internal/app/store"
	"github.com/Yarik-xxx/CodeWarsRestApi/pkg/codewars"
	"log"
	"strings"
	"time"
)

type StatCollection struct {
	challengeRepo *store.ChallengeRepo
}

type UserInfo struct {
	User          codewars.User                   `json:"user"`
	Authored      codewars.CodeChallengesAuthored `json:"authored"`
	StatisticsKyu map[string]InfoKata             `json:"statisticsKyu"`
}

type InfoKata struct {
	TotalCompleted int            `json:"TotalCompleted"`
	Rank           string         `json:"rank"`
	Score          int            `json:"score"`
	ByRank         map[string]int `json:"byRank"`
	ByTags         map[string]int `json:"byTags"`
}

func New(st *store.Store) *StatCollection {
	return &StatCollection{challengeRepo: st.Challenge()}
}

func (s *StatCollection) AllInfo(username string) (*UserInfo, error) {
	result := new(UserInfo)

	// Сбор основной информации об аккаунте
	user, err := s.GetUserInfo(username)
	if err != nil {
		return result, err
	}
	result.User = user

	// Сбор информации об составленных катах
	listAuthor, err := codewars.GetAuthoredChallenges(username)
	if err != nil {
		return result, err
	}
	result.Authored = listAuthor

	// Сбор информации о выполненных катах
	listCompleted, err := s.getCompletedList(user.Username)
	if err != nil {
		return nil, err
	}

	listKataInfo, err := s.getRankList(listCompleted)
	if err != nil {
		return nil, err
	}

	result.StatisticsKyu = s.count(listCompleted, listKataInfo, &user)
	result.User.CodeChallenges.TotalCompletedAll = result.StatisticsKyu["overall"].TotalCompleted

	return result, nil

}

func (s *StatCollection) GetUserInfo(username string) (codewars.User, error) {
	user, err := codewars.GetUser(username)

	for err != nil && (err.Error() == "429" || strings.TrimSpace(strings.ToLower(err.Error())) == "retry later") {
		time.Sleep(5 * time.Second)
		user, err = codewars.GetUser(username)
	}

	if err != nil {
		return user, err
	}

	return user, nil
}

func (s *StatCollection) GetAuthoredInfo(username string) (codewars.CodeChallengesAuthored, error) {
	list, err := codewars.GetAuthoredChallenges(username)

	for err != nil && (err.Error() == "429" || strings.TrimSpace(strings.ToLower(err.Error())) == "retry later") {
		log.Println("Спать")
		time.Sleep(5 * time.Second)
		list, err = codewars.GetAuthoredChallenges(username)
	}

	if err != nil {
		return list, err
	}

	return list, nil
}

func (s *StatCollection) getCompletedList(username string) (map[string]map[string]interface{}, error) {
	// {"python": {"123sd12312":nil, "1qda4f1313":nil}}
	result := make(map[string]map[string]interface{})

	totalPage := 1
	for page := 0; page < totalPage; page++ {
		completed, err := codewars.GetCompletedChallenges(username, page)

		for err != nil && (err.Error() == "429" || strings.TrimSpace(strings.ToLower(err.Error())) == "retry later") {
			time.Sleep(5 * time.Second)
			completed, err = codewars.GetCompletedChallenges(username, page)
		}
		if err != nil {
			return result, err
		}

		totalPage = completed.TotalPages

		for _, info := range completed.Data {
			// Перебор всех языков на котором выполнена ката
			for _, lang := range info.CompletedLanguages {
				if _, ok := result[lang]; !ok {
					result[lang] = make(map[string]interface{})
				}
				result[lang][info.Id] = nil
			}
		}
	}
	return result, nil
}

func (s *StatCollection) getRankList(list map[string]map[string]interface{}) (map[string]codewars.Kata, error) {
	result := make(map[string]codewars.Kata)

	for _, data := range list {
		for id, _ := range data {
			result[id] = codewars.Kata{}
		}
	}

	//todo можно ускорить при необходимости
	for id, _ := range result {
		kata, _, err := s.challengeRepo.Get(id)

		for err != nil && (err.Error() == "429" || strings.TrimSpace(strings.ToLower(err.Error())) == "retry later") {
			time.Sleep(5 * time.Second)
			kata, _, err = s.challengeRepo.Get(id)
		}
		if err != nil {
			log.Println(err)
			continue
		}

		result[id] = kata
	}
	return result, nil
}

func (s *StatCollection) count(list map[string]map[string]interface{}, listKata map[string]codewars.Kata, user *codewars.User) map[string]InfoKata {
	result := make(map[string]InfoKata)

	result["overall"] = InfoKata{
		TotalCompleted: 0,
		Rank:           user.Ranks.Overall.Name,
		Score:          user.Ranks.Overall.Score,
		ByRank:         make(map[string]int),
		ByTags:         make(map[string]int),
	}

	//todo Какой-то прикол, смог только так
	total := 0
	for lang, ls := range list {
		if _, ok := result[lang]; !ok {
			if user.Ranks.Languages[lang] == nil {
				user.Ranks.Languages[lang] = &codewars.UserRank{
					Rank:  0,
					Name:  "Unknown",
					Color: "Unknown",
					Score: 0,
				}
			}

			result[lang] = InfoKata{
				TotalCompleted: len(ls),
				Rank:           user.Ranks.Languages[lang].Name,
				Score:          user.Ranks.Languages[lang].Score,
				ByRank:         make(map[string]int),
				ByTags:         make(map[string]int),
			}
		}

		for id, _ := range ls {
			if _, ok := listKata[id]; !ok {
				continue
			}

			result[lang].ByRank[listKata[id].Rank.Name]++
			result["overall"].ByRank[listKata[id].Rank.Name]++
			//todo Какой-то прикол, смог только так
			total++

			for _, tag := range listKata[id].Tags {
				result[lang].ByTags[tag]++
				result["overall"].ByTags[tag]++
			}
		}
	}

	//todo Какой-то прикол, смог только так
	result["overall"] = InfoKata{
		TotalCompleted: total,
		Rank:           result["overall"].Rank,
		Score:          result["overall"].Score,
		ByRank:         result["overall"].ByRank,
		ByTags:         result["overall"].ByTags,
	}

	return result
}
