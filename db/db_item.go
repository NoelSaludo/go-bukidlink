package db

import "github.com/google/uuid"

func QueryItemByID(id string) (Item, error) {

	query := `SELECT i.id as itemid,
			i.name,
			i.description,
			i.costpkilo,
			i.category,
			i.amount,
			i.img_path,
			AVG(c.rating) AS rating
			FROM "Item" AS i
			JOIN "Review" AS C ON c.itemid = i.id
			WHERE i.id=$1
			GROUP BY i.id, i.name
			`

	row := db.QueryRow(query, id)

	var item Item
	err := row.Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.CostPKilo,
		&item.Category,
		&item.Amount,
		&item.ImgPath,
		&item.Rating,
	)

	if err != nil {
		return Item{}, err
	}

	return item, err
}

func QueryAllItem100(block int) ([]Item, error) {
	var items []Item
	// select 100 items with offset
	query := `SELECT id,
			name,
			description,
			costpkilo,
			category,
			amount,
			img_path
			FROM "Item" LIMIT 100 OFFSET $1`
	rows, err := db.Query(query, block*100)
	if err != nil {
		return items, err
	}

	defer rows.Close()

	items = getItemsFromRow(rows, items)
	return items, err
}

// category only accepts certian string
// accepted string
// fruits, vegetables, grains, livestock, dairy, others
func QueryItembyCategory(category string) ([]Item, error) {
	var items []Item
	query := `SELECT id,
				name,
				description,
				costpkilo,
				category,
				amount,
				img_path 
				FROM public."Item" WHERE category=$1`

	rows, err := db.Query(query, category)
	if err != nil {
		return nil, err
	}

	items = getItemsFromRow(rows, items)
	return items, nil
}

func InsertItem(item Item) (string, error) {
	itemId := uuid.New().String()
	query := `INSERT INTO "Item" (id, name, description, amount, costpkilo, category, img_path)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.Exec(query,
		itemId,
		item.Name,
		item.Description,
		item.Amount,
		item.CostPKilo,
		item.Category,
		item.ImgPath,
	)
	if err != nil {
		return "", err
	}

	return itemId, nil
}
