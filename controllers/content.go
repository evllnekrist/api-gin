package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	// "time"

	// "api-gin/forms"
	"api-gin/helpers"
	"api-gin/models"
	"github.com/gin-gonic/gin"
)

type ContentController struct{}

var helper_eve = new(helpers.EveHelper)
var model_redis = new(models.RedisUtil)

type Content struct{}

func DisplayInfo(c *gin.Context, err error, info string) {
	if err != nil {
		if info != "" {
			info = "oops.. sorry, error occurring"
		}
		panic(err.Error())
	}
	c.HTML(http.StatusOK, "action_info.html", gin.H{"docNecessity": info})
}

func GetTypeManual(prefix string) string {
	if prefix == "ex:string:" || prefix == "okz:content:" {
		return "string"
	} else if prefix == "ex:list:" || prefix == "okz:headline:" || prefix == "okz:breaking:" {
		return "list"
	} else if prefix == "ex:set:" {
		return "set"
	} else if prefix == "ex:zset:" {
		return "zset"
	} else if prefix == "ex:hash:" {
		return "hash"
	} else {
		return "none"
	}
}

func InputByType(c *gin.Context, err error, key_type string, key string, detail []byte) {
	if err != nil {
		DisplayInfo(c, err, "Oops.. cant input/update "+key)
	} else {
		var err2 error
		if key_type == "string" {
			err2 = model_redis.RedisSet(key, detail)
		} else if key_type == "list" { //bisa dgn Lpush juga
			err2 = model_redis.RedisRpush(key, detail)
		} else if key_type == "set" {
			err2 = model_redis.RedisSet(key, detail)
		} else if key_type == "zset" {
			DisplayInfo(c, err, "Sorry, input/update/delete/read type data ZSET not available yet")
		} else if key_type == "hash" {
			DisplayInfo(c, err, "Sorry, input/update/delete/read type data HASH not available yet")
		} else {
			DisplayInfo(c, err, "Sorry, if your input key not wrong then we just dont know what key type  \""+key_type+"\" is")
		}

		if key_type == "string" || key_type == "list" || key_type == "set" {
			DisplayInfo(c, err2, "Successfully MAKING CHANGE on -->  "+key+" ("+key_type+")")
		}
	}
}

func (ctrl ContentController) Detail(c *gin.Context, prefix string) {
	_type := c.Param("type")
	_id := c.Param("id")

	if _type == "json" {
		result, status := model_redis.RedisGet(prefix + _id)
		output, _ := helper_eve.GenerateJsonDesc(result, status)
		// c itu gabungan dr w http.ResponseWriter, r *http.Request. Jadi kalau mau ambil salah satu --> c.Writer / c.Request
		w := c.Writer                                      // r := c.Request
		w.Header().Set("Content-Type", "application/json") //tipe response
		w.Write(output)                                    //mendaftarkan data sebagai response
	} else {
		fmt.Println("type not available")
	}
}

func (ctrl ContentController) List(c *gin.Context, prefix string) {
	_type := c.Param("type")
	_id := c.Param("id")
	_limit := c.Param("limit")
	_start := c.Param("start")

	if _type == "json" {
		result, status := model_redis.RedisLrange(prefix+_id, _limit, _start) //key,limit,start
		lastIndex := len(result) - 1
		stack := "{"
		for i, each := range result {
			stack = stack + "\"" + strconv.Itoa(i) + "\":" + each
			if i != lastIndex {
				stack = stack + ","
			}
		}
		stack = stack + "}"
		output, _ := helper_eve.GenerateJsonDesc(stack, status)
		w := c.Writer                                      // r := c.Request
		w.Header().Set("Content-Type", "application/json") //tipe response
		w.Write(output)                                    //mendaftarkan data sebagai response
	} else {
		fmt.Println("type not available")
	}
}

/*_____________________________.  CRUD  start  ._____________________________*/
/*___________________________________________________________________________*/

func (ctrl ContentController) Create(c *gin.Context) { // available only via GUI
	r := c.Request
	if r.Method == "POST" {
		prefix := r.FormValue("prefix")
		detail := []byte(r.FormValue("detail"))
		id, err := model_redis.RedisIncr(prefix) //get the last key of group
		if err != nil {
			DisplayInfo(c, err, "Oops.. system cant create a key")
		} else {
			key := prefix + strconv.Itoa(id)
			key_type := GetTypeManual(prefix)
			InputByType(c, nil, key_type, key, detail) //input & update redis itu sama
		}
	} else {
		fmt.Println("only available in POST method")
	}
}

func (ctrl ContentController) Update(c *gin.Context) {
	r := c.Request
	if r.Method == "POST" {
		prefix := r.FormValue("prefix")
		detail := []byte(r.FormValue("detail"))
		id := r.FormValue("id")

		key := prefix + id
		key_type, err := model_redis.RedisType(key)
		InputByType(c, err, key_type, key, detail) //input & update redis itu sama
	} else {
		fmt.Println("only available in POST method")
	}
}

func (ctrl ContentController) Drop(c *gin.Context) {
	r := c.Request
	if r.Method == "POST" {
		prefix := r.FormValue("prefix")
		id := r.FormValue("id")
		key := prefix + id
		err := model_redis.RedisDelete(key)
		DisplayInfo(c, err, "Successfully DELETING -->  "+key)
	}
}

func (ctrl ContentController) Headline(c *gin.Context) { // EXAMPLE :: http://localhost:8080/okz/headline/json/0/4/0
	ctrl.List(c, "okz:headline:")
}

func (ctrl ContentController) Breaking(c *gin.Context) { // EXAMPLE :: http://localhost:8080/okz/breaking/json/16/7/0
	ctrl.List(c, "okz:breaking:")
}

func (ctrl ContentController) Read(c *gin.Context) { // EXAMPLE :: http://localhost:8080/okz/read/json/1555275
	ctrl.Detail(c, "okz:content:")
}

/*______________________________.  CRUD  end  .______________________________*/
/*___________________________________________________________________________*/
