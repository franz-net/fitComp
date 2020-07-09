package main

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/rs/xid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Username string `form:"username" json:"username" xml:"username" db:"username" binding:"required"`
	Password string `form:"password" json:"password" xml:"password" db:"password" bidning:"required"`
	AuthType string `form:"authtype" json:"authtype" xml:"authtype" db:"authtype"`
	Code     string `form:"code" json:"code" xml:"code" db:"code"`
}

type userDeletion struct {
	Username string `form:"username" json:"username" xml:"username" db:"username" binding:"required"`
}

type measurement struct {
	ID       int       `form:"id" json:"id" xml:"id" db:"id"`
	Username string    `form:"username" json:"username" xml:"username" db:"username"`
	Date     time.Time `form:"date" json:"date" xml:"date" db:"date"`
	Waist    int       `form:"waist" json:"waist" xml:"waist" db:"waist" binding:"required"`
	Weight   int       `form:"weight" json:"weight" xml:"weight" db:"weight" binding:"required"`
}

type prize struct {
	ID       int       `form:"id" json:"id" xml:"id" db:"id"`
	Prize    int       `form:"prize" json:"prize" xml:"prize" db:"prize" binding:"required"`
	Increase int       `form:"increase" json:"increase" xml:"increase" db:"increase" binding:"required"`
	Username string    `form:"username" json:"username" xml:"username" db:"username"`
	Date     time.Time `form:"date" json:"date" xml:"date" db:"date" binding:"required"`
}

type inviteCode struct {
	ID       int    `form:"id" json:"id" xml:"id" db:"id"`
	Code     string `form:"code" json:"code" xml:"code" db:"code"`
	Username string `form:"username" json:"username" xml:"username" db:"username"`
}

func userExists(u user) bool {
	db := dbConn()
	var storedUsername string
	result := db.QueryRow("select password from users where username=?", u.Username)
	err := result.Scan(&storedUsername)
	if err == sql.ErrNoRows {
		return false
	}
	defer db.Close()
	return true
}

func deleteUser(u user) bool {
	db := dbConn()
	if _, err := db.Query("delete from users where username=?", u.Username); err != nil {
		return false
	}

	if _, err := db.Query("delete from measurements where username=?", u.Username); err != nil {
		return false
	}

	if _, err := db.Query("delete from prize where username=?", u.Username); err != nil {
		return false
	}

	if _, err := db.Query("delete from invites where username=?", u.Username); err != nil {
		return false
	}

	return true
}

func getAuthType(u user) string {
	db := dbConn()
	result := db.QueryRow("select authtype from users where username=?", u.Username)

	authType := ""

	err := result.Scan(&authType)
	if err != nil {
		return ""
	}
	defer db.Close()
	return authType
}

func addUser(u user) bool {
	db := dbConn()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

	if _, err = db.Query("insert into users values (?, ?, ?)", u.Username, string(hashedPassword), u.AuthType); err != nil {
		return false
	}
	defer db.Close()
	return true
}

func newMeasurement(m measurement) bool {
	db := dbConn()

	if _, err := db.Query("insert into measurements values (null, ?, ?, ?, ?)", m.Username, m.Date.Format("2006-01-02"), m.Waist, m.Weight); err != nil {
		return false
	}
	defer db.Close()
	return true
}

func checkIfPays(m measurement) bool {
	db := dbConn()
	var previousM int
	res := db.QueryRow("select coalesce(max(waist), 0) as waist from measurements where username=?", m.Username)
	err := res.Scan(&previousM)
	if err != nil {
		return false
	}
	if previousM == 0 {
		return false
	} else if m.Waist > previousM {
		return true
	}
	return false
}

func increasePrize(username string, date time.Time, increment int) error {
	db := dbConn()
	var previousPrize int
	res := db.QueryRow("select coalesce(max(prize), 0) from prize")
	err := res.Scan(&previousPrize)
	if err != nil {
		return errors.New("error: could not retrieve previous prize")
	}
	var pI = prize{
		Prize:    (previousPrize + increment),
		Increase: increment,
		Username: username,
		Date:     date,
	}
	if _, err := db.Query("insert into prize values(null, ?, ?, ?, ?)", pI.Prize, pI.Increase, pI.Username, pI.Date.Format("2006-01-02")); err != nil {
		return errors.New("error: could not insert new prize into db")

	}
	defer db.Close()
	return nil
}

func validCredentials(u user) bool {

	db := dbConn()

	result := db.QueryRow("select password from users where username=?", u.Username)

	storedPass := &user{}

	err := result.Scan(&storedPass.Password)
	if err != nil {
		return false
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPass.Password), []byte(u.Password)); err != nil {
		return false
	}
	defer db.Close()
	return true
}

func updatePassword(u user) bool {

	db := dbConn()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)

	if _, err = db.Query("update users set password=? where username=?", string(hashedPassword), u.Username); err != nil {
		return false
	}
	defer db.Close()
	return true
}

func getAllMeasurements() []measurement {
	db := dbConn()
	result, err := db.Query("Select * FROM measurements ORDER BY date DESC")
	if err != nil {
		panic(err.Error())
	}
	m := measurement{}
	res := []measurement{}
	layout := "2006-01-02"
	for result.Next() {
		var id, waist, weight int
		var username, tempDate string
		var date time.Time
		err = result.Scan(&id, &username, &tempDate, &waist, &weight)
		if err != nil {
			panic(err.Error())
		}
		date, err := time.Parse(layout, tempDate)
		if err != nil {
			panic(err.Error())
		}
		m.ID = id
		m.Username = username
		m.Date = date
		m.Waist = waist
		m.Weight = weight

		res = append(res, m)
	}
	defer db.Close()
	return res
}

func getAllUsers() []user {
	db := dbConn()
	result, err := db.Query("Select username, authtype FROM users")
	if err != nil {
		panic(err.Error())
	}
	u := user{}
	res := []user{}
	for result.Next() {
		var username, authtype string
		err = result.Scan(&username, &authtype)
		if err != nil {
			panic(err.Error())
		}
		u.Username = username
		u.AuthType = authtype

		res = append(res, u)
	}
	defer db.Close()
	return res
}

func getAllInviteCodes() []inviteCode {
	db := dbConn()
	result, err := db.Query("Select * from invites")
	if err != nil {
		panic(err.Error())
	}
	ic := inviteCode{}
	res := []inviteCode{}
	for result.Next() {
		var id int
		var code, username string
		err = result.Scan(&id, &code, &username)
		if err != nil {
			panic(err.Error())
		}
		ic.ID = id
		ic.Code = code
		ic.Username = username

		res = append(res, ic)
	}
	defer db.Close()
	return res
}

func getUserMeasurements(u user) []measurement {
	db := dbConn()
	result, err := db.Query("Select * FROM measurements where username=? ORDER BY date DESC", u.Username)
	if err != nil {
		panic(err.Error())
	}
	m := measurement{}
	res := []measurement{}
	layout := "2006-01-02"
	for result.Next() {
		var id, waist, weight int
		var username, tempDate string
		var date time.Time
		err = result.Scan(&id, &username, &tempDate, &waist, &weight)
		if err != nil {
			panic(err.Error())
		}
		date, err := time.Parse(layout, tempDate)
		if err != nil {
			panic(err.Error())
		}
		m.ID = id
		m.Username = username
		m.Date = date
		m.Waist = waist
		m.Weight = weight

		res = append(res, m)
	}
	defer db.Close()
	return res
}

func getPrizeHistory() []prize {
	db := dbConn()
	result, err := db.Query("Select * FROM prize ORDER BY date DESC")
	if err != nil {
		panic(err.Error())
	}
	p := prize{}
	res := []prize{}
	layout := "2006-01-02"
	for result.Next() {
		var id, current, increase int
		var increasedBy, tempDate string
		var date time.Time
		err = result.Scan(&id, &current, &increase, &increasedBy, &tempDate)
		if err != nil {
			panic(err.Error())
		}
		date, err := time.Parse(layout, tempDate)
		if err != nil {
			panic(err.Error())
		}
		p.ID = id
		p.Prize = current
		p.Increase = increase
		p.Username = increasedBy
		p.Date = date

		res = append(res, p)
	}
	defer db.Close()
	return res
}

func getRankings() []measurement {
	db := dbConn()

	result, err := db.Query("select coalesce(min(waist), 0) as waist, weight, username from measurements where username!='admin' group by username order by waist ASC")
	if err != nil {
		panic(err.Error())
	}
	m := measurement{}
	res := []measurement{}
	for result.Next() {
		var waist int
		var username string

		err = result.Scan(&waist, &username)
		if err != nil {
			panic(err.Error())
		}
		m.Username = username
		m.Waist = waist

		res = append(res, m)
	}
	defer db.Close()
	return res
}

func genInviteCode() (inviteCode, error) {
	db := dbConn()
	var ic = inviteCode{
		Code: xid.New().String(),
	}
	if _, err := db.Query("insert into invites values(null, ?, null )", ic.Code); err != nil {
		return inviteCode{}, errors.New("error: could not generate a new invite code")
	}
	defer db.Close()
	return ic, nil
}

func assignInviteCode(ic inviteCode) bool {
	db := dbConn()

	if _, err := db.Query("update invites set username=? where code=?", ic.Username, ic.Code); err != nil {
		fmt.Println(err.Error())
		return false
	}
	defer db.Close()
	return true
}

func validateInviteCode(ic inviteCode) bool {
	db := dbConn()
	var storedCode string
	result := db.QueryRow("select code from invites where code=?", ic.Code)
	err := result.Scan(&storedCode)
	if err == sql.ErrNoRows {
		return false
	}
	var storedUser string
	result = db.QueryRow("select ifnull(username, '') from invites where code=?", ic.Code)
	err = result.Scan(&storedUser)
	if err == sql.ErrNoRows || storedUser == "" {
		return true
	}

	defer db.Close()
	return false
}
