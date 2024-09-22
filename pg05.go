package pg05

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

type User struct {
	Id          int
	UserName    string
	Name        string
	SurName     string
	Description string
}

var (
	Hostname = ""
	Port     = 5432
	Username = ""
	Password = ""
	Database = ""
)

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", Hostname, Port, Username, Password, Database)

	db, err := sql.Open("postgres", conn)

	if err != nil {
		return nil, err
	}
	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userId := -1
	statement := fmt.Sprintf(`SELECT id FROM users WHERE username='%s'`, username)

	rows, err := db.Query(statement)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	for rows.Next() {
		var ID int
		if err := rows.Scan(&ID); err != nil {
			fmt.Println("Scan", err)
			return -1
		}
		userId = ID
	}
	defer rows.Close()
	return userId
}

func AddUser(d User) int {
	d.UserName = strings.ToLower(d.UserName)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userId := exists(d.UserName)
	if userId == -1 {
		fmt.Println("user already exists", d.UserName)
		return -1
	}

	insertStatement := `INSERT INTO "users" ("username") VALUES ($1)`

	_, err = db.Exec(insertStatement, d.UserName)
	if err != nil {
		fmt.Println(err)
		return -1
	}

	userId = exists(d.UserName)
	if userId == -1 {
		return userId
	}

	insertStatement = `INSERT INTO "user_data" ("user_id", "name", "surname", "description") VALUES ($1, $2, $3, $4)`

	_, err = db.Exec(insertStatement, userId, d.Name, d.SurName, d.Description)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	return userId
}

func DeleteUser(id int) error {
	db, err := openConnection()

	if err != nil {
		return err
	}
	defer db.Close()

	statement := fmt.Sprintf(`SELECT "username" FROM "users" WHERE id = "%s"`, id)

	rows, err := db.Query(statement)
	if err != nil {
		return err
	}
	defer rows.Close()
	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}

	// лишний раз проверяем
	if id != exists(username) {
		return fmt.Errorf("user with id %d does not exists", id)
	}

	statement = `DELETE FROM "users" WHERE id=$1`
	_, err = db.Exec(statement, id)
	if err != nil {
		return err
	}

	statement = `DELETE FROM "user_data" WHERE id=$1`
	_, err = db.Exec(statement, id)
	if err != nil {
		return err
	}

	return nil
}

func List() ([]User, error) {
	Data := []User{}

	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "id", "username", "name", "surname", "description" FROM "users", "user_data" WHERE users.id=user_data.user_id`)
	if err != nil {
		return Data, err
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string

		err = rows.Scan(&id, &username, &name, &surname, &description)
		if err != nil {
			return Data, err
		}

		Data = append(Data, User{
			Id:          id,
			UserName:    username,
			Name:        name,
			SurName:     surname,
			Description: description,
		})
	}

	return Data, nil

}

func Update(u User) error {
	u.UserName = strings.ToLower(u.UserName)

	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userId := exists(u.UserName)
	if userId == -1 {
		return errors.New("User does not exist")
	}

	statement := `UPDATE "user_data" SET "name"=$1, "surname"=$2, "description"$3 where "user_id"=%4`

	_, err = db.Exec(statement, u.Name, u.SurName, u.Description, userId)
	if err != nil {
		return err
	}
	return nil
}
