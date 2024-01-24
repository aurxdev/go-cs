/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	port = flag.Int("port", 50051, "The server port")
	name string
	list string
	db   *sql.DB
) // var

// structure du serveur
type server struct {
	pb.UnimplementedGreeterServer
} // server

// ajoute un membre dans la base de données
func addDatabase(text1 string, text2 string) (int64, error) {
	result, err := db.Exec("INSERT INTO membre (nom, liste) VALUES (?, ?)", text1, text2)
	if err != nil {
		return 0, fmt.Errorf("addDatabase: %v", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("addDatabse: %v", err)
	}
	return id, nil
} // addDatabase

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	name = in.GetName()
	list = in.GetList()
	log.Printf("name: %v", name)
	log.Printf("list: %v", list)
	var msg string
	if name == "" {
		msg = "erreur: nom vide"
	} else {
		msg = "ok"
		addDatabase(name, list)
	}
	return &pb.HelloReply{Message: msg}, nil
} // SayHello

// fonction principale
func main() {
	// -------- connexion à la base de données -------------
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               " ",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "projet",
		AllowNativePasswords: true,
	}

	// ouverture de la connexion
	var err2 error
	db, err2 = sql.Open("mysql", cfg.FormatDSN())
	if err2 != nil {
		log.Fatal(err2)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected!")

	// ----------connexion au serveur----------------
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
} // main
