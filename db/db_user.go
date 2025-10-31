package db

import (
	"time"

	"github.com/google/uuid"
)

// QueryUser retrieves a single user and their details by username.
func QueryUser(username string) (User, error) {
	var user User
	query := `
		SELECT u.id,
		u.username,
		u.password,
		u.email,
		u.profile_pic_url,
		ud.address,
		ud.first_name,
		ud.last_name,
		ud.contact_number,
		ud.created_date
		FROM "User" u
		JOIN "UserUserDetail" uud ON u.id = uud.user_id
		JOIN "UserDetail" ud ON uud.detail_id = ud.id
		WHERE u.username = $1`

	row := db.QueryRow(query, username)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.ProfilePicPath,
		&user.Details.Address,
		&user.Details.FirstName,
		&user.Details.LastName,
		&user.Details.ContactNumber,
		&user.Details.CreatedDate,
	)

	if err != nil {
		return user, err
	}

	return user, nil
}

// InsertUser creates a new user, their details, and the link between them in a transaction.
func InsertUser(user User) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 1. Insert user details
	detailsId := uuid.New().String()
	detailsQuery := `
		INSERT INTO "UserDetail" (id, address, first_name, last_name, contact_number, created_date)
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.Exec(detailsQuery,
		detailsId,
		user.Details.Address,
		user.Details.FirstName,
		user.Details.LastName,
		user.Details.ContactNumber,
		time.Now())

	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// 2. Insert user
	userId := uuid.New().String()
	userQuery := `INSERT INTO "User" (id, username, password, email, profile_pic_url) VALUES ($1, $2, $3, $4, $5)`
	_, err = tx.Exec(userQuery, userId, user.Username, user.Password, user.Email, user.ProfilePicPath)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// 3. Link them in the UserUserDetail table
	linkQuery := `INSERT INTO "UserUserDetail" (user_id, detail_id) VALUES ($1, $2)`
	_, err = tx.Exec(linkQuery, userId, detailsId)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	return tx.Commit()
}

// DeleteUser deletes a user, their details, and the link in a transaction.
func DeleteUser(userId string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	// 1. Get the detail_id from the link table first
	var detailsId string
	linkQuery := `SELECT detail_id FROM "UserUserDetail" WHERE user_id = $1`
	err = tx.QueryRow(linkQuery, userId).Scan(&detailsId)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// 2. Delete from UserUserDetail link table first
	deleteLinkQuery := `DELETE FROM "UserUserDetail" WHERE user_id = $1`
	_, err = tx.Exec(deleteLinkQuery, userId)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// 3. Delete from User
	userQuery := `DELETE FROM "User" WHERE id = $1`
	_, err = tx.Exec(userQuery, userId)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	// 4. Delete from UserDetail
	detailsQuery := `DELETE FROM "UserDetail" WHERE id = $1`
	_, err = tx.Exec(detailsQuery, detailsId)
	if err != nil {
		return rollbackAndReturn(tx, err)
	}

	return tx.Commit()
}
