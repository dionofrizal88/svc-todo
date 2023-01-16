package services

import (
	"context"
	"log"
	pbUser "svc-todo/pb/user"
	"svc-todo/server/db"
)

type UserService struct {
	pbUser.UnimplementedUserServiceServer
}

func (p *UserService) Get(context.Context, *pbUser.Empty) (*pbUser.UserResponse, error){
	log.Println("Get User Was Invoked")

	var users []*pbUser.User

	con := db.CreateCon()	

	stmt := "SELECT u.name, u.uuid from user u"

	rows, err := con.Query(stmt)
	if err != nil{
		return nil, err
	}
	defer rows.Close()
	for rows.Next(){
		var user pbUser.User
		err = rows.Scan(&user.Name, &user.Uuid)
		if err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	response := &pbUser.UserResponse{
		Status: true,
		Data: users,
	}

	return response, nil
}

func (p *UserService) Add(ctx context.Context, user *pbUser.AddUserRequest) (*pbUser.AddUserResponse, error){
	log.Println("Add User Was Invoked")

	con := db.CreateCon()

	sqlStatement := "INSERT user (name) VALUES (?)"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(user.GetName())
	if err != nil {
		return nil, err
	}

	response := &pbUser.AddUserResponse{
		Status: true,
		Message: "Success",
	}

	return response, nil
}

func FindUserData(uuid string) (int64, error){
	con := db.CreateCon()

	sqlStatement := "SELECT id FROM user WHERE uuid=?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return 0, err
	}

	result, err := stmt.Query(uuid)
	if err != nil {
		return 0, err
	}

	var id int64

	for result.Next() {   
		err = result.Scan(&id)
		if err != nil {
			return 0, err
		}
	}

	return id, nil
}

	// users := &pbUser.UserResponse{
	// 	Status: true,
	// 	Data: []*pbUser.User{
	// 		{
	// 			Id: 1,
	// 			Name: "Dio",
	// 		},
	// 		{
	// 			Id: 2,
	// 			Name: "Dani",
	// 		},
	// 	},
	// }