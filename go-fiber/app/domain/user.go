package domain

type User struct {
	Name     string `json:"name" bson:"name" validate:"required"`
	Username string `json:"username" bson:"username" validate:"required,min=6,username"`
	Password string `json:"password" bson:"password"  validate:"required,min=8,password"`
	RoleID   string `json:"role_id" bson:"role_id"`
	Role     Role
}
