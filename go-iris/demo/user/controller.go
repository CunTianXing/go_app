package user

const (
	pathMyProfile = "/user/me"
	pathRegister  = "/user/register"
)

type Controller struct {
	AuthController
}

// GetRegister handles GET:/user/register.
func (c *Controller) GetRegister() {
	if c.isLoggedIn() {
		c.logout()
		return
	}

	c.Data["Title"] = "User Registration"
	c.Tmpl = pathRegister + ".html"
}

// PostRegister handles POST:/user/register.
func (c *Controller) PostRegister() {
	// we can either use the `c.Ctx.ReadForm` or read values one by one.
	var (
		firstname = c.Ctx.FormValue("firstname")
		username  = c.Ctx.FormValue("username")
		password  = c.Ctx.FormValue("password")
	)

	user, err := c.createOrUpdate(firstname, username, password)
	if err != nil {
		c.fireError(err)
		return
	}

	c.Session.Set(sessionIDKey, user.ID)

	c.Status = 303 // "See Other" RFC 7231

	
	c.Path = pathMyProfile
}

// GetLogin handles GET:/user/login.
func (c *Controller) GetLogin() {
	if c.isLoggedIn() {
		c.logout()
		return
	}
	c.Data["Title"] = "User Login"
	c.Tmpl = PathLogin + ".html"
}

// PostLogin handles POST:/user/login.
func (c *Controller) PostLogin() {
	var (
		username = c.Ctx.FormValue("username")
		password = c.Ctx.FormValue("password")
	)

	user, err := c.verify(username, password)
	if err != nil {
		c.fireError(err)
		return
	}

	c.Session.Set(sessionIDKey, user.ID)
	c.Path = pathMyProfile
}

func (c *Controller) AnyLogout() {
	c.logout()
}

func (c *Controller) GetMe() {
	id, err := c.Session.GetInt64(sessionIDKey)
	if err != nil || id <= 0 {
		c.Path = PathLogin
		return
	}

	u, found := c.Source.GetByID(id)
	if !found {
		c.logout()
		return
	}

	c.User = u
	c.Data["Title"] = "Profile of " + u.Username
	c.Tmpl = pathMyProfile + ".html"
}

func (c *Controller) renderNotFound(id int64) {
	c.Status = 404
	c.Data["Title"] = "User Not Found"
	c.Data["ID"] = id
	c.Tmpl = "user/notfound.html"
}

func (c *Controller) GetBy(userID int64) {
	if user, found := c.Source.GetByID(userID); !found {
		c.renderNotFound(userID)
	} else {
		c.Ctx.JSON(user)
	}
}