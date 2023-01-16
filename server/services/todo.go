package services

import (
	"context"
	"log"
	pbTodo "svc-todo/pb/todo"
	pbUser "svc-todo/pb/user"
	"svc-todo/server/db"
)

type TodoService struct {
	pbTodo.UnimplementedTodoServiceServer
}

func (p *TodoService)	List(context.Context, *pbUser.Empty) (*pbTodo.TodoResponse, error) {
	log.Println("List Todo Was Invoked")

	var todos []*pbTodo.Todo

	con := db.CreateCon()	

	stmt := "SELECT t.title, t.status, t.done_at, u.name, u.uuid, t.uuid  from todo t left join `user` u on u.id = t.user_id"

	rows, err := con.Query(stmt)
	if err != nil{
		return nil,	err
	}
	defer rows.Close()
	for rows.Next(){
		var todo pbTodo.Todo
		var user pbUser.User
		err = rows.Scan(&todo.Title, &todo.Status, &todo.DoneAt, &user.Name, &user.Uuid, &todo.Uuid)
		if err != nil {
			return nil, err
		}
		todo.User = &user 
		todos = append(todos, &todo)
	}

	response := &pbTodo.TodoResponse{
		Status: true,
		Data: todos,
	}

	return response, nil
}
func (p *TodoService)	AddTodo(ctx context.Context, todo *pbTodo.AddTodoRequest) (*pbTodo.AddTodoResponse, error){
	log.Println("Add Todo Was Invoked")

	con := db.CreateCon()

	sqlStatement := "INSERT todo (title, status, done_at) VALUES (?, ?, ?)"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(todo.GetTitle(), todo.GetStatus(), todo.GetDoneAt())
	if err != nil {
		return nil, err
	}

	response := &pbTodo.AddTodoResponse{
		Status: true,
		Message: "Success",
	}

	return response, nil
}

func (p *TodoService) Assign(ctx context.Context, todo *pbTodo.AssignRequest) (*pbTodo.AddTodoResponse, error){
	log.Println("Assign Todo to user Was Invoked")
	
	user_id, err := FindUserData(todo.GetUserUuid())
	if err != nil {
		return nil, err
	}

	con := db.CreateCon()

	sqlStatement := "UPDATE todo SET user_id = ? WHERE uuid = ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec(user_id, todo.GetTodoUuid())
	if err != nil {
		return nil, err
	}
	response := &pbTodo.AddTodoResponse{
		Status: true,
		Message: "Success",
	}

	return response, nil
}

func (p *TodoService) MarkAsDone(ctx context.Context, todo *pbTodo.MarkAsDoneRequest) (*pbTodo.AddTodoResponse, error){
	log.Println("Mark As Done Todo Was Invoked")
	
	con := db.CreateCon()

	sqlStatement := "UPDATE todo SET status = ? WHERE uuid = ?"

	stmt, err := con.Prepare(sqlStatement)
	if err != nil {
		return nil, err
	}

	_, err = stmt.Exec("Done", todo.GetTodoUuid())
	if err != nil {
		return nil, err
	}
	response := &pbTodo.AddTodoResponse{
		Status: true,
		Message: "Success",
	}

	return response, nil
}


