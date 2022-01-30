package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	handleRequests()
	http.ListenAndServe(":8087", nil)
}

func handleRequests() {
	http.HandleFunc("/hello", homePage)
	http.HandleFunc("/risk-analysis", riskAnalysis)
}

func riskAnalysis(w http.ResponseWriter, r *http.Request) {
	var t Transactions

	err := json.NewDecoder(r.Body).Decode(&t)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var risks []string
	users := make(map[int]*User)
	for _, eachT := range t.Tr {

		var u User
		u.ID = eachT.UserId
		if !u.userExist(&users) {
			users[u.ID] = &u
		}

		risk := ""
		dollars := (eachT.Amount / 100)
		u.Total = users[u.ID].Total + dollars
		var cardNumbers = u.cardNumbers(users, eachT.CardId)
		if dollars > 10000 || users[u.ID].Total > 20000 || cardNumbers > 2 {
			risk = "high"
			risks = append(risks, risk)
			continue
		}

		if dollars > 5000 || users[u.ID].Total > 10000 || cardNumbers > 1 {
			risk = "medium"
			risks = append(risks, risk)
			continue
		}

		risk = "low"
		risks = append(risks, risk)
	}

	resp := map[string]interface{}{
		"risk_analysis": risks,
	}

	res, err := json.MarshalIndent(resp, "", "")
	fmt.Fprintf(w, string(res))
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello\n")
}

type Transaction struct {
	Id     int     `json:"id"`
	UserId int     `json:"user_id"`
	Amount float64 `json:"amount_us_cents"`
	CardId int     `json:"card_id"`
}

type Transactions struct {
	Tr []Transaction `json:"transactions"`
}

type User struct {
	ID    int
	Total float64
	Cards []int
}

func (u *User) userExist(users *map[int]*User) bool {
	for _, us := range *users {
		if us.ID == u.ID {
			return true
		}
	}

	return false
}

func (u *User) cardNumbers(users map[int]*User, cardID int) int {

	if contains(users[u.ID].Cards, cardID) {
		return len(users[u.ID].Cards)
	} else {
		u := users[u.ID]
		u.Cards = append(u.Cards, cardID)
		users[u.ID] = u
		return len(users[u.ID].Cards)
	}
}

func contains(s []int, str int) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
