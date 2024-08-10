package verificationcode

type CreateForm struct {
	Type     string `json:"type" binding:"required"`
	Username string `json:"username" binding:"required"`
}

type UseForm struct {
	Type     string `json:"type" binding:"required"`
	Username string `json:"username" binding:"required"`
	Code     string `json:"code" binding:"required"`
	NewPass  string `json:"newPass"`
}
