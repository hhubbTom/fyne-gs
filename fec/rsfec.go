// package fec
// import (
//     "github.com/twmb/reedsolomon"
//     "github.com/pion/rtcp"
// 	"github.com/pion/webrtc/v3"
// 	"github.com/pion/webrtc/v3/pkg/media"

// )

// // fec/rsfec.go
// type RSFECEncoder struct {
//     enc reedsolomon.Encoder
// }

// func (e *RSFECEncoder) Encode(data []byte) [][]byte {
//     shards, _ := e.enc.Split(data)
//     e.enc.Encode(shards)
//     return shards
// }

// // 在 WebRTC 拦截器中应用
// // interceptor/fec_interceptor.go
// func (f *FECInterceptor) BindRTCPWriter(writer interceptor.RTCPWriter) interceptor.RTCPWriter {
//     return interceptor.RTCPWriterFunc(func(pkts []rtcp.Packet, attributes interceptor.Attributes) (int, error) {
//         // 对数据包应用 FEC
//         fecData := f.encoder.Encode(serializePackets(pkts))
//         return writer.Write(fecData, attributes)
//     })
// }