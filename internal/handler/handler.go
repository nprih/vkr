package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"vkr/internal/db"
	"vkr/internal/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var req Credentials
var value Val

type Credentials struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Val struct {
	Login   string
	IsLogin bool
	IsAdmin bool
}

type Img struct {
	Id        int64  `json:"id"`
	Url       string `json:"url"`
	Author    string `json:"author"`
	CreatedAt string `json:"createdAt"`
}

func setValue(c *gin.Context) {
	session := sessions.Default(c)
	login, ok := session.Get("login").(string)
	is_login, ok := session.Get("is_login").(bool)
	is_admin, ok := session.Get("is_admin").(bool)
	if !ok {
		log.Println("Значения не установлены")
	}

	value = Val{
		Login:   login,
		IsLogin: is_login,
		IsAdmin: is_admin,
	}
}

func formatValues(images []db.UserImage) map[string]any {
	return map[string]any{
		"login":    value.Login,
		"is_login": value.IsLogin,
		"is_admin": value.IsAdmin,
		"images":   formatImages(images),
	}
}

func formatImages(images []db.UserImage) (urls []Img) {
	for _, image := range images {
		urls = append(urls, Img{
			Id:        image.Id,
			Url:       image.FilePath,
			Author:    image.Author,
			CreatedAt: image.CreatedAt.Format("02.01.2006 15:04:05"),
		})
	}
	return urls
}

func IndexHandler(c *gin.Context) {
	setValue(c)

	images, err := db.GetAllUserImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}
	c.HTML(http.StatusOK, "main.html", gin.H{
		"values": formatValues(images),
	})
}

func GetLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{})
}

func PostLoginHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("New login try:", req.Login, req.Password)
	user, err := db.GetUserByLogin(req.Login)
	if err != nil {
		log.Printf("Пользователь %s не найден", req.Login)
		c.JSON(http.StatusOK, gin.H{
			"fail": "Пользователь с указанным логин/паролем не найден",
		})
		return
	}
	if service.CheckPassword(req.Password, user.Password) {
		log.Println("Пароль совпадает")

		service.SetSession(c, user)

		c.JSON(http.StatusOK, gin.H{
			"message":  fmt.Sprintf("Пользователь %s авторизован", req.Login),
			"redirect": "/",
		})
	} else {
		log.Println("Не верный пароль")
		c.JSON(http.StatusOK, gin.H{
			"fail": "Пользователь с указанным логин/паролем не найден",
		})
	}
}

func GetRegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.html", gin.H{})
}

func PostRegisterHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Не удалось прочитать отправленные данные",
		})
		return
	}
	log.Println("New registry try:", req.Login, req.Password)

	password, err := service.HashPassword(req.Password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось хэшировать пароль",
		})
		return
	}
	message, err := db.Register(req.Login, password)
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Не удалось создать пользователя",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message":  message,
		"redirect": "/login",
	})
}

func GetAdminHandler(c *gin.Context) {
	setValue(c)

	images, err := db.GetAllUserImages()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	users, err := db.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println(users)
	c.HTML(http.StatusOK, "admin.html", gin.H{
		"values": formatValues(images),
		"users":  users,
	})
}

func PostLogoutHandler(c *gin.Context) {
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	log.Println("New logout try")

	service.UnsetSession(c)

	c.JSON(http.StatusOK, gin.H{
		"redirect": "/",
	})
}

func GetProfileHandler(c *gin.Context) {
	setValue(c)
	if !value.IsLogin || value.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	user, err := db.GetUserByLogin(value.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Пользователь не найден",
		})
		return
	}
	images, err := db.GetUserImages(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.HTML(http.StatusOK, "profile.html", gin.H{
		"values": formatValues(images),
	})
}

func PostUploadHandler(c *gin.Context) {
	setValue(c)
	if !value.IsLogin || value.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	user, err := db.GetUserByLogin(value.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Пользователь не найден",
		})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Файл не загружен"})
		return
	}

	uploaded, err := service.SaveUserImage(file, user.Id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	image := &db.Image{
		UserId:       user.Id,
		Filename:     uploaded.Filename,
		OriginalName: uploaded.OriginalName,
		FilePath:     uploaded.FilePath,
		FileSize:     int(uploaded.FileSize),
		MimeType:     uploaded.MimeType,
		Description:  "",
	}

	if err := db.SaveImageInfo(image); err != nil {
		service.DeleteUserImage(uploaded.FilePath)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось сохранить информацию об изображении"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "Изображение успешно загружено",
		"image_id":  image.Id,
		"image_url": uploaded.URL,
		"redirect":  "/gallery",
	})
}

func GetImagesHandler(c *gin.Context) {
	setValue(c)
	if !value.IsLogin || value.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	user, err := db.GetUserByLogin(value.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Пользователь не найден",
		})
		return
	}

	var images []db.UserImage
	if value.IsAdmin {
		images, err = db.GetAllUserImages()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		images, err = db.GetUserImages(user.Id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.HTML(http.StatusOK, "photosBody.html", gin.H{
		"values": formatValues(images),
	})
}

func GetImagesAdminHandler(c *gin.Context) {
	setValue(c)
	if !value.IsLogin || value.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	userId := c.Param("user_id")
	intUserId, err := strconv.Atoi(userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var images []db.UserImage
	if intUserId == 0 {
		images, err = db.GetAllUserImages()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	} else {
		images, err = db.GetUserImages(int64(intUserId))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
	}

	c.HTML(http.StatusOK, "photosBody.html", gin.H{
		"values": formatValues(images),
	})
}

func DeleteImagesHandler(c *gin.Context) {
	setValue(c)
	if !value.IsLogin || value.Login == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	user, err := db.GetUserByLogin(value.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Пользователь не найден",
		})
		return
	}

	imageId := c.Param("image_id")
	intImageId, err := strconv.Atoi(imageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	image, err := db.GetImageByID(intImageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = db.DeleteImage(intImageId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = os.Remove(image.FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "file not found"})
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	images, err := db.GetUserImages(user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	c.HTML(http.StatusOK, "photosBody.html", gin.H{
		"values": formatValues(images),
	})
}
