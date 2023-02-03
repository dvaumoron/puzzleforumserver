/*
 *
 * Copyright 2023 puzzleforumserver authors.
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
package forumserver

import (
	"context"
	"errors"

	"github.com/dvaumoron/puzzleforumserver/model"
	pb "github.com/dvaumoron/puzzleforumservice"
	"gorm.io/gorm"
)

const dbAccessMsg = "Failed to access database :"

var errInternal = errors.New("internal service error")

// server is used to implement puzzlerightservice.RightServer.
type server struct {
	pb.UnimplementedForumServer
	db *gorm.DB
}

func New(db *gorm.DB) pb.ForumServer {
	db.AutoMigrate(&model.Thread{}, &model.Message{})
	return server{db: db}
}

func (s server) CreateThread(ctx context.Context, request *pb.CreateRequest) (*pb.Response, error) {
	// TODO
	return nil, nil
}

func (s server) CreateMessage(ctx context.Context, request *pb.CreateRequest) (*pb.Response, error) {
	// TODO
	return nil, nil
}

func (s server) GetThread(ctx context.Context, request *pb.IdRequest) (*pb.Content, error) {
	// TODO
	return nil, nil
}

func (s server) GetThreads(ctx context.Context, request *pb.SearchRequest) (*pb.Contents, error) {
	// TODO
	return nil, nil
}

func (s server) GetMessages(ctx context.Context, request *pb.SearchRequest) (*pb.Contents, error) {
	// TODO
	return nil, nil
}

func (s server) DeleteThread(ctx context.Context, request *pb.IdRequest) (*pb.Response, error) {
	// TODO
	return nil, nil
}

func (s server) DeleteMessage(ctx context.Context, request *pb.IdRequest) (*pb.Response, error) {
	// TODO
	return nil, nil
}
