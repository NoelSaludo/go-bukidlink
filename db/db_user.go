package db

func QueryUsers(username string) ([]User, error) {
	var temp []User
	query := `SELECT id, username, password, email, address
				FROM "User" 
				WHERE username=$1`

	rows, err := db.Query(query, username)
	if err != nil {
		return temp, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.Password,
			&user.Email,
			&user.Address)

		if err != nil {
			return temp, err
		}

		temp = append(temp, user)
	}

	return temp, err
}

func InsertUser(user User) error {
	query := `
	INSERT INTO "User" (
	id, username, password, email, address)
	VALUES ($1, $2, $3, $4, $5)`

	_, err := db.Exec(query,
		user.Id,
		user.Username,
		user.Password,
		user.Email,
		user.Address)

	if err != nil {
		return err
	}

	return err
}

func DeleteUser(id string) error {
	query := `DELETE FROM "User" WHERE id=$1`

	_, err := db.Exec(query, id)
	if err != nil {
		return err
	}

	return nil
}
