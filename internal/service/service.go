package service

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"vkr/internal/db"

	"github.com/alexedwards/argon2id"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	MaxImageSize  = 10 * 1024 * 1024 // 10 MB
	BaseUploadDir = "uploads/users"
)

var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/jpg":  true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

type UploadedImage struct {
	Filename     string
	OriginalName string
	FilePath     string
	FileSize     int64
	MimeType     string
	URL          string
}

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		log.Fatal(err)
	}
	return hash, nil
}

func CheckPassword(password string, hash string) bool {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		log.Fatal(err)
	}
	return match
}

func SetSession(c *gin.Context, user *db.User) {
	session := sessions.Default(c)
	session.Set("login", user.Login)
	session.Set("is_login", true)
	session.Set("is_admin", user.Is_admin)
	session.Save()
}

func UnsetSession(c *gin.Context) {
	session := sessions.Default(c)
	session.Delete("login")
	session.Delete("is_login")
	session.Delete("is_admin")
	session.Save()
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		login := session.Get("login")
		if login == nil {
			log.Println("Пользователь не авторизован, редирект на /login")
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		isAdmin, ok := session.Get("is_admin").(bool)
		if !ok || !isAdmin {
			log.Println(" Админ не авторизован, редирект на /")
			c.HTML(http.StatusForbidden, "error.html", gin.H{
				"error":    "Доступ запрещен",
				"message":  "Требуются права администратора",
				"redirect": "/",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func SaveUserImage(fileHeader *multipart.FileHeader, userID int64) (*UploadedImage, error) {
	if fileHeader.Size > MaxImageSize {
		return nil, fmt.Errorf("файл слишком большой. Максимальный размер: %d MB", MaxImageSize/1024/1024)
	}

	file, err := fileHeader.Open()
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть файл: %v", err)
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}

	mimeType := http.DetectContentType(buffer)
	if !AllowedImageTypes[mimeType] {
		return nil, fmt.Errorf("неподдерживаемый тип файла. Разрешены: JPEG, PNG, GIF, WEBP")
	}

	file.Seek(0, 0)

	userDir := filepath.Join(BaseUploadDir, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("не удалось создать директорию: %v", err)
	}

	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		switch mimeType {
		case "image/jpeg", "image/jpg":
			ext = ".jpg"
		case "image/png":
			ext = ".png"
		case "image/gif":
			ext = ".gif"
		case "image/webp":
			ext = ".webp"
		}
	}

	timestamp := time.Now().Unix()
	uniqueID := uuid.New().String()[:8]
	filename := fmt.Sprintf("%d_%s%s", timestamp, uniqueID, ext)

	filePath := filepath.Join(userDir, filename)
	webPath := "/" + filepath.ToSlash(filePath) // для URL

	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("не удалось создать файл: %v", err)
	}
	defer dst.Close()

	written, err := io.Copy(dst, file)
	if err != nil {
		return nil, fmt.Errorf("ошибка сохранения файла: %v", err)
	}

	log.Printf("Изображение сохранено: %s (пользователь %d, размер %d байт)", filePath, userID, written)

	return &UploadedImage{
		Filename:     filename,
		OriginalName: fileHeader.Filename,
		FilePath:     filePath,
		FileSize:     written,
		MimeType:     mimeType,
		URL:          webPath,
	}, nil
}

func DeleteUserImage(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil // Файла нет
	}

	err := os.Remove(filePath)
	if err != nil {
		log.Printf("Ошибка удаления файла %s: %v", filePath, err)
		return err
	}

	log.Printf("Файл изображения удален: %s", filePath)
	return nil
}
