package handlers

// PasswordRequirements is a minimal swagger-friendly representation of password rules.
type PasswordRequirements struct {
	MinLength      int  `json:"min_length" example:"8"`
	RequireNumber  bool `json:"require_number"`
	RequireUpper   bool `json:"require_upper"`
	RequireSpecial bool `json:"require_special"`
}

// UserSettings is a minimal swagger-friendly representation mirroring internal models.UserSettings
type UserSettings struct {
	ID       uint   `json:"id" example:"1"`
	UserID   uint   `json:"user_id" example:"1"`
	Language string `json:"language" example:"en"`
	Timezone string `json:"timezone" example:"UTC"`
	Theme    string `json:"theme" example:"light"`
	Settings string `json:"settings" example:"{}"`
}

// UserSettingsUpdateRequest minimal representation for swagger
type UserSettingsUpdateRequest struct {
	Language string `json:"language" example:"en"`
	Timezone string `json:"timezone" example:"UTC"`
	Theme    string `json:"theme" example:"light"`
	Settings string `json:"settings" example:"{}"`
}
