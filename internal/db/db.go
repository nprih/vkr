package db

import (
	"database/sql"
	"fmt"
	"log"
	"vkr/internal/repository"

	_ "modernc.org/sqlite"
)

func Register(login string, password string) (message string, err error) {
	db := repository.GetDB()

	_, err = db.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf("Пользователь %s добавлен", login), nil
}

func GetUserByLogin(login string) (*User, error) {
	db := repository.GetDB()

	var user User
	err := db.QueryRow("SELECT * FROM users WHERE login = $1", login).
		Scan(&user.Id, &user.Login, &user.Password, &user.Is_admin)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("пользователь с Login %s не найден", login)
		}
		return nil, fmt.Errorf("ошибка при получении пользователя: %w", err)
	}

	return &user, nil
}

func SaveImageInfo(image *Image) error {
	db := repository.GetDB()

	query := `
        INSERT INTO user_images (user_id, filename, original_name, file_path, file_size, mime_type, description)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	result, err := db.Exec(query,
		image.UserId, image.Filename, image.OriginalName,
		image.FilePath, image.FileSize, image.MimeType, image.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err == nil {
		image.Id = id
	}
	return err
}

func GetAllUserImages() ([]UserImage, error) {
	db := repository.GetDB()

	query := `SELECT ui.id, ui.file_path, u.login author, ui.created_at 
				FROM users u 
				    INNER JOIN user_images ui ON u.id = ui.user_id 
				ORDER BY ui.created_at DESC`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []UserImage
	for rows.Next() {
		var img UserImage
		err := rows.Scan(&img.Id, &img.FilePath, &img.Author, &img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func GetUserImages(userId int64) ([]UserImage, error) {
	db := repository.GetDB()

	query := `SELECT ui.id, ui.file_path, ui.created_at, u.login author 
				FROM users u 
				    INNER JOIN user_images ui ON u.id = ui.user_id 
				WHERE u.id = ? 
				ORDER BY ui.created_at DESC`
	rows, err := db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []UserImage
	for rows.Next() {
		var img UserImage
		err := rows.Scan(&img.Id, &img.FilePath, &img.CreatedAt, &img.Author)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

func GetAllUsers() ([]User, error) {
	db := repository.GetDB()

	rows, err := db.Query("SELECT id, login, is_admin FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Login, &user.Is_admin)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func GetImageByID(imageID int) (*Image, error) {
	db := repository.GetDB()
	query := `SELECT * FROM user_images WHERE id = ?`
	var img Image
	err := db.QueryRow(query, imageID).Scan(
		&img.Id, &img.UserId, &img.Filename, &img.OriginalName, &img.FilePath,
		&img.FileSize, &img.MimeType, &img.Description, &img.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &img, nil
}

func DeleteImage(imageID int) error {
	db := repository.GetDB()
	query := `DELETE FROM user_images WHERE id = ?`
	_, err := db.Exec(query, imageID)
	return err
}

func DeleteUser(userId int) error {
	db := repository.GetDB()
	_, err := db.Exec("DELETE FROM user_images WHERE user_id = ?", userId)
	_, err = db.Exec("DELETE FROM users WHERE id = ?", userId)
	return err
}
