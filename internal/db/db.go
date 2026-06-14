package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
)

func connect() (*sql.DB, error) {
	db, err := sql.Open("sqlite", "volumes/photo_bank")
	if err != nil {
		log.Println(err)
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		log.Println("ping db:", err)
		return nil, err
	}
	return db, nil
}

func Register(login string, password string) (message string, err error) {
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer conn.Close()

	_, err = conn.Exec("INSERT INTO users (login, password) VALUES ($1, $2)", login, password)
	if err != nil {
		log.Println(err)
		return "", err
	}
	return fmt.Sprintf("Пользователь %s добавлен", login), nil
}

func GetUserByLogin(login string) (*User, error) {
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()

	var user User
	err = conn.QueryRow("SELECT * FROM users WHERE login = $1", login).
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
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return err
	}
	defer conn.Close()

	query := `
        INSERT INTO user_images (user_id, filename, original_name, file_path, file_size, mime_type, description)
        VALUES (?, ?, ?, ?, ?, ?, ?)
    `
	result, err := conn.Exec(query,
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

func GetUserImages(userID int64) ([]Image, error) {
	conn, err := connect()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer conn.Close()

	query := `
       SELECT id, user_id, filename, original_name, file_path, file_size, mime_type, description, created_at
       FROM user_images
       WHERE user_id = ?
       ORDER BY created_at DESC
   `
	rows, err := conn.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		err := rows.Scan(
			&img.Id, &img.UserId, &img.Filename, &img.OriginalName,
			&img.FilePath, &img.FileSize, &img.MimeType, &img.Description,
			&img.CreatedAt)
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, nil
}

//// DeleteImage - удаляет изображение из БД
//func DeleteImage(imageID, userID int) error {
//	query := `DELETE FROM user_images WHERE id = ? AND user_id = ?`
//	_, err := DB.Exec(query, imageID, userID)
//	return err
//}

//// GetImageByID - получает изображение по ID
//func GetImageByID(imageID int) (*Image, error) {
//	query := `
//        SELECT id, user_id, filename, original_name, file_path, file_size, mime_type, description, created_at
//        FROM user_images
//        WHERE id = ?
//    `
//	var img Image
//	err := DB.QueryRow(query, imageID).Scan(
//		&img.ID, &img.UserID, &img.Filename, &img.OriginalName,
//		&img.FilePath, &img.FileSize, &img.MimeType, &img.Description,
//		&img.CreatedAt)
//	if err != nil {
//		return nil, err
//	}
//	return &img, nil
//}
