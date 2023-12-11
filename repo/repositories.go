package repo

import (
	"context"
	"database/sql"
	"fmt"
	"learn/httpserver/model"
	"log"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
)

type Repositories interface {
	// GetData() ([]model.User, error)
	GetData() ([]model.GetUser, error)
	// CreateData(model.User) (bool, error)
	CreateEmployee(model.User, pgx.Tx) error
	CreateEmployeeServicePair(int64, []int64, pgx.Tx) error
	DeleteData(string) (bool, error)
	UpdateData(model.User, string) (bool, error)
	CheckUserExist(model.Login) (string, error)
	CreateNewService(model.Service) (bool, error)
}

// get all data
func (u User) GetData() ([]model.GetUser, error) {

	getQuery := "SELECT e.id,e.email,e.name,e.age,e.address,p.serviceid,s.service FROM employee e JOIN employeeservicepair p ON e.id=p.userid JOIN services s ON s.sid=p.serviceid"

	rows, err := u.db.Query(context.Background(), getQuery)

	if err != nil {
		log.Fatal(err)
	}

	var userTags = make(map[int64]model.GetUser)

	for rows.Next() {
		var getUser model.GetUser
		var (
			serviceId *int64
			service   *string
		)
		rows.Scan(&getUser.Id, &getUser.Email, &getUser.Name, &getUser.Age, &getUser.Address, &serviceId, &service)

		if currentUser, ok := userTags[getUser.Id]; ok {
			//if avaialable for same user then append
			if serviceId != nil {
				currentUser.Service = append(currentUser.Service, model.Service{Sid: *serviceId, Service: *service})
			}
			userTags[getUser.Id] = currentUser
		} else {
			user := userTags[getUser.Id]
			user = model.GetUser{
				Id:      getUser.Id,
				Email:   getUser.Email,
				Name:    getUser.Name,
				Age:     getUser.Age,
				Address: getUser.Address,
			}

			if serviceId != nil {
				user.Service = []model.Service{
					model.Service{Sid: *serviceId, Service: *service},
				}
			}

			userTags[getUser.Id] = user
		}
	}

	var users = []model.GetUser{}

	for userId, userData := range userTags {
		users = append(users, model.GetUser{
			Id:      userId,
			Email:   userData.Email,
			Name:    userData.Name,
			Age:     userData.Age,
			Address: userData.Address,
			Service: userData.Service,
		})
	}

	defer rows.Close()
	return users, nil
}

// using transaction
func (u User) CreateEmployee(newUserData model.User, tx pgx.Tx) error {
	//call dal to get newUserData
	fmt.Println(newUserData)

	_, err := tx.Exec(context.Background(), "INSERT INTO employee(id,email,password,name,age,address)values($1,$2,$3,$4,$5,$6)", newUserData.Id, newUserData.Email, newUserData.Password, newUserData.Name, newUserData.Age, newUserData.Address)

	// createdTag, err := u.db.Exec(context.Background(), "INSERT INTO employee(id,email,password,name,age,address)values($1,$2,$3,$4,$5,$6)", newUserData.Id, newUserData.Email, newUserData.Password, newUserData.Name, newUserData.Age, newUserData.Address)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func generateEmployeeServicePairParametersToInsert(UserId int64, servicesId []int64) (string, []interface{}) {
	var standardsParameters []string
	parameterValues := []interface{}{}
	var count = 1

	for _, sid := range servicesId {
		dtoString := fmt.Sprintf("($%s,$%s)", strconv.Itoa(count), strconv.Itoa(count+1))
		standardsParameters = append(standardsParameters, dtoString)
		parameterValues = append(parameterValues, UserId, sid)
		count += 2
	}
	strParametersNames := strings.Join(standardsParameters, ",")

	return strParametersNames, parameterValues

}

// CreateEmployeeServicePair
func (u User) CreateEmployeeServicePair(UserId int64, ServicesId []int64, tx pgx.Tx) error {
	//call dal to get newUserData
	fmt.Println(UserId, ServicesId)

	strParametersNames, parameterValues := generateEmployeeServicePairParametersToInsert(UserId, ServicesId)
	//insert into employee-service-pair table
	_, err := tx.Exec(context.Background(), "INSERT INTO employeeservicepair(userid,serviceid)values "+strParametersNames, parameterValues...)

	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// CreateNewService
func (u User) CreateNewService(newServiceData model.Service) (bool, error) {

	createdTag, err := u.db.Exec(context.Background(), "INSERT INTO services(sid,service)values($1,$2)", newServiceData.Sid, newServiceData.Service)
	if err != nil {
		log.Fatal(err)
	}

	if createdTag.Insert() {
		return true, nil
	}

	return false, nil

}

// DeleteData
func (u User) DeleteData(id string) (bool, error) {
	deletedTag, err := u.db.Exec(context.Background(), "DELETE FROM employee WHERE id=$1", id)
	if err != nil {
		log.Fatal(err)
	}

	if deletedTag.Delete() {
		return true, nil
	}

	return false, nil

}

// UpdateData
func (u User) UpdateData(updateData model.User, id string) (bool, error) {
	updatedTag, err := u.db.Exec(context.Background(), "UPDATE employee SET name=$2,age=$3,address=$4,email=$5,sid=$6 WHERE id=$1", id, updateData.Name, updateData.Age, updateData.Address, updateData.Email, updateData.Sid)
	if err != nil {
		log.Fatal(err)
	}

	if updatedTag.Update() {
		return true, nil
	}

	return false, nil

}

func (u User) CheckUserExist(user model.Login) (string, error) {

	var name string
	row := u.db.QueryRow(context.Background(), "SELECT name FROM employee WHERE email=$1 AND password=$2", user.Email, user.Password)
	err := row.Scan(&name)
	if err != nil || err != sql.ErrNoRows {
		fmt.Println(err)
		return "", err
	}

	// //generate token
	// loggedInToken, tokenError := utils.GenerateToken(strconv.FormatInt(user.Id,10))
	// if tokenError != nil {
	// 	return "",tokenError
	// }

	return "success", nil

}

// func (u User) CreateData(newUserData model.User) (bool, error) {
// 	//call dal to get newUserData
// 	fmt.Println(newUserData)
// 	createdTag, err := u.db.Exec(context.Background(), "INSERT INTO employee(id,email,password,name,age,address)values($1,$2,$3,$4,$5,$6)", newUserData.Id, newUserData.Email, newUserData.Password, newUserData.Name, newUserData.Age, newUserData.Address)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if createdTag.Insert() {
// 		return true, nil
// 	}

// 	return false, nil

// }
