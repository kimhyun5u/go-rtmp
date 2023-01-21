# go-rtmp

## go-rtmp 란?
RTMP(Real Time Messaging Protocol)을 golang 으로 구현한 패키지입니다.

### RTMP(Real Time Messaging Protocol)
Adobe 가 멀티미디어 전송 스트림을 위해 전송 프로토콜을 multiplexing 및 패킷화하기 위해 설계한 응용프로그램 수준의 프로토콜입니다. 

## Handshake
RTMP 는 handshake 부터 시작합니다. 동적 크기의 chunk 를 가지고 있는 다른 프로토콜과는 다르게 3개의 정적 크기의 chunk 로 구성되있습니다.
클라이언트와 서버는 3개의 같은 chunk 를 보냅니다. 클라이언트에서 보내진 chunk 는 C0, C1, C2이고, 서버에서 보내진 chunk 는 S0, S1, S2 입니다.

### Handshake 시퀸스
handshake 는 클라이언트가 C0, C1 chunk 를 보내면서 시작합니다.

클라이언트는 C2를 보내기 전에 **필히** S1을 기다려야 합니다. 클라이언트는 다른 어떤 데이터를 보내기 전에 **필히** S2를 기다려야 합니다.

서버는 S0, S1를 보내기 전에 C0을 **필히** 기다려야 하고, C1의 도착까지 기다릴 수도 있습니다. 서버는 S2를 보내기전에 C1을 **필히** 기다려야합니다. 서버는 다른 어떤 데이터를 보내기전에 **필히** C2를 기다려야 합니다. 
### C0 와 S0 포멧
8bit 정수로 구성되있습니다.

Version 을 들 담고 있는데 C0에서는 클라이언트에서 요청하는 `RTMP의 버전`이다. S0에서는 서버에서 선택된 `RTMP의 버전`입니다

- `version`은 3으로 정해져 있습니다. 0~2는 이전에 사용되어 더 이상 사용되지 않습니다. 4~31은 미래의 버전을 위해 남겨져 있습니다. 그리고 나머지 32~255는 허용되지 않습니다.(RTMP 를 다른 text 기반의 프로토콜과 구분하기 위해 허용하고 있습니다.) 
- 만약 서버가 클라이언트가 요청한 서버를 인식하지 못했다면 **반드시** 3을 반환해야 합니다.
- 클라이언트는 3으로 바꾸던가 handshake 를 중지해야 합니다.


### C1 과 S1 포멧
총 1536 바이트의 길이이고, TIME, ZERO, RANDOM DATA 로 구성되있습니다.

`TIME(4 bytes)`는  미래에 endpoint 에서 전송될 모든 chunk 를 위한 epoch 로 **반드시** 사용되는 timestamp 를 의미합니다. 이 값은 0 혹은 임의의 수 입니다. 다수의 chunk stream 을 동기화 하기 위해, endpoint 는 다른 chunk stream 의 현재의 값을 보내려고 할 수 있습니다.

`ZERO(4 bytes)`는 **필히** 모두 0 입니다.

`RANDOM DATA(1528 bytes)`는 아무 임의의 값들을 의미합니다. 각각의 endpoint 는 자신에게서 시작한 handshake 에 대한 응답과 peer 에서 시작된 handshake 를 구분해야 하기 때문에 이 값은 **반드시** 충분히 무작위의 값이여야 합니다. 그렇다고 cryptographically-secure randomness 혹은 dynamic values 가 필요한 건 아닙니다.

### C2 와 S2 포멧
총 1536 바이트의 길이이고, 거의 S1과 C1의 메아리라고 볼 수 있고, TIME, TIME2, RANDOM ECHO 로 구성되있습니다.

`TIME(4 bytes)`는 **필히** peer 에서 보낸 (C2를 위한) S1 또는 (S2를 위한) C1에서 보내진 timestamp 를 포함하고 있어야 합니다.

`TIME2(4 bytes)`는 **필히** peer 에서 보낸 이전 패킷(S1 또는 C1)의 읽은 timestamp 를 포함해야합니다.

`RANDOM ECHP(1528 bytes)`는 peer 에서 보낸 (C2를 위한)S1 또는 (C1을 위한)S2에서 보낸 무작위 데이터를  **필히** 포함해야합니다. peer 는 현재의 timestamp 와 함께 time 및 time2를 연결의 대역폭, 대기 시간에 대한 빠른 추정에 사용 할 수 있습니다. 하지만 그다지 유용하지 않습니다.

### Handshake 다이어그램
<img width="651" alt="image" src="https://user-images.githubusercontent.com/38347891/213843599-12321b50-16be-4ab5-9ba3-195d6d48ec47.png">


`초기화되지 않음` 이 단계에서 프로토콜의 버전이 전송됩니다. 클라이언트와 서버 모두 초기화되지 않은 상태입니다. 클라이언트가 C0 패킷에 프로토콜의 버전을 담아 전송합니다. 만약 서버가 그 버전을 지원한다면, S0와 S1를 클라이언트에 보냅니다. 만약 지원하지 않는다면, 적절한 조치를 취합니다. RTMP 에서는, 연결을 끝는 것이 적절합니다.

`Version 전송됨` 클라이언트와 서버 모두 초기 상태에서 버전이 전송되있는 상태로 바뀝니다. 클라이언트는 S1 패킷을 기다리고 서버는 C1 패킷을 기다리고 있습니다. 기다리고 있는 패킷들을 받았을 때, 클라이언트는 C2 패킷을 서버는 S2 패킷을 보냅니다. 그러면 `Ack 전송됨` 상태가 시작됩니다.

`Ack 전송됨` 클라이언트와 서버는 각각 S2와 C2를 기다립니다.

`Handshake 완료` 클라이언트와 서버가 message 를 주고 받습니다.

## 참고
- https://rtmp.veriskope.com/docs/spec/#71rtmp-message-types
- https://github.com/yutopp/go-rtmp
