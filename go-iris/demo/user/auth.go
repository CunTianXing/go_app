package user

import (
	"errors"
	"strconv"
	"strings"

	"github.com/kataras/iris"
)

const (
	PathLogin  = "/user/login"
	PathLogout = "/user/logout"
	sessionIDKey = "UserID"
)

type AuthController struct {
	iris.SessionController

	Source *DataSource
	User   Model `iris:"model"`
}

func (c *AuthController) BeginRequest(ctx iris.Context) {
	c.SessionController.BeginRequest(ctx)

	if userID := c.Session.Get(sessionIDKey); userID != nil {
		ctx.Values().Set(sessionIDKey, userID)
	}
}

func (c *AuthController) fireError(err error) {
	if err != nil {
		c.Ctx.Application().Logger().Debug(err.Error())

		c.Status = 400
		c.Data["Title"] = "User Error"
		c.Data["Message"] = strings.ToUpper(err.Error())
		c.Tmpl = "shared/error.html"
	}
}

func (c *AuthController) redirectTo(id int64) {
	if id > 0 {
		c.Path = "/user/" + strconv.Itoa(int(id))
	}
}

func (c *AuthController) createOrUpdate(firstname, username, password string) (user Model, err error) {
	username = strings.Trim(username, " ")
	if username == "" || password == "" || firstname == "" {
		return user, errors.New("empty firstname, username or/and password")
	}

	userToInsert := Model{
		Firstname: firstname,
		Username:  username,
		password:  password,
	} 

	newUser, err := c.Source.InsertOrUpdate(userToInsert)
	if err != nil {
		return user, err
	}

	return newUser, nil
}

func (c *AuthController) isLoggedIn() bool {
	return c.Values.Get(sessionIDKey) != nil
}

func (c *AuthController) verify(username, password string) (user Model, err error) {
	if username == "" || password == "" {
		return user, errors.New("please fill both username and password fields")
	}

	u, found := c.Source.GetByUsername(username)
	if !found {
		return user, errors.New("user with that username does not exist")
	}

	if ok, err := ValidatePassword(password, u.HashedPassword); err != nil || !ok {
		return user, errors.New("please try to login with valid credentials")
	}

	return u, nil
}

func (c *AuthController) logout() {
	if c.isLoggedIn() {
		c.Manager.DestroyByID(c.Session.ID())
		return
	}

	c.Path = PathLogin
}

func AllowUser(ctx iris.Context) {
	if ctx.Values().Get(sessionIDKey) != nil {
		ctx.Next()
		return
	}
	ctx.Redirect(PathLogin)
}