# go-rtmp

## go-rtmp란?
RTMP(Real Time Messaging Protocol)을 golang으로 구현한 패키지입니다.

### RTMP(Real Time Messaging Protocol)
Adobe가 멀티미디어 전송 스트림을 위해 전송 프로토콜을 multiplexing 및 패킷화하기 위해 설계한 응용프로그램 수준의 프로토콜입니다. 

## Handshake
RTMP는 handshake부터 시작합니다. 동적 크기의 chunk를 가지고 있는 다른 프로토콜과는 다르게 3개의 정적 크기의 chunk로 구성되있습니다.
클라이언트와 서버는 3개의 같은 chunk를 보냅니다. 클라이언트에서 보내진 chunk는 C0, C1, C2이고, 서버에서 보내진 chunk는 S0, S1, S2 입니다.

### Handshake 다이어그램
![img](https://user-images.githubusercontent.com/38347891/213145826-a6662bc8-9127-475c-883a-d7b12007805e.png)
## 참고
- https://rtmp.veriskope.com/docs/spec/#71rtmp-message-types
- https://github.com/yutopp/go-rtmp
