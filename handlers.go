package main

import "fmt"

func addUserHandler(req *Request, res *Response) {
	print(req.Headers)
	fmt.Println("This is from addUserHandler")
	res.Body = []byte("Hello from add user")
}

func allUserHandler(req *Request, res *Response) {
	fmt.Println("This is from allUserHandler")
	res.Body = []byte("Hello from all user")
}

func getUserById(req *Request, res *Response) {
	print(req.Headers)
	print("Path Params =", req.Headers["PathParams"])
	print("Query Params =", req.Headers["QueryParams"])

	res.Body = []byte("Hello from all user")
}

func getUserByNameAndId(req *Request, res *Response) {
	print("from get users by name and id handler")

	print("Path Params =", req.Headers["PathParams"])
	print("Query Params =", req.Headers["QueryParams"])

	res.Body = []byte("Hello from all user")
}
