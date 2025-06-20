// Licensed to the LF AI & Data foundation under one
// or more contributor license agreements. See the NOTICE file
// distributed with this work for additional information
// regarding copyright ownership. The ASF licenses this file
// to you under the Apache License, Version 2.0 (the
// "License"); you may not use this file except in compliance
// with the License. You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package session

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type NodeManagerSuite struct {
	suite.Suite

	nodeManager *NodeManager
}

func (s *NodeManagerSuite) SetupTest() {
	s.nodeManager = NewNodeManager()
}

func (s *NodeManagerSuite) TearDownTest() {
}

func (s *NodeManagerSuite) TestNodeOperation() {
	s.nodeManager.Add(NewNodeInfo(ImmutableNodeInfo{
		NodeID:   1,
		Address:  "localhost",
		Hostname: "localhost",
	}))
	s.nodeManager.Add(NewNodeInfo(ImmutableNodeInfo{
		NodeID:   2,
		Address:  "localhost",
		Hostname: "localhost",
	}))
	s.nodeManager.Add(NewNodeInfo(ImmutableNodeInfo{
		NodeID:   3,
		Address:  "localhost",
		Hostname: "localhost",
	}))

	s.NotNil(s.nodeManager.Get(1))
	s.Len(s.nodeManager.GetAll(), 3)
	s.nodeManager.Remove(1)
	s.Nil(s.nodeManager.Get(1))
	s.Len(s.nodeManager.GetAll(), 2)

	s.nodeManager.Stopping(2)
	s.True(s.nodeManager.IsStoppingNode(2))
	node := s.nodeManager.Get(2)
	node.SetState(NodeStateNormal)
	s.False(s.nodeManager.IsStoppingNode(2))
}

func (s *NodeManagerSuite) TestNodeInfo() {
	node := NewNodeInfo(ImmutableNodeInfo{
		NodeID:   1,
		Address:  "localhost",
		Hostname: "localhost",
	})
	s.Equal(int64(1), node.ID())
	s.Equal("localhost", node.Addr())
	node.setChannelCnt(1)
	node.setSegmentCnt(1)
	s.Equal(1, node.ChannelCnt())
	s.Equal(1, node.SegmentCnt())

	node.UpdateStats(WithSegmentCnt(5))
	node.UpdateStats(WithChannelCnt(5))
	s.Equal(5, node.ChannelCnt())
	s.Equal(5, node.SegmentCnt())

	node.SetLastHeartbeat(time.Now())
	s.NotNil(node.LastHeartbeat())
}

// TestCPUNumFunctionality tests the newly added CPU core number functionality
func (s *NodeManagerSuite) TestCPUNumFunctionality() {
	node := NewNodeInfo(ImmutableNodeInfo{
		NodeID:   1,
		Address:  "localhost:19530",
		Hostname: "test-host",
	})

	// Test initial CPU core number
	s.Equal(int64(0), node.CPUNum())

	// Test WithCPUNum option
	node.UpdateStats(WithCPUNum(8))
	s.Equal(int64(8), node.CPUNum())

	// Test updating CPU core number
	node.UpdateStats(WithCPUNum(16))
	s.Equal(int64(16), node.CPUNum())

	// Test multiple stats update including CPU core number
	node.UpdateStats(
		WithSegmentCnt(100),
		WithChannelCnt(5),
		WithMemCapacity(4096.0),
		WithCPUNum(32),
	)
	s.Equal(int64(32), node.CPUNum())
	s.Equal(100, node.SegmentCnt())
	s.Equal(5, node.ChannelCnt())
	s.Equal(4096.0, node.MemCapacity())
}

// TestMemCapacityFunctionality tests memory capacity related methods
func (s *NodeManagerSuite) TestMemCapacityFunctionality() {
	node := NewNodeInfo(ImmutableNodeInfo{
		NodeID:   1,
		Address:  "localhost:19530",
		Hostname: "test-host",
	})

	// Test initial memory capacity
	s.Equal(float64(0), node.MemCapacity())

	// Test WithMemCapacity option
	node.UpdateStats(WithMemCapacity(1024.5))
	s.Equal(1024.5, node.MemCapacity())

	// Test updating memory capacity
	node.UpdateStats(WithMemCapacity(2048.75))
	s.Equal(2048.75, node.MemCapacity())
}

func TestNodeManagerSuite(t *testing.T) {
	suite.Run(t, new(NodeManagerSuite))
}
