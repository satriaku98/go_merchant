package model

type Customer struct {
	CustomerId string `db:"customer_id" json:"customer_id,omitempty"`
	Name       string `db:"name" json:"name,omitempty"`
	Type       string `db:"type" json:"type,omitempty"`
	Username   string `db:"username" json:"username,omitempty,required"`
	Password   string `db:"password" json:"password,omitempty,required"`
	Balance    int    `db:"balance" json:"balance,omitempty"`
	Token      string `db:"token" json:"token,omitempty"`
}
