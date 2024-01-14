package flow

import (
	"time"
)

type ActionModel interface {
	GetID() int64
	GetType() string
	GetObjectType() *string
	GetObjectId() *string
	GetObjectIdField() *string
	GetObjectRelation() *string
	GetSubjectType() *string
	GetSubjectId() *string
	GetSubjectIdField() *string
	GetScript() *string
	GetVerifiationRequired() bool
	GetOrderWeight() int64
	GetRemoveAction() bool
	GetActionEnabled() bool
	GetRunOnNewUser() bool
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	GetDeletedAt() *time.Time
	ToEventActionSpec() *EventActionsSpec
}

type EventAction struct {
	ID                  int64      `mysql:"id" postgres:"id" sqlite:"id"`
	Type                string     `mysql:"type" postgres:"type" sqlite:"type"`
	ObjectType          *string    `mysql:"objectType" postgres:"object_type" sqlite:"objectType"`
	ObjectId            *string    `mysql:"objectId" postgres:"object_id" sqlite:"objectId"`
	ObjectIdField       *string    `mysql:"objectIdField" postgres:"object_id_field" sqlite:"objectIdField"`
	ObjectRelation      *string    `mysql:"objectRelation" postgres:"object_relation" sqlite:"objectRelation"`
	SubjectType         *string    `mysql:"subjectType" postgres:"subject_type" sqlite:"subjectType"`
	SubjectId           *string    `mysql:"subjectId" postgres:"subject_id" sqlite:"subjectId"`
	SubjectIdField      *string    `mysql:"subjectIdField" postgres:"subject_id_field" sqlite:"subjectIdField"`
	Script              *string    `mysql:"script" postgres:"script" sqlite:"script"`
	VerifcationRequired bool       `mysql:"verificationRequired" postgres:"verification_required" sqlite:"verificationRequired"`
	OrderWeight         int64      `mysql:"orderWeight" postgres:"order_weight" sqlite:"orderWeight"`
	RemoveAction        bool       `mysql:"removeAction" postgres:"remove_action" sqlite:"removeAction"`
	ActionEnabled       bool       `mysql:"actionEnabled" postgres:"action_enabled" sqlite:"actionEnabled"`
	RunOnNewUser        bool       `mysql:"runOnNewUser" postgres:"run_on_new_user" sqlite:"runOnNewUser"`
	CreatedAt           time.Time  `mysql:"createdAt" postgres:"created_at" sqlite:"createdAt"`
	UpdatedAt           time.Time  `mysql:"updatedAt" postgres:"updated_at" sqlite:"updatedAt"`
	DeletedAt           *time.Time `mysql:"deletedAt" postgres:"deleted_at" sqlite:"deletedAt"`
}

func (event EventAction) GetID() int64 {
	return event.ID
}

func (event EventAction) GetType() string {
	return event.Type
}

func (event EventAction) GetObjectType() *string {
	return event.ObjectType
}

func (event EventAction) GetObjectId() *string {
	return event.ObjectId
}

func (event EventAction) GetObjectIdField() *string {
	return event.ObjectIdField
}

func (event EventAction) GetObjectRelation() *string {
	return event.ObjectRelation
}

func (event EventAction) GetSubjectType() *string {
	return event.SubjectType
}

func (event EventAction) GetSubjectId() *string {
	return event.SubjectId
}

func (event EventAction) GetSubjectIdField() *string {
	return event.SubjectIdField
}

func (event EventAction) GetScript() *string {
	return event.Script
}

func (event EventAction) GetVerifiationRequired() bool {
	return event.VerifcationRequired
}

func (event EventAction) GetOrderWeight() int64 {
	return event.OrderWeight
}

func (event EventAction) GetRemoveAction() bool {
	return event.RemoveAction
}

func (event EventAction) GetActionEnabled() bool {
	return event.ActionEnabled
}

func (event EventAction) GetRunOnNewUser() bool {
	return event.RunOnNewUser
}

func (event EventAction) GetCreatedAt() time.Time {
	return event.CreatedAt
}

func (event EventAction) GetUpdatedAt() time.Time {
	return event.UpdatedAt
}

func (event EventAction) GetDeletedAt() *time.Time {
	return event.DeletedAt
}

func (event EventAction) ToEventActionSpec() *EventActionsSpec {
	return &EventActionsSpec{
		Type:           event.Type,
		ObjectType:     *event.ObjectType,
		ObjectId:       *event.GetObjectId(),
		ObjectIdField:  *event.ObjectIdField,
		ObjectRelation: *event.ObjectRelation,
		SubjectType:    *event.SubjectType,
		SubjectId:      *event.SubjectId,
		SubjectIdField: *event.SubjectIdField,
		Script:         *event.Script,
		OrderWeight:    event.OrderWeight,
		RemoveAction:   event.RemoveAction,
		ActionEnabled:  event.ActionEnabled,
	}
}
