# go-rtmp

## go-rtmp란?
RTMP(Real Time Messaging Protocol)을 golang으로 구현한 패키지입니다.

### RTMP(Real Time Messaging Protocol)
Adobe가 멀티미디어 전송 스트림을 위해 전송 프로토콜을 multiplexing 및 패킷화하기 위해 설계한 응용프로그램 수준의 프로토콜입니다. 

## Handshake
RTMP는 handshake부터 시작합니다. 동적 크기의 chunk를 가지고 있는 다른 프로토콜과는 다르게 3개의 정적 크기의 chunk로 구성되있습니다.
클라이언트와 서버는 3개의 같은 chunk를 보냅니다. 클라이언트에서 보내진 chunk는 C0, C1, C2이고, 서버에서 보내진 chunk는 S0, S1, S2 입니다.

### C0 와 S0 포멧
8bit 정수로 구성되있습니다.

Version을 들 담고 있는데 C0에서는 클라이언트에서 요청하는 `RTMP의 버전`이다. S0에서는 서버에서 선택된 `RTMP의 버전`입니다

- `version`은 3으로 정해져 있습니다. 0~2는 이전에 사용되어 더 이상 사용되지 않습니다. 4~31은 미래의 버전을 위해 남겨져 있습니다. 그리고 나머지 32~255는 허용되지 않습니다.(RTMP를 다른 text 기반의 프로토콜과 구분하기 위해 허용하고 있습니다.) 
- 만약 서버가 클라이언트가 요청한 서버를 인식하지 못했다면 **반드시** 3을 반환해야 합니다.
- 클라이언트는 3으로 바꾸던가 handshake를 중지해야 합니다.

### Handshake 다이어그램
![img](https://user-images.githubusercontent.com/38347891/213145826-a6662bc8-9127-475c-883a-d7b12007805e.png)
## 참고
- https://rtmp.veriskope.com/docs/spec/#71rtmp-message-types
- https://github.com/yutopp/go-rtmp
