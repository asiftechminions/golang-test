package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	handleRequests()
	http.ListenAndServe(":8086", nil)
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
	var users map[int]User
	for _, eachT := range t.Tr {

		var u User
		u.ID = eachT.UserId
		if !u.userExist(&users) {
			users[u.ID] = u
		}

		risk := ""
		if eachT.Amount > 10000 {
			risk = "high"
		} else if eachT.Amount > 5000 {
			risk = "medium"
		} else {
			risk = "low"
		}

		u.Total = users[u.ID].Total + eachT.Amount

		if users[u.ID].Total > 20000 {
			risk = "high"
		} else if users[u.ID].Total > 10000 {
			risk = "medium"
		} else {
			risk = "low"
		}

		var cardNumbers = u.cardNumbers(users, eachT.CardId)
		if cardNumbers > 2 {
			risk = "high"
		} else if cardNumbers > 1 {
			risk = "medium"
		} else {
			risk = "low"
		}

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

func (u *User) userExist(users *map[int]User) bool {
	for _, us := range *users {
		if us.ID == u.ID {
			return true
		}
	}

	return false
}

func (u *User) cardNumbers(users map[int]User, cardID int) int {

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
