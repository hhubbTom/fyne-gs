// package signaling

// import (
// 	"encoding/json"
// 	"log/slog"
// 	"net/url"

// 	"github.com/gorilla/websocket"
// 	"github.com/pion/webrtc/v4"
// )

// type GamepadStateDto struct {
// 	ButtonNorth bool    `json:"buttonNorth"`
// 	ButtonSouth bool    `json:"buttonSouth"`
// 	ButtonWest  bool    `json:"buttonWest"`
// 	ButtonEast  bool    `json:"buttonEast"`
// 	DpadUp      bool    `json:"dpadUp"`
// 	DpadDown    bool    `json:"dpadDown"`
// 	DpadLeft    bool    `json:"dpadLeft"`
// 	DpadRight   bool    `json:"dpadRight"`
// 	LeftStickX  float64 `json:"leftStickX"`
// 	LeftStickY  float64 `json:"leftStickY"`
// 	// 根据需要添加更多控制器状态
// }

// 需要做:  建立websockt连接
// 发送游戏配置
// SDP协商
// ICE交换候选
// 数据通道

// type SignalClient struct {
// 	conn             *websocket.Conn
// 	peerConnection   *webrtc.PeerConnection
// 	gameConfig       *config.GameConfig
// 	codecConfig      *config.CodecConfig
// 	onSDPOffer       func(webrtc.SessionDescription)
// 	onICECandidate   func(webrtc.ICECandidateInit)
// 	onControllerData func(GamepadStateDto)
// 	dataChannel      *webrtc.DataChannel
// }

// func NewSignalClient(gameConfig *config.GameConfig, codecConfig *config.CodecConfig) *SignalClient {
// 	return &SignalClient{
// 		gameConfig:  gameConfig,
// 		codecConfig: codecConfig,
// 	}
// }

// func (s *SignalClient) Connect(serverURL string) error {
// 	u, err := url.Parse(serverURL)
// 	if err != nil {
// 		return err
// 	}

// 	// 将 http(s):// 转换为 ws(s)://
// 	wsURL := "ws://" + u.Host + "/webrtc"
// 	if u.Scheme == "https" {
// 		wsURL = "wss://" + u.Host + "/webrtc"
// 	}

// 	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 	if err != nil {
// 		return err
// 	}
// 	s.conn = conn

// 	// 连接成功后发送初始配置
// 	configMsg := map[string]interface{}{
// 		"game_config":  s.gameConfig,
// 		"codec_config": s.codecConfig,
// 	}
// 	jsonMsg, err := json.Marshal(configMsg)
// 	if err != nil {
// 		return err
// 	}

// 	if err := s.conn.WriteMessage(websocket.TextMessage, jsonMsg); err != nil {
// 		return err
// 	}
// 	slog.Info("Sent initial game and codec configuration")

// 	// 启动消息处理
// 	go s.readMessages()

// 	return nil
// }

// func (s *SignalClient) readMessages() {
// 	for {
// 		_, msg, err := s.conn.ReadMessage()
// 		if err != nil {
// 			slog.Error("WebSocket read error", "error", err)
// 			return
// 		}

// 		// 尝试解析为 SDP
// 		var sdpMsg webrtc.SessionDescription
// 		if err := json.Unmarshal(msg, &sdpMsg); err == nil && sdpMsg.Type != "" {
// 			slog.Info("Received SDP message", "type", sdpMsg.Type)
// 			if s.onSDPOffer != nil {
// 				s.onSDPOffer(sdpMsg)
// 			}
// 			continue
// 		}

// 		// 尝试解析为 ICE candidate
// 		var candidateMsg webrtc.ICECandidateInit
// 		if err := json.Unmarshal(msg, &candidateMsg); err == nil && candidateMsg.Candidate != "" {
// 			slog.Info("Received ICE candidate")
// 			if s.onICECandidate != nil {
// 				s.onICECandidate(candidateMsg)
// 			}
// 			continue
// 		}

// 		slog.Warn("Received unknown message type", "message", string(msg))
// 	}
// }

// func (s *SignalClient) SendSDPAnswer(answer webrtc.SessionDescription) error {
// 	jsonMsg, err := json.Marshal(answer)
// 	if err != nil {
// 		return err
// 	}

// 	return s.conn.WriteMessage(websocket.TextMessage, jsonMsg)
// }

// func (s *SignalClient) SendICECandidate(candidate webrtc.ICECandidateInit) error {
// 	jsonMsg, err := json.Marshal(candidate)
// 	if err != nil {
// 		return err
// 	}

// 	return s.conn.WriteMessage(websocket.TextMessage, jsonMsg)
// }

// func (s *SignalClient) SetupPeerConnection(iceServers []webrtc.ICEServer) (*webrtc.PeerConnection, error) {
// 	config := webrtc.Configuration{
// 		ICEServers: iceServers,
// 	}

// 	pc, err := webrtc.NewPeerConnection(config)
// 	if err != nil {
// 		return nil, err
// 	}

// 	s.peerConnection = pc

// 	// 设置ICE候选事件处理
// 	pc.OnICECandidate(func(candidate *webrtc.ICECandidate) {
// 		if candidate == nil {
// 			return
// 		}

// 		candidateInit := candidate.ToJSON()
// 		slog.Info("Local ICE candidate found", "candidate", candidateInit.Candidate)

// 		if err := s.SendICECandidate(candidateInit); err != nil {
// 			slog.Error("Failed to send ICE candidate", "error", err)
// 		}
// 	})

// 	// 设置数据通道处理
// 	pc.OnDataChannel(func(d *webrtc.DataChannel) {
// 		slog.Info("New data channel created", "label", d.Label())

// 		if d.Label() == "controller" {
// 			s.dataChannel = d

// 			d.OnMessage(func(msg webrtc.DataChannelMessage) {
// 				var gamepadState GamepadStateDto
// 				if err := json.Unmarshal(msg.Data, &gamepadState); err != nil {
// 					slog.Error("Failed to parse gamepad state", "error", err)
// 					return
// 				}

// 				if s.onControllerData != nil {
// 					s.onControllerData(gamepadState)
// 				}
// 			})
// 		}
// 	})

// 	return pc, nil
// }

// func (s *SignalClient) HandleSDPOffer(offer webrtc.SessionDescription) error {
// 	if s.peerConnection == nil {
// 		return nil
// 	}

// 	// 设置远程SDP描述
// 	if err := s.peerConnection.SetRemoteDescription(offer); err != nil {
// 		return err
// 	}

// 	// 创建应答
// 	answer, err := s.peerConnection.CreateAnswer(nil)
// 	if err != nil {
// 		return err
// 	}

// 	// 设置本地SDP描述
// 	if err := s.peerConnection.SetLocalDescription(answer); err != nil {
// 		return err
// 	}

// 	// 发送SDP应答到服务器
// 	return s.SendSDPAnswer(answer)
// }

// func (s *SignalClient) OnSDPOffer(handler func(webrtc.SessionDescription)) {
// 	s.onSDPOffer = handler
// }

// func (s *SignalClient) OnICECandidate(handler func(webrtc.ICECandidateInit)) {
// 	s.onICECandidate = handler
// }

// func (s *SignalClient) OnControllerData(handler func(GamepadStateDto)) {
// 	s.onControllerData = handler
// }

// func (s *SignalClient) Close() error {
// 	if s.peerConnection != nil {
// 		if err := s.peerConnection.Close(); err != nil {
// 			return err
// 		}
// 	}

// 	if s.conn != nil {
// 		return s.conn.Close()
// 	}

// 	return nil
// }
