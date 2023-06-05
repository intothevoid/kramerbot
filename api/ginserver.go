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
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"POST", "GET", "OPTIONS"},
	// 	AllowCredentials: true,
	// 	AllowHeaders:     []string{"Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization", "accept", "Origin", "Cache-Control", "X-Requested-With", "access-control-allow-origin", "access-control-allow-headers"},
	// }))

	router.OPTIONS("/authenticate", func(c *gin.Context) {
		// Respond with 200 OK and the necessary CORS headers
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, Origin, Cache-Control, X-Requested-With")
		c.Writer.WriteHeader(http.StatusOK)
	})

	router.POST("/signup", gs.signup)
	router.POST("/preferences", gs.setPreferences)
	router.POST("/authenticate", gs.authenticate)
	router.GET("/users", gs.getUsers)
	router.GET("/users/:id", gs.getUserById)
	router.GET("/deals/:count", gs.getDeals)
	router.GET("/", gs.rootHello)

	port := ":" + gs.Config.GetString("ginserver.port")

	if port != "" {
		router.Run(port)
	} else {
		// Failed to read config, use default port
		router.Run(":3179")
	}
}

func (gs *GinServer) rootHello(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"result": "Hello from @kramerbot!"})
}

// In this function if user i.e. chatId exists,
// then update username_chosen and password,
// else throw error as the user does not have a valid
// chatId obtained from @kramerbot
func (gs *GinServer) signup(c *gin.Context) {
	// obtain params from POST request
	var signup SignupRequest
	if err := c.BindJSON(&signup); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
		return
	}

	// get user based on chatId
	currentUser, err := gs.userExists(signup.ChatId)
	if currentUser == nil || err != nil {
		errMsg := "user not found"
		if err != nil {
			errMsg = err.Error()
		}
		c.JSON(http.StatusInternalServerError, gin.H{"result": errMsg})
		return
	}

	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
		return
	}

	// check if username already exists
	for _, user := range users.Users {
		if user.UsernameChosen == signup.Username {
			c.JSON(http.StatusInternalServerError, gin.H{"result": "username already exists"})
			return
		}
	}

	// if we reach here, username is unique and valid
	// update user credentials
	currentUser.UsernameChosen = signup.Username
	currentUser.Password = signup.Password
	gs.UserStoreDB.UpdateUser(currentUser)
	c.JSON(http.StatusOK, gin.H{"result": "user registered"})
}

func (gs *GinServer) userExists(chatId string) (*models.UserData, error) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		return nil, err
	}

	chatIdInt, err := strconv.ParseInt(chatId, 10, 64)

	if err != nil {
		return nil, err
	}

	_, ok := users.Users[chatIdInt]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	user := users.Users[chatIdInt]
	return user, nil
}

func (gs *GinServer) setPreferences(c *gin.Context) {
	// get preferences from post request
	var preferences PreferencesRequest
	if err := c.BindJSON(&preferences); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
		return
	}

	// check if user exists and get chatId
	user, err := gs.userExists(preferences.ChatId)
	if user == nil || err != nil {
		errMsg := "user does not exist"
		if err != nil {
			errMsg = err.Error()
		}
		c.JSON(http.StatusInternalServerError, gin.H{"result": errMsg})
		return
	}

	// update user
	user.OzbGood = preferences.OzbGood
	user.OzbSuper = preferences.OzbSuper
	user.AmzDaily = preferences.AmzDaily
	user.AmzWeekly = preferences.AmzWeekly
	user.Keywords = preferences.Keywords
	gs.UserStoreDB.UpdateUser(user)
	c.JSON(http.StatusOK, gin.H{"result": "preferences updated"})
}

func (gs *GinServer) authenticate(c *gin.Context) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
		return
	}

	// get username and password from post request
	var auth AuthRequest
	if err := c.BindJSON(&auth); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"result": err.Error()})
		return
	}

	username := auth.Username
	password := auth.Password

	for user := range users.Users {
		if users.Users[user].UsernameChosen != "" {
			if users.Users[user].UsernameChosen == username && users.Users[user].Password == password {
				// If we find user and password is correct, we return chatId i.e user = id here
				c.JSON(http.StatusOK, gin.H{"result": user})
				return
			}
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"result": "user not found"})
}

func (gs *GinServer) getUsers(c *gin.Context) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": users.Users})
}

func (gs *GinServer) getUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	iid, _ := strconv.ParseInt(id, 10, 64)

	// read store
	users, err := gs.UserStoreDB.ReadUserStore()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"result": err.Error()})
	} else {
		val, ok := users.Users[iid]
		if ok {
			c.JSON(http.StatusOK, gin.H{"result": val})
		} else {
			c.JSON(http.StatusNotFound, gin.H{"result": "user not found"})
		}
	}
}

func (gs *GinServer) getDeals(c *gin.Context) {
	// get no. of deals to send back
	sCount := c.Params.ByName("count")
	i64Count, _ := strconv.ParseInt(sCount, 10, 64)
	iCount := int(i64Count)

	// get deals
	var deals = map[string]interface{}{}
	ozbDeals := gs.OzbScraper.GetData()
	amzDeals := gs.CCCScraper.GetData()

	if iCount < len(ozbDeals) && iCount < len(amzDeals) {
		deals["OZB"] = gs.OzbScraper.GetData()[:iCount]
		deals["AMZ"] = gs.CCCScraper.GetData()[:iCount]
	} else {
		deals["OZB"] = gs.OzbScraper.GetData()
		deals["AMZ"] = gs.CCCScraper.GetData()
	}

	if len(deals) > 0 {
		c.JSON(http.StatusOK, gin.H{"result": deals})
	} else {
		c.JSON(http.StatusNotFound, gin.H{"result": "no deals found"})
	}

}
