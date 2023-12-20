import constants from "../constants";
import notification from "../notification";
import startUpdatingGamepadState from "../input";
import log from "../log";

export class WebRTC {
  constructor(display) {
    this.conn = null;
    this.remoteCandidates = [];
    this.localCandidates = [];
    this.dataChannel = null;
    this.mediaStream = null;

    this.display = display;
    this.onData = null;
  }
  init(offer) {
    this.conn = new RTCPeerConnection({
      iceServers: [{ urls: "stun:stun.l.google.com:19302" }],
    });
    this.conn.ontrack = (e) => {
      this.mediaStream = document.createElement(e.track.kind);
      this.mediaStream.srcObject = e.streams[0];
      this.mediaStream.autoplay = true;
      this.mediaStream.playsinline = true;
      this.mediaStream.controls = false;
      this.display.appendChild(this.mediaStream);
      this.mediaStream.play();
    };
    this.conn.addTransceiver("video", { direction: "recvonly" });
    this.conn.addTransceiver("audio", { direction: "recvonly" });

    this.conn.ondatachannel = (e) => {
      log.debug(`[rtc] ondatachannel, ${e.channel.label}`);
      this.dataChannel = e.channel;
      this.dataChannel.onopen = () => {
        log.debug("[rtc] the input channel has been opened");
        notification.broadcast(
          constants.RTC_EVENT.RTC_INPUT_READY,
          this.dataChannel.send.bind(this.dataChannel)
        );
      };
      this.dataChannel.onerror = (e) => {
        log.error("[rtc] the input channel has been closed", e);
      };
      this.dataChannel.onclose = () => {
        log.info("[rtc] the input channel has been closed");
      };
    };

    this.conn.oniceconnectionstatechange = (e) => {
      log.info(`[rtc] State: ${this.conn.iceConnectionState}`);
      switch (this.conn.iceConnectionState) {
        case "connected":
          notification.broadcast(constants.RTC_EVENT.RTC_CONNECTION_READY,constants.RTC_SIGNAL.CONNECTION_READY);
          break;
        case "disconnected":
        case "failed":
          this.stop();
          break;
      }
    };
    this.conn.onicecandidate = (e) => {
      if (e.candidate) {
        this.localCandidates.push(e.candidate);
      }
    };
    this.conn.onicegatheringstatechange = (e) => {
      log.info(`[rtc] Gathering state: ${this.conn.iceGatheringState}`);
      switch (this.conn.iceGatheringState) {
        case "complete":
          this.#sendCandidtes();
          break;
      }
    };
    this.#setSDP(offer);
  }
  async #setSDP(offer) {
    await this.conn
      .setRemoteDescription(offer)
      .then(() => log.info("[rtc] Remote description has been seted"))
      .catch((e) => {
        log.error("[rtc] setRemoteDescription", e);
      });
    const answer = await this.conn.createAnswer();
    await this.conn.setLocalDescription(answer).then(() => {
      log.info("[rtc] Local description has been seted");
      notification.broadcast(constants.RTC_EVENT.RTC_SDP_ANSWER_CREATED,answer);
    });
  }
  addice() {
    this.remoteCandidates.forEach((candidate) => {this.conn.addIceCandidate(candidate)});
  }

  #sendCandidtes = () => {
    this.localCandidates.forEach((candidate) => {
      notification.broadcast(constants.RTC_EVENT.RTC_ICE_CANDIDATE_FOUND,candidate);
    });
  };
  stop() {
    if (this.mediaStream) {
      this.display = null;
      this.mediaStream.pause();
      this.mediaStream.remove();
      this.mediaStream.srcObject = null;
      this.mediaStream = null;
    }
    if (this.conn) {
      this.conn.close();
      this.conn = null;
    }
    if (this.dataChannel) {
      this.dataChannel.close();
      this.dataChannel = null;
      this.onData = null;
    }
    this.localCandidates = null;
    this.remoteCandidates = null;
    log.info("[rtc] connection has been closed");
  }

  input(data) {
    this.dataChannel.send(startUpdatingGamepadState());
  }

  onData(fn) {
    this.onData = fn;
  }
}

export default WebRTC;
