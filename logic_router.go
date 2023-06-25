package fuse

import (
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/phuhao00/spoor/logger"
	"math"
	"strconv"
	"sync"
)

type LogicRouter struct {
	littleEndian       bool
	forwardMsgHandlers *sync.Map
	msgHandlers        map[uint64]MsgHandler
	sceneMsgHandlers   map[uint64]SceneMsgHandler
	serverMsgHandlers  map[uint64]ServerMsgHandler
	defaultHandler     MsgHandler
	preHandler         MsgHandler
	postHandler        MsgHandler
	scenePreHandler    SceneMsgHandler
	scenePostHandler   SceneMsgHandler
	serverPreHandler   ServerMsgHandler
	serverPostHandler  ServerMsgHandler
	//forwardMsgHandlers map[uint16]ForwardMessageHandler
}

type MsgHandler func(msgID uint64, data []byte)

type ForwardMessageHandler func(srvID string, userID uint64, msgID uint64, data []byte)

type SceneMsgHandler func(msgID uint64, data []byte)

type ServerMsgHandler func(srvID string, msgID uint64, data []byte)

func NewLogicRouter() *LogicRouter {
	r := new(LogicRouter)
	r.littleEndian = true
	r.msgHandlers = make(map[uint64]MsgHandler)
	//r.forwardMsgHandlers = make(map[uint16]ForwardMessageHandler)
	r.forwardMsgHandlers = &sync.Map{}
	r.sceneMsgHandlers = make(map[uint64]SceneMsgHandler)
	r.serverMsgHandlers = make(map[uint64]ServerMsgHandler)
	return r
}

func (r *LogicRouter) Clear() {
	r.msgHandlers = nil
	r.forwardMsgHandlers = nil
	r.sceneMsgHandlers = nil
	r.serverMsgHandlers = nil
}

func (r *LogicRouter) SetDefaultHandler(handler MsgHandler) {
	r.defaultHandler = handler
}

func (r *LogicRouter) SetPreHandler(handler MsgHandler) {
	r.preHandler = handler
}

func (r *LogicRouter) SetPostHandler(handler MsgHandler) {
	r.postHandler = handler
}

func (r *LogicRouter) SetScenePreHandler(handler SceneMsgHandler) {
	r.scenePreHandler = handler
}

func (r *LogicRouter) SetScenePostHandler(handler SceneMsgHandler) {
	r.scenePostHandler = handler
}

func (r *LogicRouter) SetServerPreHandler(handler ServerMsgHandler) {
	r.serverPreHandler = handler
}

func (r *LogicRouter) SetServerPostHandler(handler ServerMsgHandler) {
	r.serverPostHandler = handler
}

func (r *LogicRouter) SetByteOrder(littleEndian bool) {
	r.littleEndian = littleEndian
}

func (r *LogicRouter) Register(msgID uint64, msgHandler MsgHandler) bool {
	if msgID >= math.MaxUint16 {
		logger.Error("too many protobuf messages (max = %v)", math.MaxUint16)
		return false
	}

	handler, ok := r.msgHandlers[msgID]
	if ok {
		logger.Warn("message %v is already registered handler:%v", msgID, handler)
	}
	r.msgHandlers[msgID] = msgHandler
	return true
}

func (r *LogicRouter) RegisterForwardHandler(msgID uint16, msgHandler ForwardMessageHandler) bool {
	if msgID >= math.MaxUint16 {
		logger.Error("too many protobuf messages (max = %v)", math.MaxUint16)
		return false
	}

	_, ok := r.forwardMsgHandlers.Load(msgID)
	if ok {
		logger.Warn("message %v is already registered", msgID)
	}
	r.forwardMsgHandlers.Store(msgID, msgHandler)
	return true
}

func (r *LogicRouter) RegSceneMsgHandler(msgID uint64, msgHandler SceneMsgHandler) bool {
	if msgID >= math.MaxUint16 {
		logger.Error("too many protobuf messages (max = %v)", math.MaxUint16)
		return false
	}

	handler, ok := r.sceneMsgHandlers[msgID]
	if ok {
		logger.Warn("message %v is already registered handler:%v", msgID, handler)
	}
	r.sceneMsgHandlers[msgID] = msgHandler
	return true
}

func (r *LogicRouter) RegServerMsgHandler(msgID uint64, msgHandler ServerMsgHandler) bool {
	if msgID >= math.MaxUint16 {
		logger.Error("too many protobuf messages (max = %v)", math.MaxUint16)
		return false
	}

	handler, ok := r.serverMsgHandlers[msgID]
	if ok {
		logger.Warn("message %v is already registered handler:%v", msgID, handler)
	}
	r.serverMsgHandlers[msgID] = msgHandler
	return true
}

func (r *LogicRouter) Route(data []byte) (uint64, error) {

	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}

	var msgID uint64
	if r.littleEndian {
		msgID = binary.LittleEndian.Uint64(data)
	} else {
		msgID = binary.BigEndian.Uint64(data)
	}

	handler, ok := r.msgHandlers[msgID]

	if ok && handler != nil {
		if r.preHandler != nil {
			r.preHandler(msgID, data)
		}
		handler(msgID, data)
		if r.postHandler != nil {
			r.postHandler(msgID, data)
		}
	} else {
		if r.defaultHandler != nil {
			r.defaultHandler(msgID, data)
			return msgID, nil
		}
		return msgID, errors.New("unknown msg msgID:" + strconv.Itoa(int(msgID)))
	}

	return msgID, nil
}

func (r *LogicRouter) ForwardRoute(srvID string, userID uint64, data []byte) (uint64, error) {

	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}

	var msgID uint64
	if r.littleEndian {
		msgID = binary.LittleEndian.Uint64(data)
	} else {
		msgID = binary.BigEndian.Uint64(data)
	}

	handler, ok := r.forwardMsgHandlers.Load(msgID)

	if ok && handler != nil {
		msgHandler, ok := handler.(ForwardMessageHandler)
		if ok {
			msgHandler(srvID, userID, msgID, data)
		}
	} else {
		if r.defaultHandler != nil {
			r.defaultHandler(msgID, data)
			return msgID, nil
		}
		return msgID, errors.New("unknown msg msgID:" + strconv.Itoa(int(msgID)))
	}

	return msgID, nil
}

func (r *LogicRouter) SceneRoute(data []byte) (uint64, error) {

	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}

	var msgID uint64
	if r.littleEndian {
		msgID = binary.LittleEndian.Uint64(data)
	} else {
		msgID = binary.BigEndian.Uint64(data)
	}

	handler, ok := r.sceneMsgHandlers[msgID]

	if ok && handler != nil {
		if r.scenePreHandler != nil {
			r.scenePreHandler(msgID, data)
		}
		handler(msgID, data)
		if r.scenePostHandler != nil {
			r.scenePostHandler(msgID, data)
		}
	} else {
		return msgID, errors.New("unknown msg msgID:" + strconv.Itoa(int(msgID)))
	}

	return msgID, nil
}

func (r *LogicRouter) ServerRoute(srvID string, data []byte) (uint64, error) {

	if len(data) < 2 {
		return 0, errors.New("protobuf data too short")
	}

	var msgID uint64
	if r.littleEndian {
		msgID = binary.LittleEndian.Uint64(data)
	} else {
		msgID = binary.BigEndian.Uint64(data)
	}

	handler, ok := r.serverMsgHandlers[msgID]

	if ok && handler != nil {
		if r.serverPreHandler != nil {
			r.serverPreHandler(srvID, msgID, data)
		}
		handler(srvID, msgID, data)
		if r.serverPostHandler != nil {
			r.serverPostHandler(srvID, msgID, data)
		}
	} else {
		return msgID, errors.New("unknown msg msgID:" + strconv.Itoa(int(msgID)))
	}

	return msgID, nil
}

func (r *LogicRouter) GetRegisterMsgids() []uint32 {
	msgIds := make([]uint32, 0, 64)
	for msgID := range r.msgHandlers {
		msgIds = append(msgIds, uint32(msgID))
	}
	return msgIds
}

func (r *LogicRouter) Unmarshal(data []byte, msg interface{}) error {
	if len(data) < 2 {
		return errors.New("protobuf data too short")
	}
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return fmt.Errorf("msg is not protobuf message")
	}
	return proto.UnmarshalMerge(data[2:], pbMsg)
}

func (r *LogicRouter) Marshal(msgID uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	// data
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	// 4byte = len(flag)[2byte] + len(msgID)[2byte]
	buf := make([]byte, 4+len(data))
	if r.littleEndian {
		binary.LittleEndian.PutUint16(buf[0:2], 0)
		binary.LittleEndian.PutUint16(buf[2:], msgID)
	} else {
		binary.BigEndian.PutUint16(buf[0:2], 0)
		binary.BigEndian.PutUint16(buf[2:], msgID)
	}
	copy(buf[4:], data)
	return buf, err
}

func (r *LogicRouter) InnerServerMarshal(msgID uint16, msg interface{}) ([]byte, error) {
	pbMsg, ok := msg.(proto.Message)
	if !ok {
		return []byte{}, fmt.Errorf("msg is not protobuf message")
	}
	data, err := proto.Marshal(pbMsg)
	if err != nil {
		return data, err
	}
	buf := make([]byte, 2+len(data))
	if r.littleEndian {
		binary.LittleEndian.PutUint16(buf[0:2], msgID)
	} else {
		binary.BigEndian.PutUint16(buf[0:2], msgID)
	}
	copy(buf[2:], data)
	return buf, err
}
