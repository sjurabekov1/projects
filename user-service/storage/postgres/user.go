package postgres

import (
	"github.com/jmoiron/sqlx"
	pb "github.com/template-service/genproto"
)

type userRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(user *pb.User) (*pb.User, error) {
	var (
		ruser = pb.User{}
	)
	insertUserQuery := `INSERT INTO users (id, first_name, last_name) values($1, $2, $3) RETURNING id, first_name, last_name`
	err := r.db.QueryRow(insertUserQuery, user.Id, user.FirstName, user.LastName).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
	)
	if err != nil {
		return &pb.User{}, err
	}
	return &ruser, nil
}

func (r *userRepo) UpdateUser(user *pb.User) (string, error) {

	insertUserQuery := `UPDATE users SET first_name = $1, last_name = $2 WHERE id = $3 `
	_, err := r.db.Query(insertUserQuery, user.FirstName, user.LastName, user.Id)
	if err != nil {
		return "", err
	}
	return "", err
}

func (r *userRepo) GetUserById(ID string) (*pb.User, error) {
	var ruser pb.User
	getByIdQuery := `SELECT id, first_name, last_name FROM users WHERE id = $1`
	err := r.db.QueryRow(getByIdQuery, ID).Scan(
		&ruser.Id,
		&ruser.FirstName,
		&ruser.LastName,
	)
	if err != nil {
		return &pb.User{}, err
	}

	return &ruser, nil
}

func (r *userRepo) GetUserList(limit, page int64) ([]*pb.User, int64, error) {

	var users []*pb.User
	var count int64

	offset := (page - 1) * limit

	queryGetUserList := `SELECT id, first_name, last_name FROM users OFFSET $1 LIMIT $2`

	rows, err := r.db.Query(queryGetUserList, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var user pb.User
		err := rows.Scan(
			&user.Id,
			&user.FirstName,
			&user.LastName,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, &user)
	}

	countQuery := `SELECT count(*) from users`
	err = r.db.QueryRow(countQuery).Scan(&count)

	return users, count, nil
}
