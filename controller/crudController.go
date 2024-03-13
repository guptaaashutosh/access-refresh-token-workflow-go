package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	// "log"
	"net/http"
	"strconv"

	"learn/httpserver/constants"
	"learn/httpserver/model"
	"learn/httpserver/repo"
	"learn/httpserver/setup"
	validation "learn/httpserver/validations"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func GetAllDataFromRedis(c *gin.Context) ([]model.User, error) {

	var allData []model.User

	ctxRedisClient, ctxExist := c.Get("redis-client")
	if !ctxExist {
		return allData, errors.New("redis client not available")
	}

	redisClient, ok := ctxRedisClient.(*redis.Client)
	if !ok {
		panic("not a valid redis client")
	}

	getMapData, err := redisClient.HGetAll(c, "all-data").Result()
	if err != nil {
		return allData, err
	}

	for _, v := range getMapData {
		var mapData model.User
		err := json.Unmarshal([]byte(v), &mapData)
		if err != nil {
			panic(err)
		}
		allData = append(allData, mapData)
	}
	return allData, nil

}

// // get api
// func Get(c *gin.Context) {
// 	var getData []model.User
// 	var err error

// 	allData, getErr := GetAllDataFromRedis(c)

// 	if getErr == nil {
// 		DB := setup.ConnectDB()
// 		//repositories initialization
// 		repos := repo.UserRepo(DB)
// 		getData, err = repos.GetData()
// 		if err != nil {
// 			panic(err)
// 		}
// 		c.JSON(http.StatusOK, gin.H{
// 			"getData-database": getData,
// 		})
// 		return
// 	}

//		c.JSON(http.StatusOK, gin.H{
//			"getMapData-redis": allData,
//		})
//	}
//
// get api
func Get(c *gin.Context) {
	var getData []model.GetUser
	var err error

	allData, getErr := GetAllDataFromRedis(c)

	if getErr == nil {
		DB := setup.ConnectDB()
		//repositories initialization
		repos := repo.UserRepo(DB)
		getData, err = repos.GetData()
		if err != nil {
			panic(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"getData-database": getData,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"getMapData-redis": allData,
	})
}



// ------------- put all data in redis -------------
func PutAllDataInRedis(redisClient *redis.Client) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}

	for k, v := range getData {
		fmt.Println(k, v)
		redisKey := strconv.Itoa(int(v.Id))
		// err := redisClient.HSet(context.Background(), "user:"+redisKey, "Id", v.Id, "Email", v.Email, "Name", v.Name, "Age", v.Age, "Address", v.Address).Err()
		err := redisClient.HSet(context.Background(), "user:"+redisKey, "Id", v.Id, "Email", v.Email, "Name", v.Name, "Age", v.Age, "Address", v.Address).Err()
		if err != nil {
			fmt.Printf("HSet Error: %s", err)
		}
	}
}



// create - applying transaction
func Create(c *gin.Context) {

	//vaidate data 
	var requestParam = map[string]string{
		constants.IdKey:  "number|len:2",
		constants.EmailKey: "string",
		constants.PasswordKey : "string|minLen:4|maxLen:8",
		constants.NameKey : "string",
		constants.AgeKey : "string",
		constants.AddressKey : "string",
		constants.ServiceIdKey : "slice|maxLen:20",
	}

	var userReq model.User

	err, invalidParameters, invalidParamsErrMsg := validation.ValidateParameters(c.Request, &userReq, &requestParam, nil, nil, nil, nil)

	if err != nil || invalidParameters !=nil || invalidParamsErrMsg != nil {
		fmt.Print("error : ",err)
		return
	}
	

	// DB := setup.ConnectDB()
	// repos := repo.UserRepo(DB)
	//check data
	// var user model.User
	// err := c.BindJSON(&user)
	// if err != nil {
	// 	panic(err)
	// }

	// //redis-client
	// ctxRedisClient, redisConnected := c.Get("redis-client")
	// redisClient := ctxRedisClient.(*redis.Client)

	// //apply transaction
	// tx, err := DB.Begin(c)
	// if err != nil {
	// 	// return err
	// 	log.Fatal("Error in transaction begin : ", err)
	// 	return
	// }

	// //insert into employee table
	// err = repos.CreateEmployee(user, tx)
	// if err != nil {
	// 	tx.Rollback(c)
	// 	c.JSON(500, gin.H{
	// 		"isCreated": false,
	// 	})
	// 	return
	// }
	// //insert into employee-service-pair
	// err = repos.CreateEmployeeServicePair(user.Id, user.Sid, tx)
	// if err != nil {
	// 	tx.Rollback(c)
	// 	c.JSON(500, gin.H{
	// 		"isCreated": false,
	// 	})
	// 	return
	// }

	// err = tx.Commit(c)

	// if err != nil {
	// 	log.Fatal(err)
	// 	c.JSON(500, gin.H{
	// 		"isCreated": false,
	// 		"message ":  "transaction failed",
	// 	})
	// 	return
	// }

	// fmt.Println("-- transaction committed --")

	// //if redis connected then insert only
	// if redisConnected {
	// 	redisKey := strconv.Itoa(int(user.Id))
	// 	createErr := redisClient.HSet(c, "user:"+redisKey, "Id", user.Id, "Email", user.Email, "Name", user.Name, "Age", user.Age, "Address", user.Address).Err()
	// 	if createErr != nil {
	// 		fmt.Printf("HSet create Error: %s", err)
	// 	}
	// }
	c.JSON(http.StatusOK, gin.H{
		"isCreated": true,
		"vaid-data": userReq,

	})
}

// ------------------------AssignNewServiceToUser -------------------------
func AssignNewServiceToUser(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.ServiceRepo(DB)

	//check data
	var NewService model.Service
	err := c.BindJSON(&NewService)
	if err != nil {
		panic(err)
	}

	//insert into database only
	isCreated, creationError := repos.CreateNewService(NewService)
	if creationError != nil {
		panic(creationError)
	}

	c.JSON(http.StatusOK, gin.H{
		"isCreatedService": isCreated,
		"created-service":  "successfully created service",
	})
}

// ------------------------------------------------------------------------

// --------------------- Delete with redis --------------------
func Delete(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)
	var id = c.Param("id")

	//delete from database
	isCreated, deletionError := repos.DeleteData(id)
	if deletionError != nil {
		panic(deletionError)
	}

	//redis-client
	ctxRedisClient, redisConnected := c.Get("redis-client")
	redisClient := ctxRedisClient.(*redis.Client)

	if redisConnected {
		// delErr := redisClient.HDel(c, "all-data", "user:"+id).Err()
		delErr := redisClient.Del(c, "user:"+id).Err()
		if delErr != nil {
			fmt.Printf("HSet delete Error: %s", delErr)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isDeleted":                 isCreated,
		"delete-response-database ": "deleted from database",
	})
}

// Update
func Update(c *gin.Context) {
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)

	id := c.Param("id")

	var user model.User
	err := c.BindJSON(&user)
	if err != nil {
		panic(err)
	}

	//update in database
	isUpdated, updationError := repos.UpdateData(user, id)
	if updationError != nil {
		panic(updationError)
	}

	//redis-client
	ctxRedisClient, redisConnected := c.Get("redis-client")
	redisClient := ctxRedisClient.(*redis.Client)

	//if redis connected then update only
	if redisConnected {
		// userValue, _ := json.Marshal(user)
		// createErr := redisClient.HMSet(c, "all-data", key, userValue).Err()
		createErr := redisClient.HSet(c, "user:"+id, "Id", id, "Email", user.Email, "Name", user.Name, "Age", user.Age, "Address", user.Address, "Sid", user.Sid).Err()
		if createErr != nil {
			fmt.Printf("HSet create Error: %s", err)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"isUpdated": isUpdated,
	})

}

// Login
func Login(c *gin.Context) {
	session := sessions.Default(c)
	//check if user has already loggedIn
	if session.Get("isAuthenticated") == true {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "you are already loggedIn",
		})
		return
	}

	var loginUserData model.Login
	err := c.BindJSON(&loginUserData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": err.Error(),
		})
		return
	}
	//database operation
	DB := setup.ConnectDB()
	repos := repo.UserRepo(DB)

	//if data is not in redis hit database
	loggedInStatus, loggedInError := repos.CheckUserExist(loginUserData)
	if loggedInError != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"Error": loggedInError.Error(),
		})
		return
	}
	//save authentication in session
	session.Set("isAuthenticated", true)
	session.Set("loggedInStatus", loggedInStatus)
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success":                true,
		"message-db-login-check": "successfully loggedIn",
	})

}

// Logout
func Logout(c *gin.Context) {
	//save authentication in session

	session := sessions.Default(c)

	if session.Get("isAuthenticated") == false {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request to logout - already logout",
		})
		return
	}

	session.Set("isAuthenticated", false)
	session.Set("loggedInToken", "")
	session.Save()
	// session.Clear()
	// session.Save()

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "successfully logout",
	})
}

// AuthData
func AuthData(c *gin.Context) {
	DB := setup.ConnectDB()
	//repositories initialization
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{
		"message":   "successfully authenticated and authorized",
		"Auth Data": getData,
	})
}

// SessionTest
func SessionTest(c *gin.Context) {
	DB := setup.ConnectDB()
	//repositories initialization
	repos := repo.UserRepo(DB)
	getData, err := repos.GetData()
	if err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"session-test-message": "successfully test session",
		"AllData":              getData,
	})
}




// func Create(c *gin.Context) {
// 	DB := setup.ConnectDB()
// 	repos := repo.UserRepo(DB)
// 	//check data
// 	var user model.User
// 	err := c.BindJSON(&user)
// 	if err != nil {
// 		panic(err)
// 	}
// 	//redis-client
// 	ctxRedisClient, redisConnected := c.Get("redis-client")
// 	redisClient := ctxRedisClient.(*redis.Client)

// 	//apply transaction
// 	//insert into database only in employee table
// 	isCreated, creationError := repos.CreateData(user)
// 	if creationError != nil {
// 		panic(creationError)
// 	}

// 	//if redis connected then insert only
// 	if redisConnected {
// 		redisKey := strconv.Itoa(int(user.Id))
// 		// createErr := redisClient.HSet(c, "user:"+redisKey, "Id", user.Id, "Email", user.Email, "Name", user.Name, "Age", user.Age, "Address", user.Address,"Sid",user.Sid).Err()
// 		createErr := redisClient.HSet(c, "user:"+redisKey, "Id", user.Id, "Email", user.Email, "Name", user.Name, "Age", user.Age, "Address", user.Address).Err()
// 		if createErr != nil {
// 			fmt.Printf("HSet create Error: %s", err)
// 		}
// 	}
// 	c.JSON(http.StatusOK, gin.H{
// 		"isCreated": isCreated,
// 	})
// }