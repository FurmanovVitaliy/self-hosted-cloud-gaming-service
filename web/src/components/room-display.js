import constants from '../common/constants';
import WebRtc  from '../common/network/webrtc';
import event from '../common/event';
import log from "../common/log";


class RoomDisplay extends HTMLElement {
  constructor(){
    super();
    const shadow = this.attachShadow({mode: 'open'});

    const display = document.createElement('div');
    display.setAttribute('class', 'main-display');
    const stream = document.createElement('div');
    stream.setAttribute('id', 'remote-video');



    const style = document.createElement('style');
    style.textContent = `
    `;

    shadow.appendChild(style);
    shadow.appendChild(display);
    shadow.appendChild(stream);
  }

  connectedCallback(){
    this.upgardeConnection();

  }
  async upgardeConnection(){
   
    const shadow = this.shadowRoot;
    const url = window.location.pathname;
    const uuid = url.split('/')[2];

    const remoteVideo = shadow.querySelector('#remote-video');

    const webrtc = new WebRtc(remoteVideo);
    const socket = new WebSocket(`ws://localhost:8000/room/join/${uuid}?username=lock`);

    window.addEventListener('beforeunload', () => {webrtc.stop(), socket.close()});
    window.addEventListener('unload', () => {webrtc.stop(), socket.close()});

    socket.onopen = () => {
      log.info ('[ws] Socket is open');
        
      event.execute(constants.RTC_SDP_ANSWER_CREATED, (answer => {
      socket.send(JSON.stringify(answer))
      log.info('[ws] -> Answer was send');
      }))

      event.execute(constants.RTC_ICE_CANDIDATE_FOUND, (candidate=> {
      socket.send(JSON.stringify(candidate));
      log.info('[ws] -> Candidate was send');
      }));
    }
    
    socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      switch (true) {
        case data.hasOwnProperty('sdp'):
          log.info('[ws] <- Offer received');
           webrtc.start(data);
           break;
        
        case data.hasOwnProperty('candidate'):
          log.info('[ws] <- Candidate received');
          webrtc.setCandidate(data);
          break;
      
        default:
          log.warn('[ws] <- Unknown message received');
          //log.warn(JSON.stringify(e.data))
          break;
      }
    }
    socket.onclose = () => {log.info ('[ws] Socket is closed');      }
    socket.onerror = (error) => {log.error (`[ws] Socket error: ${error}`);}
  }
}
customElements.define('room-display', RoomDisplay);
