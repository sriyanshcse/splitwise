package main

import (
	//"database/sql"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	//_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	//"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"log"
	"net/http"
	//"os"
	"time"
	//"flag"
	"net/smtp"
	"strconv"
)

var Map map[int]int

var cur_id int

var mutex *sync.Mutex

type handlerError struct {
	Error   error
	Message string
	Code    int
}

type node struct {
	Id    int      `json:"id,omitempty"`
	Name  string   `json:"name,omitempty"`
	Email string   `json:"email,omitempty"`
	Tedge []*tEdge `json:"tedge,omitempty"`
	Redge []*rEdge `json:"fedge,omitempty"`
	Nedge []*nEdge `json:"nedge,omitempty'`
}

type rEdge struct {
	To int
}

type tEdge struct {
	Name  string    `json:"name,omitempty"` // bill name
	To    int       `json:"to,omitempty"`
	Value int       `json:"value,omitempty"`
	Time  time.Time `json:"time,omitempty"`
}

type nEdge struct {
	To    int `json:"to,omitempty"`
	Value int `json:"value,omitempty"`
}

type graph struct {
	Vertices []*node
}

type modifyGraph interface {
	AddVertex(node)
	//RemoveVertex(id int, ty string)
	AddUndirEdge(int, int)
	AddDirEdge(int, int, int, string)
	//RemoveEdge(from int, to int)
	UpdateValueEdge(int, int)
	AddTransaction(int, int, int, string)
	AddRelationship(int, int)
	UpdateNetAmount(int, int, int)
	GetNetAmount(int, int)
}

/*
type Env struct {
	DB *sql.DB
	Host string
	Port string
}*/

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func (g *graph) AddVertex(Node node) int {

	cur_id++
	Node.Id = cur_id + 1
	Map[cur_id] = cur_id
	g.Vertices = append(g.Vertices, &Node)
	return cur_id
}

func (g *graph) AddUndirEdge(from int, to int) {
	g.Vertices[from].Redge = append(g.Vertices[from].Redge, &rEdge{
		To: to + 1,
	})
}

func (g *graph) AddDirEdge(from int, to int, value int, name string) {

	g.Vertices[from].Tedge = append(g.Vertices[from].Tedge, &tEdge{
		Name:  name,
		To:    to + 1,
		Value: value,
		Time:  time.Now(),
	})

}

func (g *graph) UpdateValueEdge(from int, to int, value int) {

	for i, val := range g.Vertices[from].Nedge {
		if val.To == to+1 {
			g.Vertices[from].Nedge[i].Value += value
			return
		}
	}
	g.Vertices[from].Nedge = append(g.Vertices[from].Nedge, &nEdge{
		To:    to + 1,
		Value: value,
	})
}

func (g *graph) AddTransaction(from int, to int, value int, name string) {
	g.AddDirEdge(from, to, value, name)
	g.AddDirEdge(to, from, -value, name)
	g.UpdateNetAmount(from, to, value)
}

func (g *graph) AddRelationship(from int, to int) {
	g.AddUndirEdge(from, to)
	g.AddUndirEdge(to, from)
}

func (g *graph) UpdateNetAmount(from int, to int, value int) {
	g.UpdateValueEdge(from, to, value)
	g.UpdateValueEdge(to, from, -value)
}

func (g *graph) GetNetAmount(from int, to int) int {
	for _, val := range g.Vertices[from].Nedge {
		if val.To == to+1 {
			return val.Value
		}
	}
	return 0
}

type handler struct {
	*graph
	H func(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError)
}

func (fn handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	response, err := fn.H(fn.graph, w, r)
	log.Println("Hello")
	if err != nil {
		log.Printf("ERROR: %v\n", err.Error)
		http.Error(w, fmt.Sprintf(`{"error":"%s"}`), err.Code)
		return
	}

	if response == nil {
		log.Printf("ERROR: response from method is nil\n")
		http.Error(w, "Internal Server Error. Check Logs.", http.StatusInternalServerError)
		return
	}

	bytes, e := json.MarshalIndent(response, "", " ")
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}
	log.Printf("hello")
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func userExist(eid string, g *graph) (bool, int, string) {
	exist := false
	id := 0
	var name string
	log.Printf("%s", eid)
	for _, user := range g.Vertices {
		if user.Email == eid {
			exist = true
			id = user.Id
			name = user.Name
			break
		}
	}
	return exist, id, name
}

func GetFriends(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	id--
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}
	if _, ok := Map[id]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}

	friends := make([]interface{}, len(g.Vertices[id].Redge))
	log.Printf("%d", len(g.Vertices[id].Redge))
	cnt := 0
	for _, frds := range g.Vertices[id].Redge {
		var name string
		if name = g.Vertices[frds.To-1].Name; name == "" {
			name = g.Vertices[frds.To-1].Email
		}
		temp := struct {
			Name  string
			Id    int
			Value int
		}{
			Name:  name,
			Id:    g.Vertices[frds.To-1].Id,
			Value: g.GetNetAmount(id, frds.To-1),
		}
		log.Printf("%d", g.GetNetAmount(id, frds.To-1))
		friends[cnt] = temp
		cnt++
	}
	log.Printf("%d\n", cnt)
	return friends, nil
}

func GetGroups(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	/*	param := mux.Vars(r)["id"]
		id, e := strconv.Atoi(param)
		if e != nil {
			return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
		}
		if _, ok := Map[id]; !ok {
			return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
		}

		var groups []int

		for _, grps := range g.Vertices[id].Redge {
			if g.Vertices[grps.To].Type == "g" {
				groups = append(groups, grps.To)
			}
		}
		// DB calls

	*/return make(map[string]string), nil
}

func AddUserBill(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	id--
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}

	if _, ok := Map[id]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}

	data, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return nil, &handlerError{e, "Bill Not Found ", http.StatusNotFound}
	}
	var payload tEdge
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return nil, &handlerError{e, "Could not Parse JSON", http.StatusBadRequest}
	}
	g.AddTransaction(id, payload.To-1, payload.Value/2, payload.Name)
	return make(map[string]string), nil
}

func SettlePayment(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	param := mux.Vars(r)["id1"]
	id1, e := strconv.Atoi(param)
	id1--
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}
	param = mux.Vars(r)["id2"]
	id2, e := strconv.Atoi(param)
	id2--
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}

	if _, ok := Map[id1]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}
	if _, ok := Map[id2]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}

	g.AddTransaction(id1, id2, -g.GetNetAmount(id1, id2), "Settled Payment")
	return make(map[string]string), nil
}

func GetUserTransactionHistory(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	param := mux.Vars(r)["id1"]
	id1, e := strconv.Atoi(param)
	id1--
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}

	param = mux.Vars(r)["id2"]
	id2, e := strconv.Atoi(param)
	id2--
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}

	if _, ok := Map[id1]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}
	if _, ok := Map[id2]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}
	txns := make([]interface{}, len(g.Vertices[id1].Tedge))
	cnt := 0
	for _, txn := range g.Vertices[id1].Tedge {
		if txn.To == id2+1 {
			txns[cnt] = txn
			cnt++
		}
	}
	txns = txns[0:cnt]
	return txns, nil
}

func AddUser(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	data, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return nil, &handlerError{e, "Details Not Found ", http.StatusNotFound}
	}
	var payload node
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return nil, &handlerError{e, "Could not Parse JSON", http.StatusBadRequest}
	}
	added := false
	id := 0
	for _, user := range g.Vertices {
		if user.Email == payload.Email {
			user.Name = payload.Name
			id = user.Id
			added = true
			break
		}
	}
	if !added {
		id = g.AddVertex(payload) + 1
	}
	response := struct {
		Id int
	}{
		Id: id,
	}
	return response, nil
}

func AddFriend(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {
	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	id--
	//log.Printf("%d\n", id)
	if e != nil {
		return nil, &handlerError{e, "Id should be an integer", http.StatusBadRequest}
	}

	if _, ok := Map[id]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}
	data, e := ioutil.ReadAll(r.Body)

	if e != nil {
		return nil, &handlerError{e, "Details Not Found ", http.StatusNotFound}
	}
	var payload node
	e = json.Unmarshal(data, &payload)
	if e != nil {
		return nil, &handlerError{e, "Could not Parse JSON", http.StatusBadRequest}
	}

	for _, frd := range g.Vertices[id].Redge {
		if g.Vertices[frd.To-1].Email == payload.Email {
			return struct{ Added bool }{Added: false}, nil
		}
	}
	added := false
	//log.Printf("%s", payload.Email)
	var id2 int
	for _, user := range g.Vertices {
		log.Printf("%s -- %s", user.Email, payload.Email)
		if user.Email == payload.Email {
			added = true
			id2 = user.Id - 1
			break
		}
	}
	if added == false {
		id2 = g.AddVertex(payload)
	}
	log.Printf("%d", id2)
	g.AddRelationship(id, id2)
	return struct{ Added bool }{Added: true}, nil
}

func GetAllTransactions(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	param := mux.Vars(r)["id"]
	id, e := strconv.Atoi(param)
	id--
	if e != nil {
		return nil, &handlerError{e, "Details Not Found ", http.StatusNotFound}
	}
	if _, ok := Map[id]; !ok {
		return nil, &handlerError{e, "Id not found", http.StatusBadRequest}
	}
	txns := make([]interface{}, len(g.Vertices[id].Tedge))
	cnt := 0
	for _, txn := range g.Vertices[id].Tedge {
		txns[cnt] = txn
		cnt++
	}
	return txns[0:cnt], nil
}

func CheckUserExists(g *graph, w http.ResponseWriter, r *http.Request) (interface{}, *handlerError) {

	eid := mux.Vars(r)["eid"]
	exist, id, name := userExist(eid, g)
	response := struct {
		Exist bool
		Id    int
		Name  string
	}{
		Exist: exist,
		Id:    id,
		Name:  name,
	}
	return response, nil
}

func notFound(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "front/index.html")
}

func SendOTP(email string, otp string) {
	from := "deterministic007@gmail.com"
	pass := "contextfree"
	to := email

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: OTP for Splitwise Login \n\n" +
		"Hi,\nYour OTP for Splitwise Login is " + otp + ".This OTP is valid for 2 minutes.\n\n" +
		"Best Reagrds,\nTeam Splitwise."

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}

	log.Print("sent, visit gmail.com")
}

func GenerateOTP(w http.ResponseWriter, r *http.Request) {
	data, e := ioutil.ReadAll(r.Body)
	if e != nil {
		log.Printf("ERROR: body from client is nil\n")
		http.Error(w, "Internal Server Error. Check Logs.", http.StatusInternalServerError)
		return
	}
	payload := struct {
		Email string
	}{}
	e = json.Unmarshal(data, &payload)
	if e != nil {
		log.Printf("ERROR: Could not parse JSON\n")
		http.Error(w, "Internal Server Error. Check Logs.", http.StatusInternalServerError)
		return
	}
	text := payload.Email + time.Now().String()
	encrypted := GetMD5Hash(text)
	response := struct {
		OTP string
	}{
		OTP: encrypted,
	}
	bytes, e := json.MarshalIndent(response, "", " ")
	if e != nil {
		http.Error(w, "Error marshalling JSON", http.StatusInternalServerError)
		return
	}
	log.Printf(payload.Email)
	SendOTP(payload.Email, encrypted)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bytes)
}

func main() {

	//db, err := sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)")
	Map = make(map[int]int, 10000)
	cur_id = -1
	g := &graph{}
	/*env := &Env{
		DB: db,
		Port : os.Getenv("PORT"),
		Host : os.Getenv("HOST"),
	}*/
	mutex = &sync.Mutex{}

	router := mux.NewRouter()
	dir := http.Dir("./front")
	fileServer := http.FileServer(dir)
	router.Handle("/api/getFriends/{id}", handler{g, GetFriends}).Methods("GET")
	router.Handle("/api/getGroups/{id}", handler{g, GetGroups}).Methods("GET")
	router.Handle("/api/user/addBill/{id}", handler{g, AddUserBill}).Methods("POST")
	router.Handle("/api/addFriend/{id}", handler{g, AddFriend}).Methods("POST")
	//router.Handle("/group/addBill/{id}", handler{g, AddGroupBill}).Methods("POST")
	router.Handle("/api/settlePayment/{id1}/{id2}", handler{g, SettlePayment}).Methods("PUT")
	router.Handle("/api/user/getTransactionHistory/{id1}/{id2}", handler{g, GetUserTransactionHistory}).Methods("GET")
	//router.Handle("/group/getTransactionHistory/{id}", handler{g, GetGroupTransactionHistory}).Methods("GET")
	//router.Handle("/getUserDetails/{id}", handler{g, GetGroupTransaction}).Methods("GET")
	router.Handle("/api/addUser", handler{g, AddUser}).Methods("POST")
	router.Handle("/api/addFriend/{id}", handler{g, AddFriend}).Methods("POST")
	router.Handle("/api/getAllTransactionHistory/{id}", handler{g, GetAllTransactions}).Methods("GET")
	router.HandleFunc("/api/generateOTP", GenerateOTP).Methods("POST")
	router.Handle("/api/checkUserExists/{eid}", handler{g, CheckUserExists}).Methods("GET")
	//SendOTP("sriyansh.cse@gmail.com", "9879787987978")
	//router.Handle("/addGroup", handler{g, AddGroup}).Methods("POST")
	router.PathPrefix("/").Handler(fileServer)
	//router.NotFoundHandler = http.HandlerFunc(notFound)
	http.Handle("/", router)
	http.ListenAndServe(":3001", router)

}
