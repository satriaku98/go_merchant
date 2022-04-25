package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"go_merchant/model"
)

type CustomerRepo interface {
	LogoutCustomer(username string, passcode string) (int, error)
	LoginCustomer(username string, passcode string) (model.Customer, int, error)
	SendTransfer(senderAccountNumber int, amountTransfer int) error
	GetTransfer(accountNumber int, amountTransfer int) error
}

func (c *loginRepo) LoginCustomer(username string, password string) (model.Customer, int, error) {
	var count int
	var data model.Customer
	fmt.Println(username, password)
	err := c.db.Get(&count, "select count(*) from customer where username = $1 and password = $2", username, password)
	fmt.Println(username, password, err, count)
	if err != nil {
		return data, 0, err
	}
	c.db.Get(&data, "select name, type from customer where username = $1 and password = $2", username, password)
	return data, count, nil
}

func (c *loginRepo) LogoutCustomer(username string, passcode string) (int, error) {
	var count int
	err := c.db.Get(&count, "select count(*) from customer where username = $1 and password = $2", username, passcode)
	fmt.Println(username, passcode, err, count)
	if err != nil {
		return 0, err
	}
	_, errUpdate := c.db.Query("update customer set token = null where username = $1", username)
	fmt.Println(username, passcode, errUpdate, count)
	if errUpdate != nil {
		return 0, errUpdate
	}
	return count, nil
}

func (c *loginRepo) SendTransfer(senderAccountNumber int, amountTransfer int) error {
	_, err := c.db.Exec("UPDATE customers SET balance = balance - $1 WHERE account_number = $2", amountTransfer, senderAccountNumber)
	if err != nil {
		return err
	}
	return nil
}

func (c *loginRepo) GetTransfer(accountNumber int, amountTransfer int) error {
	_, err := c.db.Exec("UPDATE customers SET balance = balance + $1 WHERE account_number = $2", amountTransfer, accountNumber)
	if err != nil {
		return err
	}
	return nil
}

type loginRepo struct {
	db *sqlx.DB
}

func NewLoginRepo(db *sqlx.DB) CustomerRepo {
	return &loginRepo{db}
}
