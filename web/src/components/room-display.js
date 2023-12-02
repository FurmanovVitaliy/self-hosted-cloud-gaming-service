import constants from "../common/constants";
import WebRtc from "../common/network/webrtc";
import notification from "../common/notification";
import log from "../common/log";
import { startUpdatingGamepadState } from "../common/input";
import { stopUpdatingGamepadState } from "../common/input";

class RoomDisplay extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });

    const display = document.createElement("div");
    display.setAttribute("class", "main-display");
    display.setAttribute("autoplay", true);
    display.setAttribute("playsinline", true);
    if (display.requestFullscreen) {
      display.requestFullscreen();
    } else if (display.webkitRequestFullscreen) {
      display.webkitRequestFullscreen();
    } else if (display.mozRequestFullScreen) {
      display.mozRequestFullScreen();
    } else if (display.msRequestFullscreen) {
      display.msRequestFullscreen();
    }


    this.webrtc = new WebRtc(display);
    this.socket = null;

    const style = document.createElement("style");
    style.textContent = `
    .main-display {
        height: 100vh;
        width: 100vw;
        display: flex;
        justify-content: center;
        align-items: stretch;
        background-color: black;
      }
      .video {
        width: 100vw;
        height: 100vh;
      }

    `;

    shadow.appendChild(style);
    shadow.appendChild(display);
  }

  connectedCallback() {
    window.addEventListener("popstate", (e) => {
      stopUpdatingGamepadState;
      this.webrtc.stop();
      this.socket.close();
      this.socket = null;
      this.webrtc = null;
    });

    this.upgardeConnection();
    notification.execute(
      constants.event.RTC_INPUT_READY,
      startUpdatingGamepadState,
      true
    );
  }
  upgardeConnection() {
    const url = window.location.pathname;
    const uuid = url.split("/")[2];

    this.socket = new WebSocket(
      `wss://${constants.SERVER_IP}/room/join/${uuid}?username=lock`
    );
    //notification.execute(constants.event.RTC_CONNECTION_READY,this.socket.close ,true);

    this.socket.onopen = () => {
      log.info("[ws] Socket is open");
      this.socket.send(JSON.stringify({ signal: "Client is ready" }));
    };

    notification.execute(
      constants.event.RTC_SDP_ANSWER_CREATED,
      (answer) => {
        this.socket.send(JSON.stringify(answer));
        log.info(`[ws]->Answer was send`);
      },
      true
    );

    notification.execute(
      constants.event.RTC_ICE_CANDIDATE_FOUND,
      (candidate) => {
        this.socket.send(JSON.stringify(candidate));
        log.info(`[ws]->Candidate was send`);
      }
    );

    this.socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      switch (true) {
        case data.hasOwnProperty("sdp"):
          log.info("[ws] <- Offer received");
          this.webrtc.init(data);
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
