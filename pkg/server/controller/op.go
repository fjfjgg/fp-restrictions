package controller

import (
	"log"
	"fmt"
	"errors"
	"net/http"
	"encoding/json"
	"reflect"
	"regexp"

	plugin "github.com/fatedier/frp/pkg/plugin/server"
	"github.com/gin-gonic/gin"
	
)

type MinimalRequest struct {
	Version string      `json:"version"`
	Op      string      `json:"op"`
}

type OpController struct {
	restrictions map[string]Restriction
}

func NewOpController(r map[string]Restriction) *OpController {
	return &OpController{
		restrictions: r,
	}
}

func (c *OpController) Register(engine *gin.Engine) {
	engine.POST("/handler", MakeGinHandlerFunc(c.HandleSelector))
}

func (c *OpController) HandleSelector(ctx *gin.Context) (interface{}, error) {
	var err error
	var res interface{}
	switch op := ctx.Query("op") ; op {
	case "Login":
		log.Print("Login\n")
		res, err = c.HandleLogin(ctx)
	case "NewProxy":
		log.Print("NewProxy\n")
		res, err = c.HandleNewProxy(ctx)
	default:
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  errors.New("Unsupported operation"),
		}
	}
	return res, err
}

//Login
func (c *OpController) HandleLogin(ctx *gin.Context) (interface{}, error) {
	var r plugin.Request
	var content plugin.LoginContent
	r.Content = &content
	if err := ctx.BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}
	data, _ := json.Marshal(r)
	log.Println(string(data))
	var res plugin.Response
	token := content.Metas["token"]
	if content.User == "" || token == "" {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
	} else if c.restrictions[content.User].Token == token {
		res.Unchange = true
	} else {
		res.Reject = true
		res.RejectReason = "invalid meta token"
	}
	return res, nil
}

//NewProxy
func (c *OpController) HandleNewProxy(ctx *gin.Context) (interface{}, error) {
	var r plugin.Request
	var content plugin.NewProxyContent
	r.Content = &content
	if err := ctx.BindJSON(&r); err != nil {
		return nil, &HTTPError{
			Code: http.StatusBadRequest,
			Err:  err,
		}
	}
	data, _ := json.Marshal(r)
	log.Println(string(data))
	var res plugin.Response
	token := content.User.Metas["token"]
	if content.User.User == "" || token == "" {
		res.Reject = true
		res.RejectReason = "user or meta token can not be empty"
	} else if c.restrictions[content.User.User].Token == token {
		//test other restrictions
		//t := reflect.Indirect(reflect.ValueOf(content.NewProxy)).Type()
		//for i := 0; i < t.NumField(); i++ {
	 	  // f := t.Field(i)
		  // fmt.Println(f.Name, f.Type)
		//}
		rv := reflect.Indirect(reflect.ValueOf(content.NewProxy))
		for k, v := range c.restrictions[content.User.User].Restriction { 
		    //log.Println("Testing key-value ["+k+"] =",v)
		    //Test if exist key
		    value := rv.FieldByName(k)
		    if value.IsValid() {
		    	//log.Println("Exist [", value,"]")
		    	vType := reflect.TypeOf(v)
			switch vType.Kind() {
			case reflect.String:
			    //log.Println(k, "is a string")
			    //sv :=  value.String()
			    sv :=  fmt.Sprintf("%v", value )
			    match, _ := regexp.MatchString(v.(string), sv)
			    if match == false {
		    	    	res.Reject = true
		    	    	res.RejectReason = "restriction applied"
		    	    	log.Println("Rejected expected:",v,"requested:",sv)
			    }
			case reflect.Array:
			    log.Println(v, "is an array (unsupported type) with element type",
			                vType.Elem())
			    log.Println("Request value =", value)
			default:
			    log.Println(k, "unsupported type", vType)
			    log.Println("Request value =", value)
			}
		    		
		    }
		    if res.Reject == true {
		    	break;
		    }
		}		
		res.Unchange = true
	} else {
		log.Println("esperado" +  c.restrictions[content.User.User].Token)
		res.Reject = true
		res.RejectReason = "invalid meta token"
	}
	return res, nil
}

