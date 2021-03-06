package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"myGoProjectNew/db"
	"myGoProjectNew/models"
	"myGoProjectNew/pkg/util"
	"net/http"
	"strconv"
)

type UserC struct {
	Mgo *mongo.Database
	//RedisCli *redis.Client
}

/*login user*/
func (m UserC) Login(g *gin.Context) {
	fmt.Println("login.........")
	rsp := new(Rsp)
	name := g.PostForm("username")
	pass := g.PostForm("password")

	//var gerr *gvalid.Error
	//gerr = gvalid.Check(g.PostForm("username"), "required", nil)
	//if gerr != nil {
	//	rsp.Msg = "faild"
	//	rsp.Code = 201
	//	rsp.Data = gerr.Maps()
	//	g.JSON(http.StatusOK, rsp)
	//	return
	//}
	//gerr = gvalid.Check(g.PostForm("password"), "required", nil)
	//if gerr != nil {
	//	rsp.Msg = "faild"
	//	rsp.Code = 201
	//	rsp.Data = gerr.Maps()
	//	g.JSON(http.StatusOK, rsp)
	//	return
	//}

	findfilter := bson.D{{"username", g.PostForm("username")}, {"password", g.PostForm("password")}}
	cur, err := m.Mgo.Collection(db.User).Find(context.Background(), findfilter)
	if err != nil {
		rsp.Msg = "faild"
		rsp.Code = 201
		rsp.Data = err.Error()
		g.JSON(http.StatusOK, rsp)
		return
	}
	for cur.Next(context.Background()) {
		elme := new(models.User)
		err := cur.Decode(elme)
		if err == nil {
			if elme.Username == name && elme.Password == pass {
				var info = new(LoginInfo)
				info.User = elme
				token, err := util.GenerateToken(g.PostForm("username"), g.PostForm("password"))
				if err == nil {
					info.Token = token
				}
				rsp.Msg = "success"
				rsp.Code = 200
				rsp.Data = info
				g.JSON(http.StatusOK, rsp)
				return
			}
		}
	}

	rsp.Msg = "user is null"
	rsp.Code = 201
	rsp.Data = err
	g.JSON(http.StatusOK, rsp)
}

/* query all user */
func (m UserC) Queryalluser(g *gin.Context) {
	fmt.Println("Queryalluser.........")
	rsp := new(Rsp)
	var users []models.User
	cur, err := m.Mgo.Collection(db.User).Find(context.Background(), bson.D{}, nil)
	if err == nil {
		for cur.Next(context.Background()) {
			elme := new(models.User)
			err := cur.Decode(elme)
			if err == nil {
				users = append(users, *elme)
			}
		}
	}
	rsp.Msg = "success"
	rsp.Code = 200
	rsp.Data = users
	g.JSON(http.StatusOK, rsp)
	return
}

/* get all user */
func (m UserC) Getalluser(g *gin.Context) {
	log.Println("Getalluser.........")
	rsp := new(Rsp)

	filter := bson.M{}

	//limit, err := strconv.Atoi(g.Query("limit"))
	page, err := strconv.Atoi(g.Query("page"))
	fmt.Println(page)

	//?????? ??????1 ??????-1  ----------------------------
	opts := new(options.FindOptions)
	sortMap := make(map[string]interface{})
	sortMap["gender"] = -1
	opts.Sort = sortMap
	//opts.Limit=int64(limit)
	//?????? ??????1 ??????-1  ----------------------------

	var users []models.User
	cur, err := m.Mgo.Collection(db.User).Find(context.Background(), filter, opts)
	if err == nil {
		for cur.Next(context.Background()) {
			elme := new(models.User)
			err := cur.Decode(elme)
			if err == nil {
				users = append(users, *elme)
			}
		}
	}

	rsp.Msg = "success"
	rsp.Code = 0
	rsp.Data = users
	g.JSON(http.StatusOK, rsp)
	return
}
