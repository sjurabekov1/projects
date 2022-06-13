package postgres

import (
	"github.com/gofrs/uuid"
	"github.com/jmoiron/sqlx"
	pb "github.com/template-service/genproto"
)

type postRepo struct {
	db *sqlx.DB
}

//NewUserRepo ...
func NewPostRepo(db *sqlx.DB) *postRepo {
	return &postRepo{db: db}
}
func (r *postRepo) CreatePost(post *pb.Post) (*pb.Post, error) {
	var (
		rPost = pb.Post{}
	)
	insertPost := `INSERT INTO posts (id, name, description, user_id) VALUES($1, $2, $3, $4) RETURNING id, name, description, user_id`
	err := r.db.QueryRow(insertPost, post.Id, post.Name, post.Description, post.UserId).Scan(
		&rPost.Id,
		&rPost.Name,
		&rPost.Description,
		&rPost.UserId,
	)
	if err != nil {
		return &pb.Post{}, err
	}
	mediaQuery := `INSERT INTO post_medias (id, type, link, post_id) VALUES($1, $2, $3, $4)`
	for _, media := range post.Medias {
		id, err := uuid.NewV4()
		if err != nil {
			return nil, err
		}
		_, err = r.db.Exec(mediaQuery, id, media.Type, media.Link, rPost.Id)
		if err != nil {
			return &pb.Post{}, err
		}
	}

	return &rPost, nil
}

func (r *postRepo) GetPostById(ID string) (*pb.Post, error) {
	var (
		rPost = pb.Post{}
	)
	getPostById := `SELECT id, name, description, user_id from posts WHERE id = $1`
	row, err := r.db.Query(getPostById, ID)
	for row.Next() {
		row.Scan(
			&rPost.Id,
			&rPost.Name,
			&rPost.Description,
			&rPost.UserId,
		)
		if err != nil {
			return nil, err
		}

		var medias []*pb.Media
		getAllMedias := `SELECT id, type, link from post_medias WHERE post_id = $1`
		rows, err := r.db.Query(getAllMedias, rPost.Id)

		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err := rows.Scan(
				&media.Id,
				&media.Type,
				&media.Link,
			)
			if err != nil {
				return nil, err
			}

			medias = append(medias, &media)
		}
		rPost.Medias = medias
	}

	return &rPost, nil
}

func (r *postRepo) GetAllUserPosts(userID string) ([]*pb.Post, error) {
	var (
		posts []*pb.Post
	)
	getAllUserPosts := `SELECT id, name, description, user_id from posts WHERE user_id = $1`
	rows, err := r.db.Query(getAllUserPosts, userID)

	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var post pb.Post
		err := rows.Scan(
			&post.Id,
			&post.Name,
			&post.Description,
			&post.UserId,
		)
		if err != nil {
			return nil, err
		}

		var medias []*pb.Media
		getAllMedias := `SELECT id, type, link from post_medias WHERE post_id = $1`
		rows, err := r.db.Query(getAllMedias, post.Id)

		if err != nil {
			return nil, err
		}
		for rows.Next() {
			var media pb.Media
			err := rows.Scan(
				&media.Id,
				&media.Type,
				&media.Link,
			)
			if err != nil {
				return nil, err
			}

			post.Medias = append(medias, &media)
		}
		posts = append(posts, &post)
	}

	return posts, nil
}

func (r *postRepo) GetAllPostList(limit, page int64) ([]*pb.Post, int64, error) {
	var posts []*pb.Post
	var medias []*pb.Media
	var count int64

	offset := (page - 1) * limit

	queryGetUserList := `SELECT id, name, description, user_id FROM posts OFFSET $1 LIMIT $2`

	rows, err := r.db.Query(queryGetUserList, offset, limit)
	if err != nil {
		return nil, 0, err
	}
	for rows.Next() {
		var post pb.Post
		err := rows.Scan(
			&post.Id,
			&post.Name,
			&post.Description,
			&post.UserId,
		)
		if err != nil {
			return nil, 0, err
		}

		getAllMedias := `SELECT id, type, link from post_medias WHERE post_id = $1`
		row, err := r.db.Query(getAllMedias, post.Id)
		if err != nil {
			return nil, 0, err
		}
		for row.Next() {
			var media pb.Media
			err := row.Scan(
				&media.Id,
				&media.Type,
				&media.Link,
			)
			if err != nil {
				return nil, 0, err
			}

			post.Medias = append(medias, &media)
		}

		posts = append(posts, &post)
	}

	countQuery := `SELECT count(*) from posts`
	err = r.db.QueryRow(countQuery).Scan(&count)

	return posts, count, nil
}
