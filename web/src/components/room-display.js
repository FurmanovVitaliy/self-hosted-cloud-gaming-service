import { createRoom } from "../api/roomsApi";
import { navigateTo } from "../router"

class RoomDisplay extends HTMLElement {
    constructor(){
        super()
        const shadow = this.attachShadow({mode: 'open'})

        const display = document.createElement('div')
        display.setAttribute('class', 'main-display')
        const stream = document.createElement('div')
        stream.setAttribute('id', 'remote-video')
        const log = document.createElement('div')
        log.setAttribute('id', 'event-log')
        
     
        const style = document.createElement('style')
        style.textContent = `
        `
        
        shadow.appendChild(style)
        shadow.appendChild(display)
        shadow.appendChild(stream)
        shadow.appendChild(log)
    }

    connectedCallback(){
        const shadow = this.shadowRoot;
        this.upgardeConnection()
    }
    
    upgardeConnection(){
        const shadow = this.shadowRoot;
        const url = window.location.pathname;
        const uuid = url.split('/')[2];

        const socket = new WebSocket(`ws://localhost:8000/room/join/${uuid}?username=lock`);
        // Получаем ссылку на элемент div
        const divElement = shadow.querySelector('#event-log');
        const remoteVideo = shadow.querySelector('#remote-video');

// Функция для логирования сообщения
const log = (msg) => {
  divElement.innerHTML += msg + '<br>';
};


// Создаем объект RTCPeerConnection с серверами STUN
const peerConnection = new RTCPeerConnection();

peerConnection.ontrack = function (event) {
  const el = document.createElement(event.track.kind);
  el.srcObject = event.streams[0];
  el.autoplay = true;
  el.controls = true;
  remoteVideo.appendChild(el);
};
// Добавляем треки для видео и аудио
peerConnection.addTransceiver('video', {
  direction: 'sendrecv',
});
peerConnection.addTransceiver('audio', {
  direction: 'sendrecv',
});

peerConnection.oniceconnectionstatechange = e => log(peerConnection.iceConnectionState)
// Устанавливаем обработчик события icecandidate
peerConnection.onicecandidate = async (e) => {
  if (e.candidate) {
    await sendIceCandidate(e.candidate);
  }
};

// Функция для создания WebRTC соединения
async function createWebRTC(offer) {
  try {
    // Устанавливаем удаленное описание
    await peerConnection.setRemoteDescription(offer);
    // Устанавливаем обработчик для приема видео и аудио

    // Создаем и устанавливаем локальное описание (answer)
    const answer = await peerConnection.createAnswer();
    await peerConnection.setLocalDescription(answer);
   

    // Отправляем локальное описание на сервер через WebSocket
    await sendAnswer(answer);
  } catch (error) {
    console.error('Error creating or setting local description:', error);
  }
}

// Обработчик открытия WebSocket соединения
socket.addEventListener('open', () => {
  console.log('WebSocket соединение установлено');
});

// Обработчик приема сообщения по WebSocket
socket.addEventListener('message', async (event) => {
  const JsonData = JSON.parse(event.data);
  console.log('Полученные данные:', JsonData);
  switch (JsonData.type) {
    case 'offer':
      await createWebRTC(JsonData);
      break;
      default:
      await peerConnection.addIceCandidate(JsonData);
      break;
  }
});

// Обработчик закрытия WebSocket соединения
socket.addEventListener('close', () => {
  console.log('WebSocket соединение закрыто');
});

// Функция для отправки answer на сервер через WebSocket
async function sendAnswer(answer) {
  console.log('Отправляем answer на сервер через WebSocket:',answer);
  socket.send(JSON.stringify(answer));
}

// Функция для отправки кандидата ICE на сервер через WebSocket
async function sendIceCandidate(candidate) {
  console.log('Отправляем ice кандидата  на сервер через WebSocket:',candidate);
  socket.send(JSON.stringify(candidate));
}
    }
}
customElements.define('room-display', RoomDisplay )