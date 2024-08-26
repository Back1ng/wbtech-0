package entity

import "time"

type TrackNumber string

type Delivery struct {
	Name    string `json:"name" db:"name"`
	Phone   string `json:"phone" db:"phone"`
	Zip     string `json:"zip" db:"zip"`
	City    string `json:"city" db:"city"`
	Address string `json:"address" db:"address"`
	Region  string `json:"region" db:"region"`
	Email   string `json:"email" db:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction" db:"transaction"`
	RequestID    string `json:"request_id,omitempty" db:"request_id"`
	Currency     string `json:"currency" db:"currency"`
	Provider     string `json:"provider" db:"provider"`
	Amount       int    `json:"amount" db:"amount"`
	PaymentDT    int    `json:"payment_dt" db:"payment_dt"`
	Bank         string `json:"bank" db:"bank"`
	DeliveryCost int    `json:"delivery_cost" db:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total" db:"goods_total"`
	CustomFee    int    `json:"custom_fee" db:"custom_fee"`
}

type Item struct {
	ChrtId      int         `json:"chrt_id" db:"chrt_id"`
	TrackNumber TrackNumber `json:"track_number" db:"track_number"`
	Price       int         `json:"price" db:"price"`
	Rid         string      `json:"rid" db:"rid"`
	Name        string      `json:"name" db:"name"`
	Sale        int         `json:"sale" db:"sale"`
	Size        string      `json:"size" db:"size"`
	TotalPrice  int         `json:"total_price" db:"total_price"`
	NmId        int         `json:"nm_id" db:"nm_id"`
	Brand       string      `json:"brand" db:"brand"`
	Status      int         `json:"status" db:"status"`
}

type Order struct {
	OrderUID          string      `json:"order_uid" db:"order_uid" validate:"required"`
	TrackNumber       TrackNumber `json:"track_number" db:"track_number" validate:"required"`
	Entry             string      `json:"entry" db:"entry" validate:"required"`
	Delivery          Delivery    `json:"delivery" db:"delivery" validate:"required"`
	Payment           Payment     `json:"payment" db:"payment" validate:"required"`
	Items             []Item      `json:"items" db:"items" validate:"required"`
	Locale            string      `json:"locale" db:"locale" validate:"required"`
	InternalSignature string      `json:"internal_signature,omitempty" db:"internal_signature"`
	CustomerId        string      `json:"customer_id" db:"customer_id" validate:"required"`
	DeliveryService   string      `json:"delivery_service" db:"delivery_service" validate:"required"`
	ShardKey          string      `json:"shardkey" db:"shard_key" validate:"required"`
	SmId              int         `json:"sm_id" db:"sm_id" validate:"required"`
	DateCreated       time.Time   `json:"date_created" db:"date_created" validate:"required"`
	OofShard          string      `json:"oof_shard" db:"oof_shard" validate:"required"`
}
