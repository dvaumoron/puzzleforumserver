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

	dbclient "github.com/dvaumoron/puzzledbclient"
	"github.com/dvaumoron/puzzleforumserver/model"
	pb "github.com/dvaumoron/puzzleforumservice"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const dbAccessMsg = "Failed to access database"

var errInternal = errors.New("internal service error")

// server is used to implement puzzleforumservice.ForumServer.
type server struct {
	pb.UnimplementedForumServer
	db     *gorm.DB
	logger *zap.Logger
}

func New(db *gorm.DB, logger *zap.Logger) pb.ForumServer {
	db.AutoMigrate(&model.Thread{}, &model.Message{})
	return server{db: db, logger: logger}
}

func (s server) CreateThread(ctx context.Context, request *pb.CreateRequest) (*pb.Response, error) {
	userId := request.UserId
	thread := model.Thread{
		ObjectId: request.ContainerId, UserId: userId, Title: request.Title,
	}
	if text := request.Text; text != "" {
		thread.Messages = []model.Message{{UserId: userId, Text: text}}
	}
	if err := s.db.Create(&thread).Error; err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Response{Success: true, Id: thread.ID}, nil
}

func (s server) CreateMessage(ctx context.Context, request *pb.CreateRequest) (*pb.Response, error) {
	message := model.Message{ThreadID: request.ContainerId, UserId: request.UserId, Text: request.Text}
	if err := s.db.Create(&message).Error; err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Response{Success: true, Id: message.ID}, nil
}

func (s server) GetThread(ctx context.Context, request *pb.IdRequest) (*pb.Content, error) {
	var thread model.Thread
	if err := s.db.First(&thread, request.Id).Error; err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return convertThreadFromModel(thread), nil
}

func (s server) GetThreads(ctx context.Context, request *pb.SearchRequest) (*pb.Contents, error) {
	objectId := request.ContainerId
	filter := request.Filter
	noFilter := filter == ""

	threadRequest := s.db.Model(&model.Thread{})
	if noFilter {
		threadRequest.Where("object_id = ?", objectId)
	} else {
		filter = dbclient.BuildLikeFilter(filter)
		threadRequest.Where("object_id = ? AND title LIKE ?", objectId, filter)
	}

	var total int64
	err := threadRequest.Count(&total).Error
	if err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	if total == 0 {
		return &pb.Contents{}, nil
	}

	var threads []model.Thread
	page := dbclient.Paginate(s.db, request.Start, request.End).Order("created_at desc")
	if noFilter {
		err = page.Find(&threads, "object_id = ?", objectId).Error
	} else {
		err = page.Find(&threads, "object_id = ? AND title LIKE ?", objectId, filter).Error
	}

	if err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Contents{List: convertThreadsFromModel(threads), Total: uint64(total)}, nil
}

func (s server) GetMessages(ctx context.Context, request *pb.SearchRequest) (*pb.Contents, error) {
	threadId := request.ContainerId
	filter := request.Filter

	messageRequest := s.db.Model(&model.Message{})
	if filter == "" {
		messageRequest.Where("thread_id = ?", threadId)
	} else {
		filter = dbclient.BuildLikeFilter(filter)
		messageRequest.Where("thread_id = ? AND text LIKE ?", threadId, filter)
	}

	var total int64
	err := messageRequest.Count(&total).Error
	if err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	if total == 0 {
		return &pb.Contents{}, nil
	}

	var messages []model.Message
	page := dbclient.Paginate(s.db, request.Start, request.End).Order("created_at asc")
	if filter == "" {
		err = page.Find(&messages, "thread_id = ?", threadId).Error
	} else {
		err = page.Find(&messages, "thread_id = ? AND text LIKE ?", threadId, filter).Error
	}

	if err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Contents{List: convertMessagesFromModel(messages), Total: uint64(total)}, nil
}

func (s server) DeleteThread(ctx context.Context, request *pb.IdRequest) (*pb.Response, error) {
	if err := s.db.Delete(&model.Thread{}, request.Id).Error; err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Response{Success: true}, nil
}

func (s server) DeleteMessage(ctx context.Context, request *pb.IdRequest) (*pb.Response, error) {
	if err := s.db.Delete(&model.Message{}, request.Id).Error; err != nil {
		s.logger.Error(dbAccessMsg, zap.Error(err))
		return nil, errInternal
	}
	return &pb.Response{Success: true}, nil
}

func convertThreadFromModel(thread model.Thread) *pb.Content {
	return &pb.Content{
		Id: thread.ID, CreatedAt: thread.CreatedAt.Unix(), UserId: thread.UserId, Text: thread.Title,
	}
}

func convertThreadsFromModel(threads []model.Thread) []*pb.Content {
	resThreads := make([]*pb.Content, 0, len(threads))
	for _, thread := range threads {
		resThreads = append(resThreads, convertThreadFromModel(thread))
	}
	return resThreads
}

func convertMessagesFromModel(messages []model.Message) []*pb.Content {
	resMessages := make([]*pb.Content, 0, len(messages))
	for _, message := range messages {
		resMessages = append(resMessages, &pb.Content{
			Id: message.ID, CreatedAt: message.CreatedAt.Unix(), UserId: message.UserId, Text: message.Text,
		})
	}
	return resMessages
}
