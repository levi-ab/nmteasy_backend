package utils

import (
	"fmt"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"math/rand"
	"nmteasy_backend/models"
	"nmteasy_backend/models/migrated_models"
	"time"
)

func ResetLeagues() {
	const USERS_PER_LIGUE = 10
	const bronzeLimit = 50
	const silverLimit = 100
	const goldLimit = 200

	var users []migrated_models.User

	models.DB.Find(&users)

	var woodenLeagues []migrated_models.League
	var bronzeLeagues []migrated_models.League
	var silverLeagues []migrated_models.League
	var goldenLeagues []migrated_models.League

	for i, user := range users {
		if user.WeeklyPoints > goldLimit {
			leaguesLen := len(goldenLeagues)
			if leaguesLen == 0 {
				goldenLeagues = append(goldenLeagues, migrated_models.League{
					ID:     uuid.New(),
					Title:  "golden",
					Points: goldLimit,
					Users:  []migrated_models.User{user},
				})

				users[i].LeagueID = &goldenLeagues[0].ID
				continue
			}

			if len(goldenLeagues[leaguesLen-1].Users) < 10 {
				goldenLeagues[leaguesLen-1].Users = append(goldenLeagues[leaguesLen-1].Users, user)
				users[i].LeagueID = &goldenLeagues[leaguesLen-1].ID
				continue
			}

			goldenLeagues = append(goldenLeagues, migrated_models.League{
				ID:     uuid.New(),
				Title:  "golden",
				Points: goldLimit,
				Users:  []migrated_models.User{user},
			})
			users[i].LeagueID = &goldenLeagues[leaguesLen-1].ID
			continue
		}

		if user.WeeklyPoints > silverLimit {
			leaguesLen := len(silverLeagues)
			if leaguesLen == 0 {
				silverLeagues = append(silverLeagues, migrated_models.League{
					ID:     uuid.New(),
					Title:  "silver",
					Points: silverLimit,
					Users:  []migrated_models.User{user},
				})

				users[i].LeagueID = &silverLeagues[0].ID
				continue
			}

			if len(silverLeagues[leaguesLen-1].Users) < 10 {
				silverLeagues[leaguesLen-1].Users = append(silverLeagues[leaguesLen-1].Users, user)
				users[i].LeagueID = &silverLeagues[leaguesLen-1].ID
				continue
			}

			silverLeagues = append(silverLeagues, migrated_models.League{
				ID:     uuid.New(),
				Title:  "silver",
				Points: silverLimit,
				Users:  []migrated_models.User{user},
			})
			users[i].LeagueID = &silverLeagues[leaguesLen-1].ID
			continue
		}

		if user.WeeklyPoints > bronzeLimit {
			leaguesLen := len(bronzeLeagues)
			if leaguesLen == 0 {
				bronzeLeagues = append(bronzeLeagues, migrated_models.League{
					ID:     uuid.New(),
					Title:  "bronze",
					Points: bronzeLimit,
					Users:  []migrated_models.User{user},
				})

				users[i].LeagueID = &bronzeLeagues[0].ID
				continue
			}

			if len(bronzeLeagues[leaguesLen-1].Users) < 10 {
				bronzeLeagues[leaguesLen-1].Users = append(bronzeLeagues[leaguesLen-1].Users, user)
				users[i].LeagueID = &bronzeLeagues[leaguesLen-1].ID
				continue
			}

			bronzeLeagues = append(bronzeLeagues, migrated_models.League{
				ID:     uuid.New(),
				Title:  "bronze",
				Points: bronzeLimit,
				Users:  []migrated_models.User{user},
			})
			users[i].LeagueID = &bronzeLeagues[leaguesLen-1].ID
			continue
		}

		leaguesLen := len(woodenLeagues)
		if leaguesLen == 0 {
			woodenLeagues = append(woodenLeagues, migrated_models.League{
				ID:     uuid.New(),
				Title:  "wooden",
				Points: 0,
				Users:  []migrated_models.User{user},
			})

			users[i].LeagueID = &woodenLeagues[0].ID
			continue
		}

		if len(woodenLeagues[leaguesLen-1].Users) < 10 {
			woodenLeagues[leaguesLen-1].Users = append(woodenLeagues[leaguesLen-1].Users, user)
			users[i].LeagueID = &woodenLeagues[leaguesLen-1].ID
			continue
		}

		woodenLeagues = append(woodenLeagues, migrated_models.League{
			ID:     uuid.New(),
			Title:  "wooden",
			Points: 0,
			Users:  []migrated_models.User{user},
		})
		users[i].LeagueID = &woodenLeagues[leaguesLen-1].ID
		continue

	}

	models.DB.Save(&users)

	models.DB.Exec("DELETE FROM leagues")

	models.DB.Save(&woodenLeagues)
	models.DB.Save(&bronzeLeagues)
	models.DB.Save(&silverLeagues)
	models.DB.Save(&goldenLeagues)

}

func removeUser(users []migrated_models.User, target migrated_models.User) []migrated_models.User {
	var result []migrated_models.User
	for _, user := range users {
		if user.ID != target.ID {
			result = append(result, user)
		}
	}
	return result
}

func GenerateRandomUsers(db *gorm.DB) error {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	for i := 1; i <= 100; i++ {
		user := migrated_models.User{
			ID:           uuid.New(),
			FirstName:    fmt.Sprintf("Bot%d", i),
			LastName:     fmt.Sprintf("Bot%d", i),
			Username:     fmt.Sprintf("username%d", i),
			Email:        fmt.Sprintf("user%d@bot.com", i),
			Points:       rand.Intn(500),
			WeeklyPoints: rand.Intn(400),
			Password:     fmt.Sprintf("password%d", i),
			LeagueID:     nil, // Set league ID to nil for now
		}

		// Create user record in the database
		if err := db.Create(&user).Error; err != nil {
			return err
		}
	}

	return nil
}
