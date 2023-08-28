package video

type Video struct {
	VideoID       int64
	UserID        int64 `gorm:"foreignKey:UserID"`
	User          User
	PlayURL       string
	CoverURL      string
	FavoriteCount int64
	CommentCount  int64
	IsFavorite    bool
	Title         string
}

type User struct {
	ID              int64
	Name            string
	FollowCount     int64
	FollowerCount   int64
	IsFollow        bool
	Avatar          string
	BackgroundImage string
	Signature       string
	TotalFavorited  int64
	WorkCount       int64
	FavoriteCount   int64
}
