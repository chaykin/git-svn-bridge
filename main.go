package main

import (
	"git-svn-bridge/cmd"
	"git-svn-bridge/plock"
)

func main() {
	plock.Lock()
	cmd.Execute()

	/*	addUserOption := flag.String("add-user", "", "Add user information to the database")
		//changePassOption := flag.String("p", "", "Change usr's password in the database")
		//	resetAuthOption := flag.String("r", "", "Reset SVN auth cache with usr's credentials; option argument is usr's email; SVN URL required as non-option argument")
		syncOption := flag.Bool("sync", false, "Sync SVN and GIT repositories")
		flag.Parse()

		if *addUserOption != "" {
			addUser(*addUserOption)
		} else if *syncOption {
			flag.Args()
			fmt.Println("syncOption")
		}*/
}

/*func addUser(svnUserName string) {
	fmt.Printf("Adding/overwriting SVN usr: %s\n", svnUserName)
	pass := readPassFromInput()
	mail := readFieldFromInput("e-mail")
	//TODO: check mail valid
	gitUserName := readFieldFromInput("Git user name")

	user := usr.CreateUser(svnUserName, pass, gitUserName, mail)
	store.StoreUser(user)
}*/
