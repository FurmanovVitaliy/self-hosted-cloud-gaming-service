import constants from "../common/constants";
import WebRtc from "../common/network/webrtc";
import notification from "../common/notification";
import log from "../common/log";
import { startUpdatingGamepadState } from "../common/input";
import { stopUpdatingGamepadState } from "../common/input";
import {collectDeviceInfo} from "../common/devinfo";
import {msg} from "../common/network/msg";

class RoomDisplay extends HTMLElement {
  constructor() {
    super();
    const shadow = this.attachShadow({ mode: "open" });

    const display = document.createElement("div");
    display.setAttribute("class", "main-display");
    display.setAttribute("autoplay", true);
    display.setAttribute("playsinline", true);
   /* if (display.requestFullscreen) {
      display.requestFullscreen();
    } else if (display.webkitRequestFullscreen) {
      display.webkitRequestFullscreen();
    } else if (display.mozRequestFullScreen) {
      display.mozRequestFullScreen();
    } else if (display.msRequestFullscreen) {
      display.msRequestFullscreen();
    }

*/
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
    });

    this.upgardeConnection();
    notification.execute(
      constants.RTC_EVENT.RTC_INPUT_READY,
      startUpdatingGamepadState,
      true
    );
  }
  upgardeConnection() {
    const url = window.location.pathname;
    const uuid = url.split("/")[2];

    this.socket = new WebSocket(
      `wss://${
        constants.SERVER_IP
      }/room/join/${uuid}?username=${localStorage.getItem("username")}`
    );
    

    this.socket.onopen = () => {
      log.info("[ws] Socket is open");
      this.socket.send(msg("deviceInfo",collectDeviceInfo()));
      log.info("[ws]->Device info sent");
      notification.execute(constants.RTC_EVENT.RTC_CONNECTION_READY,
        (signal) => {
        this.socket.send(msg(constants.RTC_ENTITY_NAME.SIGNAL, signal));
        log.info(`[ws]->WS closure signal was sent`);
      },
      true
      );
    };

    notification.execute(
      constants.RTC_EVENT.RTC_SDP_ANSWER_CREATED,
      (answer) => {
        this.socket.send(msg(constants.RTC_ENTITY_NAME.ANSWER, answer));
        log.info(`[ws]->Answer was sent`);
      },
      true
    );

    notification.execute(
      constants.RTC_EVENT.RTC_ICE_CANDIDATE_FOUND,
      (candidate) => {
        this.socket.send(
          msg(constants.RTC_ENTITY_NAME.CLIENT_CANDIDATE, candidate)
        );
        log.info(`[ws]->Candidate was send`);
      }
    );

    this.socket.onmessage = (e) => {
      const data = JSON.parse(e.data);
      switch (data.content_type) {
        case constants.RTC_ENTITY_NAME.OFFER:
          log.info("[ws] <- Offer received");
          this.webrtc.init(JSON.parse(atob(data.content)));
          break;
        case constants.RTC_ENTITY_NAME.SERVER_CANDIDATE:
          log.debug(`[ws] <- Candidate received : ${atob(data.content)}`);
          this.webrtc.remoteCandidates.push(JSON.parse(atob(data.content)));
          break;
        case constants.RTC_ENTITY_NAME.SIGNAL:
          if (data.content === constants.RTC_SIGNAL.SERVER_ICE_GATHERING_COMPLETE) {
            log.info("[ws] <- Server ICE ready signal received");
            this.webrtc.addice();
          }
          break;
        default:
          log.warn("[ws] <- Unknown message received");
          log.warn(JSON.stringify(data));
          
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
