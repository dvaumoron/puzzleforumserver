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

package main

import (
	_ "embed"

	dbclient "github.com/dvaumoron/puzzledbclient"
	"github.com/dvaumoron/puzzleforumserver/forumserver"
	pb "github.com/dvaumoron/puzzleforumservice"
	grpcserver "github.com/dvaumoron/puzzlegrpcserver"
)

//go:embed version.txt
var version string

func main() {
	s := grpcserver.Make(forumserver.ForumKey, version)
	pb.RegisterForumServer(s, forumserver.New(dbclient.Create(s.Logger), s.Logger))
	s.Start()
}
