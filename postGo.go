package main

import (
	"fmt"
	"log"
	"main.go/post05"
	"net/http"
	"os"
)

//This program adds a random user the database defined the post05 postgres interface through a network connection

func main() {
	//add args for new course info including Cid, Name, PreReq
	args := os.Args
	fmt.Printf("num args: %d \n", len(args))
	var t post05.CourseData
	if len(args) > 1 {

		t = post05.CourseData{ //defines the new user based on the post05 struct
			Cid:    args[1],
			Name:   args[2],
			PreReq: args[3]}
	} else {

		t = post05.CourseData{ //defines the new user based on the post05 struct
			Cid:    "432",
			Name:   "Foundations Of Data Engineering",
			PreReq: "MSDS 410"}
	}

	//defining the network interface to connect to the database based on post0t.go
	post05.Hostname = "postgres"
	post05.Port = 5432
	post05.Username = "postgres"
	post05.Password = "happy"
	post05.Database = "msds"

	//list the users in the database
	data, err := post05.ListUsers()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, v := range data {
		fmt.Println(v)
	}

	id := post05.AddUser(t) //adds the defined user based on the post05 function
	if id == -1 {
		fmt.Println("There was an error adding user", t.Cid)
	}
	/*
		//demonstrating delete functionality
		err = DeleteUser(id)
		if err != nil {
			fmt.Println(err)
		}

			// Trying to delete it again!
			err = DeleteUser(id)
			if err != nil {
				fmt.Println(err)
			}

			//adds the user back in
			id = AddUser(t)
			if id == -1 {
				fmt.Println("There was an error adding user", t.Cid)
			}
	*/
	//changing prereq to demonstrate updating an entry
	t = post05.CourseData{
		Cid:    "432",
		Name:   "Foundations Of Data Engineering",
		PreReq: "MSDS 420"}

	err = post05.UpdateUser(t) //updating entry
	if err != nil {
		fmt.Println(err)
	}

	// now print users to http requests
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		users, err := post05.ListUsers()
		if err != nil {
			http.Error(w, "Error listing users", http.StatusInternalServerError)
			return
		}
		for _, user := range users {
			fmt.Fprintf(w, "Class: %d, %s, %s, %s \n", user.ID, user.Cid, user.Name, user.PreReq)
		}
	})

	log.Printf("Listening on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
