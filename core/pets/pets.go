package pets

import "time"

type Pet struct {
	Added   time.Time `firestore:"added" json:"added"`
	Caption string    `firestore:"caption" json:"caption"`
	Email   string    `firestore:"email" json:"email"`
	Image   string    `firestore:"image" json:"image"`
	Likes   int       `firestore:"likes" json:"likes"`
	Owner   string    `firestore:"owner" json:"owner"`
	Petname string    `firestore:"petname" json:"petname"`
	Name    string
}

type Update struct {
	Added   *time.Time `json:"added" `
	Caption *string    `json:"caption" `
	Email   *string    `json:"email" `
	Image   *string    `json:"image" `
	Likes   *int       `json:"likes" `
	Owner   *string    `json:"owner" `
	Petname *string    `json:"petname" `
}

type Create struct {
	Added   time.Time `firestore:"added" json:"added" binding:"required"`
	Caption string    `firestore:"caption" json:"caption" binding:"required"`
	Email   string    `firestore:"email" json:"email" binding:"required"`
	Image   string    `firestore:"image" json:"image" binding:"required"`
	Likes   int       `firestore:"likes" json:"likes" binding:"required"`
	Owner   string    `firestore:"owner" json:"owner" binding:"required"`
	Petname string    `firestore:"petname" json:"petname" binding:"required"`
}
