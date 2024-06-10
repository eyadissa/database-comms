package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mactsouk/post05"
)

//This program adds a random user the database defined the post05 postgres interface through a network connection

// initializing variables to make a random user ID
var MIN = 0
var MAX = 26

// defining the type of randomizer we will use (integers between MIN and MAX)
func random(min, max int) int {
	return rand.Intn(max-min) + min
}

// create a string of letters from the randomizer
func getString(length int64) string {
	startChar := "A"
	temp := ""
	var i int64 = 1
	for {
		myRand := random(MIN, MAX)
		newChar := string(startChar[0] + byte(myRand)) //bit addition of rand output creates a new random letter
		temp = temp + newChar                          //append the new letter to the string
		if i == length {                               //leave the function if we've reached our desired rand string length
			break
		}
		i++
	}
	return temp
}

func main() {
	//defining the network interface to connect to the database based on post0t.go
	post05.Hostname = "localhost"
	post05.Port = 5433
	post05.Username = "postgres"
	post05.Password = "root"
	post05.Database = "go"

	//list the users in the database
	data, err := post05.ListUsers()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(v)
	}

	SEED := time.Now().Unix()
	rand.Seed(SEED)                 //sets a random state based on time
	random_username := getString(5) //generates a random username of length 5

	t := post05.Userdata{ //defines the new user based on the post05 struct
		Username:    random_username,
		Name:        "Mihalis",
		Surname:     "Tsoukalos",
		Description: "This is me!"}

	id := post05.AddUser(t) //adds the defined user based on the post05 function
	if id == -1 {
		fmt.Println("There was an error adding user", t.Username)
	}

	//demonstrating delete functionality
	err = post05.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
	}

	// Trying to delete it again!
	err = post05.DeleteUser(id)
	if err != nil {
		fmt.Println(err)
	}

	//adds the user back in
	id = post05.AddUser(t)
	if id == -1 {
		fmt.Println("There was an error adding user", t.Username)
	}

	//changing discriptoin to demonstrate updating an entry
	t = post05.Userdata{
		Username:    random_username,
		Name:        "Mihalis",
		Surname:     "Tsoukalos",
		Description: "This might not be me!"}

	err = post05.UpdateUser(t) //updating entry
	if err != nil {
		fmt.Println(err)
	}
}
