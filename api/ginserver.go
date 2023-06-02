package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/intothevoid/kramerbot/models"
	persist "github.com/intothevoid/kramerbot/persist"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/spf13/viper"
)

type GinServer struct {
	UserStoreDB persist.DatabaseIF
	OzbScraper  *scrapers.OzBargainScraper
	CCCScraper  *scrapers.CamCamCamScraper
	Config      *viper.Viper
}

type SignupRequest struct {
	ChatId   string `json:"chatId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PreferencesRequest struct {
	ChatId    string   `json:"chatId"`
	OzbGood   bool     `json:"ozbGood"`
	OzbSuper  bool     `json:"ozbSuper"`
	AmzDaily  bool     `json:"amzDaily"`
	AmzWeekly bool     `json:"amzWeekly"`
	Keywords  []string `json:"keywords"`
}

func (gs *GinServer) StartServer() {
	router := gin.Default()
	router.GET("/users", gs.getUsers)
	router.GET("/users/:id", gs.getUserById)
	router.GET("/deals", gs.getDeals)
	router.POST("/signup", gs.signup)
	router.POST("/preferences", gs.setPreferences)
	router.POST("/authenticate", gs.authenticate)

	port := ":" + gs.Config.GetString("ginserver.port")

	if port != "" {
		router.Run(port)
	} else {
		// Failed to read config, use default port
		router.Run(":8080")
	}
}

// In this function if user i.e. chatId exists,
// then update username_chosen and password,
// else throw error as the user does not have a valid
// chatId obtained from @kramerbot
func (gs *GinServer) signup(c *gin.Context) {
	// check if user exists and get chatId
	user, err := gs.userExists(c)
	if user == nil || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
		return
	}

	// get username from post request
	username := c.Params.ByName("username")
	password := c.Params.ByName("password")
	user.UsernameChosen = username
	user.Password = password

	// write to store
	gs.UserStoreDB.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"users": "user updated"})
}

func (gs *GinServer) userExists(c *gin.Context) (*models.UserData, error) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		return nil, err
	}

	chatId := c.Params.ByName("chatId")
	chatIdInt, err := strconv.ParseInt(chatId, 10, 64)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
		return nil, err
	}

	_, ok := users.Users[chatIdInt]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"users": "user not found"})
		return nil, fmt.Errorf("user not found")
	}

	user := users.Users[chatIdInt]
	return user, nil
}

func (gs *GinServer) setPreferences(c *gin.Context) {
	// check if user exists and get chatId
	user, err := gs.userExists(c)
	if user == nil || err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
		return
	}

	// get preferences from post request
	var preferences PreferencesRequest
	if err := c.BindJSON(&preferences); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"preferences": err.Error()})
		return
	}

	// update user
	user.OzbGood = preferences.OzbGood
	user.OzbSuper = preferences.OzbSuper
	user.AmzDaily = preferences.AmzDaily
	user.AmzWeekly = preferences.AmzWeekly
	user.Keywords = preferences.Keywords
	gs.UserStoreDB.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"preferences": "preferences updated"})
}

func (gs *GinServer) authenticate(c *gin.Context) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
		return
	}

	// get username and password from post request
	var auth AuthRequest
	if err := c.BindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"authenticated": err.Error()})
		return
	}

	username := auth.Username
	password := auth.Password

	for user := range users.Users {
		if users.Users[user].UsernameChosen != "" {
			if users.Users[user].UsernameChosen == username && users.Users[user].Password == password {
				c.JSON(http.StatusOK, gin.H{"authenticated": true})
				return
			}
		}
	}

	c.JSON(http.StatusBadRequest, gin.H{"authenticated": err.Error()})
}

func (gs *GinServer) getUsers(c *gin.Context) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"users": users.Users})
	}
}

func (gs *GinServer) getUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	iid, _ := strconv.ParseInt(id, 10, 64)

	// read store
	users, err := gs.UserStoreDB.ReadUserStore()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"users": err.Error()})
	} else {
		val, ok := users.Users[iid]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": val})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": "user id not found"})
		}
	}
}

func (gs *GinServer) getDeals(c *gin.Context) {
	// get deals

	var deals = map[string]interface{}{}
	deals["OZB"] = gs.OzbScraper.GetData()
	deals["AMZ"] = gs.CCCScraper.GetData()

	if len(deals) > 0 {
		c.JSON(http.StatusOK, gin.H{"deals": deals})
	} else {
		c.JSON(http.StatusOK, gin.H{"deals": "no deals found"})
	}

}
