package post05

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
)

// Connection details
var (
	Hostname = ""
	Port     = 2345
	Username = ""
	Password = ""
	Database = ""
)

// CourseData is for holding full user data
// CourseData table + Username
type CourseData struct {
	ID     int
	Cid    string
	Name   string
	PreReq string
}

func openConnection() (*sql.DB, error) {
	// connection string
	conn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname, Port, Username, Password, Database)

	// open database
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// The function returns the ID of the Cid
// -1 if the user does not exist
func exists(cid string) int {
	//username = strings.ToLower(username)
	fmt.Printf("%s exists? \n", cid)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	ID := -1
	statement := fmt.Sprintf(`SELECT "id" FROM "classdata" where cid = '%s'`, cid)
	rows, err := db.Query(statement)
	//fmt.Println("rows found:", rows)
	if rows != nil {
		for rows.Next() {
			var id int
			err = rows.Scan(&id)
			if err != nil {
				fmt.Println("Scan", err)
				return -1
			}
			ID = id
		}
		defer rows.Close()
	}
	return ID
}

// AddUser adds a new user to the database
// Returns new User ID
// -1 if there was an error
func AddUser(d CourseData) int {
	//d.Username = strings.ToLower(d.Username)
	fmt.Printf("adding user %s \n", d.Cid)
	db, err := openConnection() //starts a connection to the database as db
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	id := exists(d.Cid) //ID happens to already exist
	if id != -1 {
		fmt.Println("User already exists:", id)
		return -1
	}

	//execute a network command the updated the full data fields in classdata
	insertStatement := `insert into "classdata" ("cid", "name", "prereq")
	values ($1, $2, $3)`
	_, err = db.Exec(insertStatement, d.Cid, d.Name, d.PreReq)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}
	id = exists(d.Cid) //get the id

	return id //if we run without errors, return the new user id
}

// DeleteUser deletes an existing user
func DeleteUser(id int) error {
	fmt.Printf("deleting user: %d \n", id)
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	// Does the ID exist?
	statement := fmt.Sprintf(`SELECT "cid" FROM "classdata" where id = %d`, id)
	rows, err := db.Query(statement)

	var cid string
	if rows != nil {
		for rows.Next() {
			err = rows.Scan(&cid)
			if err != nil {
				return err
			}
		}
	}
	defer rows.Close()

	if exists(cid) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	// Delete from CourseData
	deleteStatement := `delete from "classdata" where id=$1`
	_, err = db.Exec(deleteStatement, id)
	if err != nil {
		return err
	}

	return nil
}

// ListUsers lists all users in the database
func ListUsers() ([]CourseData, error) {
	fmt.Printf("listing users \n")
	Data := []CourseData{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT "id","cid","name","prereq"
		FROM "classdata";`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var cid string
		var Name string
		var PreReq string
		err = rows.Scan(&id, &cid, &Name, &PreReq)
		temp := CourseData{ID: id, Cid: cid, Name: Name, PreReq: PreReq}
		Data = append(Data, temp)
		if err != nil {
			return Data, err
		}
	}
	defer rows.Close()
	return Data, nil
}

// UpdateUser is for updating an existing user
func UpdateUser(d CourseData) error {
	fmt.Printf("updating user %s \n", d.Cid)
	db, err := openConnection()
	if err != nil {
		return err
	}
	defer db.Close()

	userID := exists(d.Cid)
	if userID == -1 {
		return errors.New("User does not exist")
	}
	d.ID = userID
	updateStatement := `update "classdata" set "cid"=$1, "name"=$2, "prereq"=$3 where "id"=$4`
	_, err = db.Exec(updateStatement, d.Cid, d.Name, d.PreReq, d.ID)
	if err != nil {
		return err
	}

	return nil
}
