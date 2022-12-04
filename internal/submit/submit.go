package submit

import (
	"net/http"
)

func main(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	email := request.FormValue("email")
	if email == "" {
		http.Redirect(response, request, "/", http.StatusSeeOther)
		return
	}

	signedString, err := getJWT(email)
	if err != nil {
		http.Error(response, "couldn't get JWT", http.StatusInternalServerError)
		return
	}

	cookie := http.Cookie{
		Name:  "session",
		Value: signedString,
	}

	http.SetCookie(response, &cookie)
	http.Redirect(response, request, "/", http.StatusSeeOther)
}
