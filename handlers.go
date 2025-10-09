package main

import "fmt"

func addUserHandler(req *Request, res *Response) {
	fmt.Println("This is from addUserHandler")
	res.Body = []byte("Hello from add user")
}

func allUserHandler(req *Request, res *Response) {
	fmt.Println("This is from allUserHandler")
	res.Body = []byte("Hello from all user")
}
