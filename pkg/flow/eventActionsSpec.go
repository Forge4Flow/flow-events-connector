package flow

// TODO: Better validation for required fields
type EventActionsSpec struct {
	Type                 string `json:"type,omitempty"` // TODO: Rename to have consistency with Flow event naming conventions
	ObjectType           string `json:"objectType,omitempty" validate:"required_with=ActionEnabled"`
	ObjectId             string `json:"objectId,omitempty"`
	ObjectIdField        string `json:"objectIdField,omitempty"`
	ObjectRelation       string `json:"objectRelation,omitempty"`
	SubjectType          string `json:"subjectType,omitempty"`
	SubjectId            string `json:"subjectId"`
	SubjectIdField       string `json:"subjectIdField"`
	Script               string `json:"script,omitempty"`
	VerificationRequired bool   `json:"verificationRequired,omitempty"`
	OrderWeight          int64  `json:"orderWeight" validate:"required"`
	RemoveAction         bool   `json:"removeAction,omitempty"`
	ActionEnabled        bool   `json:"actionEnabled,omitempty"`
	RunOnNewUser         bool   `json:"runOnNewUser,omitempty"`
}

func (e *EventActionsSpec) ToEventAction() EventAction {
	return EventAction{
		Type:                e.Type,
		ObjectType:          &e.ObjectType,
		ObjectId:            &e.ObjectId,
		ObjectIdField:       &e.ObjectIdField,
		ObjectRelation:      &e.ObjectRelation,
		SubjectType:         &e.SubjectType,
		SubjectId:           &e.SubjectId,
		SubjectIdField:      &e.SubjectIdField,
		Script:              &e.Script,
		VerifcationRequired: e.VerificationRequired,
		OrderWeight:         e.OrderWeight,
		RemoveAction:        e.RemoveAction,
		ActionEnabled:       e.ActionEnabled,
		RunOnNewUser:        e.RunOnNewUser,
	}
}
