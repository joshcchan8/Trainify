package queries

// User Queries

var GetUserQuery = `
	SELECT 
		user_id,
		username,
		email,
		hashed_password,
		profile_id 
	FROM users 
	WHERE user_id=?
`

var GetUserByEmailQuery = `
	SELECT
		user_id,
		username,
		email,
		hashed_password,
		profile_id
	FROM users 
	WHERE email=?
`

var CreateUserQuery = `
	INSERT INTO users 
		(username, email, hashed_password, profile_id)
	VALUES (?, ?, ?, ?)
`

var UpdateUserQuery = `
	UPDATE users
	SET 
		username=?,
		hashed_password=?
	WHERE user_id=?
`

var DeleteUserQuery = `
	DELETE FROM users 
	WHERE user_id=?
`

// Profile Queries

var GetProfileIdFromUserIdQuery = `
	SELECT 
		profile_id 
	FROM users 
	WHERE user_id=?
`

var GetProfileQuery = `
	SELECT 
		profile_id,
		age,
		weight,
		height,
		max_push_ups,
		avg_push_ups,
		max_pull_ups,
		avg_pull_ups,
		max_squat,
		avg_squat,
		max_bench,
		avg_bench,
		cardio_level
	FROM user_profiles 
	WHERE profile_id=?
`

var CreateEmptyProfileQuery = `
	INSERT INTO user_profiles
		(
			age,
			weight,
			height,
			max_push_ups,
			avg_push_ups,
			max_pull_ups, 
			avg_pull_ups,
			max_squat,
			avg_squat,
			max_bench,
			avg_bench,
			cardio_level
		)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
`

var UpdateProfileQuery = `
	UPDATE user_profiles
	SET 
		age=?,
		weight=?,
		height=?,
		max_push_ups=?,
		avg_push_ups=?,
		max_pull_ups=?,
		avg_pull_ups=?,
		max_squat=?,
		avg_squat=?,
		max_bench=?,
		avg_bench=?,
		cardio_level=?
	WHERE profile_id=?
`

var DeleteProfileQuery = `
	DELETE FROM user_profiles 
	WHERE profile_id=?
`

// Item Queries

var GetItemQuery = `
	SELECT 
		item_id,
		item_name,
		difficulty,
		minutes,
		calories_burned,
		targeted_muscle_groups,
		workout_description,
		created_by
	FROM items 
	WHERE item_id=? AND created_by=?`

var GetAllItemsQuery = `
	SELECT 
		item_id,
		item_name,
		difficulty,
		minutes,
		calories_burned,
		targeted_muscle_groups,
		workout_description,
		created_by 
	FROM items 
	WHERE created_by=?
`

var CreateItemQuery = `
	INSERT INTO items 
		(
			item_name, 
			difficulty, 
			minutes, 
			calories_burned, 
			targeted_muscle_groups, 
			workout_description, 
			created_by
		) 
	VALUES (?, ?, ?, ?, ?, ?, ?)`

var UpdateItemQuery = `
	UPDATE items 
	SET 
		item_name=?,
		difficulty=?,
		minutes=?,
		calories_burned=?,
		targeted_muscle_groups=?, 
		workout_description=?,
		created_by=?
	WHERE item_id=?
`

var DeleteItemQuery = `
	DELETE FROM items
	WHERE item_id=?
`

var DeleteAllItemsQuery = `
	DELETE FROM items 
	WHERE created_by=?
`
