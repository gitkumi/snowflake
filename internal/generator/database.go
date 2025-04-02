package generator

import (
	"fmt"
	"path/filepath"

	snowflaketemplate "github.com/gitkumi/snowflake/template"
)

type Database string

const (
	SQLite3  Database = "sqlite3"
	Postgres Database = "postgres"
	MySQL    Database = "mysql"
)

var AllDatabases = []Database{
	SQLite3,
	Postgres,
	MySQL,
}

var (
	registerFunction = `func (s *Server) Register(c *gin.Context) {
	var registerParams RegisterInput
	if err := c.BindJSON(&registerParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := registerParams.Data()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.CreateUser(c, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}`

	loginFunction = `func (s *Server) Login(c *gin.Context) {
	var login LoginInput
	if err := c.BindJSON(&login); err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, login.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid username or password.",
		})
		return
	}

	if user == (data.User{}) {
		// Fake run
		_, _ = password.VerifyPassword(login.Password, "notsecure")
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid username or password.",
		})
		return
	}

	match, err := password.VerifyPassword(login.Password, user.HashedPassword.String)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid username or password.",
		})
		return
	}

	if !match {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Invalid username or password.",
		})
		return
	}

	tokenString, err := utils.CreateJwt(user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}`

	createMagicLinkFunction = `func (s *Server) CreateMagicLink(c *gin.Context) {
	var magicLink CreateMagicLinkInput
	if err := c.BindJSON(&magicLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, magicLink.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := s.Query.CreateUserAuthToken(c, data.CreateUserAuthTokenParams{
		ID:     nanoid,
		UserID: user.ID,
		Type:   "magic_link",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	loginURL := path.Join(s.BaseURL, "auth", "email-login", token.ID)
	s.Mailer.Send(smtp.Email{
		To: user.Email, Subject: "Login", Body: loginURL,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Please check your email for the magic link.",
	})
}`

	createResetPasswordFunction = `func (s *Server) CreateResetPassword(c *gin.Context) {
	var createResetPasswordParams CreateResetPasswordInput
	if err := c.BindJSON(&createResetPasswordParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, createResetPasswordParams.Email)
	if err != nil || user == (data.User{}) {
		c.Status(http.StatusNoContent)
		return
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := s.Query.CreateUserAuthToken(c, data.CreateUserAuthTokenParams{
		ID:     nanoid,
		Type:   "reset_password",
		UserID: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resetPasswordURL := path.Join(s.BaseURL, "auth", "reset-password", token.ID)
	s.Mailer.Send(smtp.Email{
		To: user.Email, Subject: "Password Reset", Body: resetPasswordURL,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email for the reset password link.",
	})
}`

	createConfirmEmailFunction = `func (s *Server) CreateConfirmEmail(c *gin.Context) {
	var createConfirmEmailParams CreateConfirmEmailInput
	if err := c.BindJSON(&createConfirmEmailParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, createConfirmEmailParams.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	token, err := s.Query.CreateUserAuthToken(c, data.CreateUserAuthTokenParams{
		ID:     nanoid,
		UserID: user.ID,
		Type:   "confirm_email",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	confirmEmailURL := path.Join(s.BaseURL, "auth", "confirm-email", token.ID)
	s.Mailer.Send(smtp.Email{
		To: user.Email, Subject: "Confirm Email", Body: confirmEmailURL,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email for the confirmation link.",
	})
}`

	consumeMagicLinkFunction = `func (s *Server) ConsumeMagicLink(c *gin.Context) {
	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   c.Param("token"),
		Type: "magic_link",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	tokenString, err := utils.CreateJwt(token.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}`

	consumeResetPasswordFunction = `func (s *Server) ConsumeResetPassword(c *gin.Context) {
	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   c.Param("token"),
		Type: "reset_password",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var resetPasswordParams ResetPasswordInput
	if err := c.BindJSON(&resetPasswordParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUser(c, token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if resetPasswordParams.Email != user.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token.",
		})
		return
	}

	hashed, err := password.HashPassword(resetPasswordParams.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	updated, err := s.Query.UpdateUser(c, data.UpdateUserParams{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		HashedPassword: sql.NullString{
			String: hashed,
			Valid:  true,
		},
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		ConfirmedAt: user.ConfirmedAt,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}`

	consumeConfirmEmailFunction = `func (s *Server) ConsumeConfirmEmail(c *gin.Context) {
	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   c.Param("token"),
		Type: "confirm_email",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUser(c, token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	confirmed, err := s.Query.UpdateUser(c, data.UpdateUserParams{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ConfirmedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": createUserResponse(confirmed),
	})
}`

	mysqlRegisterFunction = `func (s *Server) Register(c *gin.Context) {
	var registerParams RegisterInput
	if err := c.BindJSON(&registerParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	data, err := registerParams.Data()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	res, err := s.Query.CreateUser(c, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	id, err := res.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUser(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": user,
	})
}`

	mysqlCreateMagicLinkFunction = `func (s *Server) CreateMagicLink(c *gin.Context) {
	var magicLink CreateMagicLinkInput
	if err := c.BindJSON(&magicLink); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, magicLink.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = s.Query.CreateUserAuthToken(c, data.CreateUserAuthTokenParams{
		ID:     nanoid,
		UserID: user.ID,
		Type:   "magic_link",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   nanoid,
		Type: "magic_link",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	loginURL := path.Join(s.BaseURL, "auth", "email-login", token.ID)
	s.Mailer.Send(smtp.Email{
		To: user.Email, Subject: "Login", Body: loginURL,
	})

	c.JSON(http.StatusCreated, gin.H{
		"message": "Please check your email for the magic link.",
	})
}`

	mysqlCreateResetPasswordFunction = `func (s *Server) CreateResetPassword(c *gin.Context) {
	var createResetPasswordParams CreateResetPasswordInput
	if err := c.BindJSON(&createResetPasswordParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, createResetPasswordParams.Email)
	if err != nil || user == (data.User{}) {
		c.Status(http.StatusNoContent)
		return
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = s.Query.CreateUserAuthToken(c, data.CreateUserAuthTokenParams{
		ID:     nanoid,
		Type:   "reset_password",
		UserID: user.ID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   nanoid,
		Type: "reset_password",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	resetPasswordURL := path.Join(s.BaseURL, "auth", "reset-password", token.ID)
	s.Mailer.Send(smtp.Email{
		To: user.Email, Subject: "Password Reset", Body: resetPasswordURL,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email for the reset password link.",
	})
}`

	mysqlCreateConfirmEmailFunction = `func (s *Server) CreateConfirmEmail(c *gin.Context) {
	var createConfirmEmailParams CreateConfirmEmailInput
	if err := c.BindJSON(&createConfirmEmailParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUserByEmail(c, createConfirmEmailParams.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	nanoid, err := gonanoid.New()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	_, err = s.Query.CreateUserAuthToken(c, data.CreateUserAuthTokenParams{
		ID:     nanoid,
		UserID: user.ID,
		Type:   "confirm_email",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   nanoid,
		Type: "confirm_email",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	confirmEmailURL := path.Join(s.BaseURL, "auth", "confirm-email", token.ID)
	s.Mailer.Send(smtp.Email{
		To: user.Email, Subject: "Confirm Email", Body: confirmEmailURL,
	})

	c.JSON(http.StatusOK, gin.H{
		"message": "Please check your email for the confirmation link.",
	})
}`

	mysqlConsumeResetPasswordFunction = `func (s *Server) ConsumeResetPassword(c *gin.Context) {
	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   c.Param("token"),
		Type: "reset_password",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	var resetPasswordParams ResetPasswordInput
	if err := c.BindJSON(&resetPasswordParams); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUser(c, token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if resetPasswordParams.Email != user.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid token.",
		})
		return
	}

	hashed, err := password.HashPassword(resetPasswordParams.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.Query.UpdateUser(c, data.UpdateUserParams{
		ID:       user.ID,
		Email:    user.Email,
		Username: user.Username,
		HashedPassword: sql.NullString{
			String: hashed,
			Valid:  true,
		},
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		ConfirmedAt: user.ConfirmedAt,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	updated, err := s.Query.GetUser(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}`

	mysqlConsumeConfirmEmailFunction = `func (s *Server) ConsumeConfirmEmail(c *gin.Context) {
	token, err := s.Query.GetUserAuthTokens(c, data.GetUserAuthTokensParams{
		ID:   c.Param("token"),
		Type: "confirm_email",
	})

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := s.Query.GetUser(c, token.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = s.Query.UpdateUser(c, data.UpdateUserParams{
		ID:             user.ID,
		Email:          user.Email,
		Username:       user.Username,
		HashedPassword: user.HashedPassword,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		ConfirmedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	confirmed, err := s.Query.GetUser(c, user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": createUserResponse(confirmed),
	})
}`
)

func (d Database) String() string {
	return string(d)
}

func (d Database) IsValid() bool {
	for _, db := range AllDatabases {
		if db == d {
			return true
		}
	}
	return false
}

func (d Database) ConnString(projectName string) string {
	switch d {
	case SQLite3:
		return projectName + ".db"
	case Postgres:
		return fmt.Sprintf("user=postgres password=postgres dbname=%s host=localhost sslmode=disable", projectName)
	case MySQL:
		return ""
	default:
		return fmt.Sprintf("root:root@tcp(localhost:3306)/%s?parseTime=true", projectName)
	}
}

func (d Database) Driver() string {
	switch d {
	case SQLite3:
		return "sqlite3"
	case Postgres:
		return "postgres"
	case MySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) SQLCEngine() string {
	switch d {
	case SQLite3:
		return "sqlite"
	case Postgres:
		return "postgresql"
	case MySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) GooseDriver() string {
	switch d {
	case SQLite3:
		return "sqlite3"
	case Postgres:
		return "postgresql"
	case MySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) GooseDialect() string {
	switch d {
	case SQLite3:
		return "sqlite3"
	case Postgres:
		return "postgres"
	case MySQL:
		return "mysql"
	default:
		return ""
	}
}

func (d Database) GooseDBString(projectName string) string {
	switch d {
	case SQLite3:
		return projectName + ".db"
	case Postgres:
		return fmt.Sprintf("user=postgres password=postgres dbname=%s host=localhost sslmode=disable", projectName)
	case MySQL:
		return fmt.Sprintf("root:root@tcp(localhost:3306)/%s?parseTime=true", projectName)
	default:
		return ""
	}
}

func (d Database) Import() string {
	switch d {
	case SQLite3:
		return "github.com/mattn/go-sqlite3"
	case Postgres:
		return "github.com/lib/pq"
	case MySQL:
		return "github.com/go-sql-driver/mysql"
	default:
		return ""
	}
}

func LoadDatabaseMigration(db Database, filename string) (string, error) {
	fragmentPath := filepath.Join("fragments/database", string(db), "migrations", filename)
	content, err := snowflaketemplate.DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database fragment: %w", err)
	}
	return string(content), nil
}

func LoadDatabaseQuery(db Database, filename string) (string, error) {
	fragmentPath := filepath.Join("fragments/database", string(db), "queries", filename)
	content, err := snowflaketemplate.DatabaseFragments.ReadFile(fragmentPath)
	if err != nil {
		return "", fmt.Errorf("failed to read database query fragment: %w", err)
	}
	return string(content), nil
}

func (d Database) RegisterFunction() string {
	switch d {
	case SQLite3, Postgres:
		return registerFunction
	case MySQL:
		return mysqlRegisterFunction
	default:
		return ""
	}
}

func (d Database) LoginFunction() string {
	switch d {
	case SQLite3, Postgres:
		return loginFunction
	case MySQL:
		return loginFunction // MySQL uses the same implementation
	default:
		return ""
	}
}

func (d Database) CreateMagicLinkFunction() string {
	switch d {
	case SQLite3, Postgres:
		return createMagicLinkFunction
	case MySQL:
		return mysqlCreateMagicLinkFunction
	default:
		return ""
	}
}

func (d Database) CreateResetPasswordFunction() string {
	switch d {
	case SQLite3, Postgres:
		return createResetPasswordFunction
	case MySQL:
		return mysqlCreateResetPasswordFunction
	default:
		return ""
	}
}

func (d Database) CreateConfirmEmailFunction() string {
	switch d {
	case SQLite3, Postgres:
		return createConfirmEmailFunction
	case MySQL:
		return mysqlCreateConfirmEmailFunction
	default:
		return ""
	}
}

func (d Database) ConsumeMagicLinkFunction() string {
	switch d {
	case SQLite3, Postgres:
		return consumeMagicLinkFunction
	case MySQL:
		return consumeMagicLinkFunction // MySQL uses the same implementation
	default:
		return ""
	}
}

func (d Database) ConsumeResetPasswordFunction() string {
	switch d {
	case SQLite3, Postgres:
		return consumeResetPasswordFunction
	case MySQL:
		return mysqlConsumeResetPasswordFunction
	default:
		return ""
	}
}

func (d Database) ConsumeConfirmEmailFunction() string {
	switch d {
	case SQLite3, Postgres:
		return consumeConfirmEmailFunction
	case MySQL:
		return mysqlConsumeConfirmEmailFunction
	default:
		return ""
	}
}
