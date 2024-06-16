package importer

import (
	"log"
	"os"
	"text/template"

	"github.com/ccin2p3/go-freeipa/freeipa"
	"golang.org/x/exp/slices"
)

// Entrypoint to create import users from an existing FreeIPA instance, to a terraform state.
func ImportUsers(client *freeipa.Client) error {
	ignoreUsers := []string{"admin"}
	var userList []freeipa.User
	log.Println("Importing all users.")
	res, err := client.UserFind("", &freeipa.UserFindArgs{}, &freeipa.UserFindOptionalArgs{
		Sizelimit: freeipa.Int(0),
		All:       &valTrue,
	})
	if err != nil {
		return err
	}
	log.Println("Found", len(res.Result), "users")

	//mark ignored users as ignored in the import
	for _, user := range res.Result {
		if slices.Contains(ignoreUsers, user.UID) {
			log.Println("  User", user.UID, "will be ignored")

		} else {
			log.Println("  User", user.UID, "will be imported")
			userList = append(userList, user)
		}
	}
	err = importUser(&userList)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("End of user import.")
	return nil
}

// Creates terraform resources related to a list of Users
func importUser(users *[]freeipa.User) error {
	var tmplFile = "templates/users.tf.tmpl"

	tmpl, err := template.New("users.tf.tmpl").Funcs(tmplFuncMap).ParseFiles(tmplFile)
	if err != nil {
		panic(err)
	}
	f, err := os.OpenFile("./output/users.tf", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	err = tmpl.Execute(f, users)
	if err != nil {
		panic(err)
	}
	return nil
}
