package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loginHandler(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	session := sessions.Default(c)

	if strings.Trim(user.Username, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username can't be empty"})
	}

	if strings.Trim(user.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password can't be empty"})
	}

	if !validCredentials(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid username/password combination"})
	}

	session.Set("user", user.Username)
	session.Set("authType", getAuthType(user))

	err := session.Save()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to generate session token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "authentication successful"})
}

func deleteUserHandler(c *gin.Context) {
	var userDel userDeletion

	if err := c.ShouldBind(&userDel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	if strings.Trim(userDel.Username, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Username can't be empty"})
		return
	}
	user := user{
		Username: userDel.Username,
	}
	if !userExists(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User doesn't exist"})
		return
	}

	if !deleteUser(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unable to delete user at this time"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Success!"})
}

func registrationHandler(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	if strings.Trim(user.Username, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username can't be empty"})
		return
	}

	if strings.Trim(user.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password can't be empty"})
		return
	}

	if strings.Trim(user.Code, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invite code is required to register"})
		return
	}

	user.AuthType = "user"

	if userExists(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "user already exists"})
		return
	}

	ic := inviteCode{
		Code:     user.Code,
		Username: user.Username,
	}

	if !validateInviteCode(ic) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "the invite code provided is invalid or has already been used"})
		return
	}

	fmt.Println(ic)

	if !assignInviteCode(ic) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to validate inviteCode"})
		return
	}

	err := increasePrize(user.Username, time.Now(), 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}

	if !addUser(user) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "unable to register user at this time, please try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User added, please login with the credentials provided"})
}

func logoutHandler(c *gin.Context) {
	session := sessions.Default(c)

	session.Delete("user")
	session.Delete("authType")

	err := session.Save()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to remove session token"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "successfully logged out"})
}

func measurementHistoryHandler(c *gin.Context) {
	measurements := getAllMeasurements()
	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements})
}

func listUsersHandler(c *gin.Context) {
	users := getAllUsers()
	c.JSON(http.StatusOK, gin.H{
		"users": users})
}

func listInviteCodesHandler(c *gin.Context) {
	invites := getAllInviteCodes()
	c.JSON(http.StatusOK, gin.H{
		"invite-codes": invites})
}

func myMeasurementsHandler(c *gin.Context) {
	session := sessions.Default(c)
	sessionusername := session.Get("user")
	username, ok := sessionusername.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to determine current user"})
		return
	}
	u := user{
		Username: username,
	}
	measurements := getUserMeasurements(u)
	c.JSON(http.StatusOK, gin.H{
		"measurements": measurements})
}

func newMeasurementHandler(c *gin.Context) {
	var updm measurement

	if err := c.ShouldBind(&updm); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}
	session := sessions.Default(c)
	sessionusername := session.Get("user")
	username, ok := sessionusername.(string)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to determine current user"})
		return
	}

	if updm.Waist == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Waist measurement can't be 0"})
		return
	}

	updm.Username = username
	updm.Date = time.Now()

	if checkIfPays(updm) {
		err := increasePrize(updm.Username, updm.Date, 5)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Could not increase the prize at the moment"})
			return
		}
	}

	if newMeasurement(updm) {
		c.JSON(http.StatusOK, gin.H{
			"message": "Record added"})
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "Could not create a new record, please try again later",
			"record": updm})
		return
	}

}

func prizeHistoryHandler(c *gin.Context) {
	prizes := getPrizeHistory()
	c.JSON(http.StatusOK, gin.H{
		"prize-history": prizes})
}

func standingsHandler(c *gin.Context) {
	standings := getRankings()
	c.JSON(http.StatusOK, gin.H{
		"standings": standings})
}

func newInviteCodeHandler(c *gin.Context) {
	ic, err := genInviteCode()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"invite-code": ic.Code})
}

func resetPasswordHandler(c *gin.Context) {
	var user user

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return
	}

	if strings.Trim(user.Username, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "username can't be empty"})
		return
	}

	if strings.Trim(user.Password, " ") == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "password can't be empty"})
		return
	}

	if !userExists(user) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "The user does not exist"})
		return
	}

	if !updatePassword(user) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to update the password at this time, please try again later"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "password successfully updated"})
}
