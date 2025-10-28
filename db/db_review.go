package db

func QueryReviewsOnItem(itemid string) ([]Review, error) {
	var reviews []Review

	query := `SELECT id, userid, itemid, content, rating
			FROM public."Review" WHERE itemid=$1`

	rows, err := db.Query(query, itemid)
	if err != nil {
		return nil, err
	}

	reviews, err = getReviewFromRow(rows, reviews)
	if err != nil {
		return nil, err
	}

	return reviews, nil
}

func QueryUsersItem(userid string) ([]Item, error) {
	var items []Item
	query := `
	SELECT i.id   AS itemid,
       i.name,
       i.description,
       i.costpkilo,
       i.category,
       i.amount
		FROM "User" u
		JOIN "UsersItem" ui ON u.id = ui.userid
		JOIN "Item" i ON i.id = ui.itemid
		WHERE u.id = $1; `

	rows, err := db.Query(query, userid)
	if err != nil {
		return nil, err
	}

	items = getItemsFromRow(rows, items)

	return items, nil
}
