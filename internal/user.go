package internal

import (
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"secretary/alpha/internal/config"
	"secretary/alpha/storage"
	"secretary/alpha/utils"
)

type User struct {
	UUID         string
	Username     string
	PasswordHash string
	Active       bool
	CreatedTime  string
	ModifiedTime string
}

func (u *User) SetPassword(password string) error {
	// Use configurable bcrypt cost from config
	cost := 12 // Default cost
	if config.GlobalConfig != nil {
		cost = config.GlobalConfig.Security.BcryptCost
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	u.PasswordHash = string(hash)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
	return err == nil
}

func (u *User) CreateUser(username string, password string, active bool) error {
	existingUser := u.GetUser(username)
	if existingUser != nil {
		return fmt.Errorf("username %v already exists", username)
	}

	// FIXME Add validation code here ...
	// FIXME change the error handling

	createdTime := utils.CurrentTime()

	u.UUID = utils.UUID()
	u.Username = username
	u.CreatedTime = createdTime
	u.ModifiedTime = createdTime
	u.Active = active

	err := u.SetPassword(password)
	if err != nil {
		return fmt.Errorf("setpassword error: %v", err)
	}

	query := `
		INSERT INTO user_local (uuid, username, password, active, created_time, modified_time)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	_, err = storage.DatabaseExec(query, u.UUID, u.Username, u.PasswordHash, u.Active, u.CreatedTime, u.ModifiedTime)
	if err != nil {
		return fmt.Errorf("error in createuser: %v", err)
	}

	return nil
}

func (u *User) GetUser(username string) *User {
	// Use parameterized query to prevent SQL injection
	query := `SELECT uuid, username, password, active, created_time, modified_time FROM user_local WHERE username = ?`

	rows, err := storage.DatabaseQuery(query, username)
	if err != nil {
		utils.Logger("err", fmt.Sprintf("Error querying user: %v", err))
		return nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		utils.Logger("err", err.Error())
		return nil
	}

	results, err := utils.HandleTableToJSON(columns, rows)
	if err != nil {
		utils.Logger("err", err.Error())
		return nil
	}

	if len(results) == 0 {
		return nil
	}

	return &User{
		UUID:         results[0]["uuid"].(string),
		Username:     results[0]["username"].(string),
		PasswordHash: results[0]["password"].(string),
		Active:       results[0]["active"].(bool),
		CreatedTime:  results[0]["created_time"].(time.Time).Format(time.RFC3339),
		ModifiedTime: results[0]["modified_time"].(time.Time).Format(time.RFC3339),
	}
}

func (u *User) GetAllUsers() []*User {
	query := `SELECT * FROM user_local`

	rows, err := storage.DatabaseQuery(query)
	if err != nil {
		utils.Logger("err", err.Error())
		return nil
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		utils.Logger("err", err.Error())
		return nil
	}

	results, err := utils.HandleTableToJSON(columns, rows)
	if err != nil {
		utils.Logger("err", err.Error())
		return nil
	}

	users := make([]*User, 0, len(results))
	for _, res := range results {
		user := &User{
			UUID:         res["uuid"].(string),
			Username:     res["username"].(string),
			PasswordHash: res["password"].(string),
			Active:       res["active"].(bool),
			CreatedTime:  res["created_time"].(time.Time).Format(time.RFC3339),
			ModifiedTime: res["modified_time"].(time.Time).Format(time.RFC3339),
		}
		users = append(users, user)
	}

	return users
}
