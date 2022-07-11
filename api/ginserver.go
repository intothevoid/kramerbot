package api

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/intothevoid/kramerbot/persist"
	"github.com/intothevoid/kramerbot/scrapers"
	"github.com/spf13/viper"
)

type GinServer struct {
	UserStoreDB *persist.UserStoreDB
	OzbScraper  *scrapers.OzBargainScraper
	CCCScraper  *scrapers.CamCamCamScraper
	Config      *viper.Viper
}

func (gs *GinServer) StartServer() {
	router := gin.Default()
	router.GET("/users", gs.getUsers)
	router.GET("/users/:id", gs.getUserById)
	router.GET("/deals", gs.getDeals)

	port := ":" + gs.Config.GetString("ginserver.port")

	if port != "" {
		router.Run(port)
	} else {
		// Failed to read config, use default port
		router.Run(":8080")
	}
}

func (gs *GinServer) getUsers(c *gin.Context) {
	// read store
	users, err := gs.UserStoreDB.ReadUserStore()

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
