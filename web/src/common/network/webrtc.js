import constants from "../constants";
import event from "../event";
import log from "../log";

export class WebRTC {
    constructor(remoteVideoElement) {

        this.conn = null;
        this.candidates = [];
        this.dataChannel = null;
        this.mediaStream = null;
        this.remoteVideoElement = remoteVideoElement;

        this.onData = null;
        this.#init();
    }

    #init() {
        this.conn = new RTCPeerConnection({ iceServers: [{ urls: "stun:stun.l.google.com:19302" }] });
        
        this.MediaStream = new MediaStream();
        this.conn.ontrack = (event) => {
            this.mediaStream = document.createElement(event.track.kind);
            this.mediaStream.srcObject = event.streams[0];
            this.mediaStream.autoplay = true;
            this.mediaStream.true = false;
            this.remoteVideoElement.appendChild(this.mediaStream);
        };
       // this.conn.addTransceiver('video', { direction: 'recvonly' });
       // this.conn.addTransceiver('audio', { direction: 'recvonly' });

        this.conn.ondatachannel = (e) => {
            log.debug('[rtc] ondatachannel', e.channel.label);
            this.dataChannel = e.channel;
            this.dataChannel.onopen = () => {
                this.inputReady = true;
            };
            this.dataChannel.onerror = (e) => { log.error('[rtc] the input channel has been closed', e); };
            this.dataChannel.onclose = () => { log.info('[rtc] the input channel has been closed'); };
        };

        this.state = this.conn.oniceconnectionstatechange = (e) => {
            log.info(`[rtc] State: ${this.conn.iceConnectionState}`);
            switch (this.conn.iceConnectionState) {
                case 'connected':
                    this.connected = true;
                    event.broadcast(constants.RTC_CONNECTION_READY);
                    break;
                case 'disconnected':
                    this.connected = false;
                    event.broadcast(constants.RTC_CONNECTION_CLOSED);
                    break;
            }
        }
        this.#iceHandler();
    }
    
    async start(offer){

        await  this.conn.setRemoteDescription(offer).then(() => {
            log.debug('[rtc] offer has been set')
        }).catch((error) => {
            log.error(`[rtc] Error setting remote description : ${error}`);
            this.stop();
        });

        await this.conn.createAnswer().then((answer) => {
            this.conn.setLocalDescription(answer);
            event.broadcast(constants.RTC_SDP_ANSWER_CREATED, answer);
        }).catch((error) => {
            this.stop();
            log.error(`[rtc] Error creating answer : ${error}`);
        });
    }


    setCandidate(candidate) {
        this.conn.addIceCandidate(candidate).then(() => {
                log.debug(`[rtc] Candidate added : ${candidate}`);
            })
            .catch((error) => {
                this.stop();
                log.error(`[rtc] Error adding candidate' : ${error}`);
            });
            console.log();
    }

    #iceHandler() {
        this.conn.onicecandidate = async (e) => {
            if(e.candidate) {
                this.candidates.push(e.candidate);
                log.debug(`[rtc] Candidate found : ${e.candidate}`);
            }
            event.broadcast(constants.RTC_ICE_CANDIDATE_FOUND, e.candidate);
        }
    }


    stop() {
        if (this.mediaStream) {
            this.mediaStream.getTracks().forEach(t => {
                t.stop();
                this.mediaStream.removeTrack(t);
            });
            this.mediaStream = null;
        }
        if (this.conn) {
            this.conn.close();
            this.conn = null;
        }
        if (this.dataChannel) {
            this.dataChannel.close();
            this.dataChannel = null;
        }
        this.candidates = [];
        log.info('[rtc] WebRTC has been closed');
    }

    input(data) {

    }

    onData(fn) {
        this.onData = fn;
    }
}

export default WebRTC;