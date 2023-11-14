import constants from "../common/constants";
import WebRtc from "../common/network/webrtc";
import notification from "../common/notification";
import log from "../common/log";

class RoomDisplay extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });

    const display = document.createElement("div");
    display.setAttribute("class", "main-display");
    this.webrtc = new WebRtc(display);
    this.socket = null;

    const style = document.createElement("style");
    style.textContent = `
    `;

    shadow.appendChild(style);
    shadow.appendChild(display);
  }

  connectedCallback() {
    this.upgardeConnection();
    notification.execute(constants.event.RTC_INPUT_READY,()=>{
      console.log(this.webrtc.dataChannel.send("test"))
    },true);
    
  }


  upgardeConnection() {
    
    const url = window.location.pathname;
    const uuid = url.split("/")[2];
    
    this.socket = new WebSocket(
      `ws://localhost:8000/room/join/${uuid}?username=lock`
      );

      this.socket.onopen = () => {
        log.info("[ws] Socket is open");
        
      };

      notification.execute(constants.event.RTC_SDP_ANSWER_CREATED,(answer) => {
         this.socket.send(JSON.stringify(answer));
         log.info(`[ws]->Answer was send`);
      },true);
      
    notification.execute(constants.event.RTC_ICE_CANDIDATE_FOUND,(candidate) => {
      this.socket.send(JSON.stringify(candidate));
      log.info(`[ws]->Candidate was send`);
    });

    this.socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      switch (true) {
        case data.hasOwnProperty("sdp"):
          log.info("[ws] <- Offer received");
          this.webrtc.init(data)
          break;
        case data.hasOwnProperty("candidate"):
          log.debug(`[ws] <- Candidate received : ${JSON.stringify(data)}`);
         
          this.webrtc.remoteCandidates.push(data);
          break;
        case data.hasOwnProperty("signal"):
          if (data.signal === "Server ICE gathering is complete") {
            log.info("[ws] <- Ready signal received");
            this.webrtc.addice();
          }
          break;
        default:
          log.warn("[ws] <- Unknown message received");
          log.warn(JSON.stringify(e.data));
          break;
      }
    };
    this.socket.onclose = () => {
      log.info("[ws] Socket is closed");
    };
    this.socket.onerror = (error) => {
      log.error(`[ws] Socket error: ${error}`);
    };
  }
}
customElements.define("room-display", RoomDisplay);
