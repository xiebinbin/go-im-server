package relation

type Relation struct {
	ID         string `bson:"_id" json:"id"` //
	AK         string `bson:"ak" json:"ak,omitempty"`
	AUId       string `bson:"auid" json:"auid,omitempty"`
	TargetAUId string `bson:"target_auid" json:"target_auid,omitempty"`
	Status     int8   `bson:"status" json:"status,omitempty"`
	CreatedAt  int64  `bson:"create_time" json:"create_time,omitempty"`
	UpdatedAt  int64  `bson:"update_time" json:"update_time,omitempty"`
}

func (r Relation) TableName() string {
	return "user_relation"
}

func New() *Relation {
	return new(Relation)
}

func (r *Relation) WithAKey(ak string) *Relation {
	r.AK = ak
	return r
}
