import constants from "@/common/constants";
import log from "@/common/log";
import eventBus from "@/common/event-bus";
import { startUpdatingGamepadState, stopUpdatingGamepadState } from "@comp/templates/input";

class WebRTC {
    private conn = new RTCPeerConnection({
        iceServers: [
            { urls: "stun:stun.l.google.com:19302" },
            { urls: "stun:stun1.l.google.com:19302" }
        ],
        
        iceCandidatePoolSize: 10 // Увеличенный пул ICE кандидатов
    });
    private dataChannel: RTCDataChannel | undefined;
    private display: HTMLVideoElement;

    constructor(display: HTMLVideoElement) {
        this.display = display;
    }

    init(offer: RTCSessionDescriptionInit) {
        this.conn.ontrack = (e) => {
            if (e.streams && e.streams[0]) {
                this.display.srcObject = e.streams[0];
                this.display.autoplay = true;
                this.display.muted = true;
                this.display.playsInline = true;
                this.display.controls = true;
                this.display.play();
            } else {
                log.warn("No video stream found");
            }
        };

        // Добавление трансиверов для видео и аудио
        const transceiver = this.conn.addTransceiver('video', { direction: 'recvonly' });
        const params = transceiver.sender.getParameters();
        params.encodings[0].maxBitrate = 30000000; // 30 Mbps
        transceiver.sender.setParameters(params);

        this.conn.addTransceiver("audio", { direction: "recvonly" });

        this.conn.ondatachannel = (e) => {
            log.debug(`[rtc] ondatachannel, ${e.channel.label}`);
            this.dataChannel = e.channel;
            this.dataChannel.onopen = () => {
                log.debug("[rtc] the input channel has been opened");
                startUpdatingGamepadState(this.dataChannel!.send.bind(this.dataChannel));
            };
            this.dataChannel.onerror = (e) => {
                log.error(`[rtc] the input channel has been closed due to error ${e}`);
            };
            this.dataChannel.onclose = () => {
                log.info("[rtc] the input channel has been closed");
            };
        };

        this.conn.oniceconnectionstatechange = (e) => {
            log.info(`[rtc] State: ${this.conn.iceConnectionState}`);
            switch (this.conn.iceConnectionState) {
                case "connected":
                    eventBus.notify(constants.RTC_EVENT.CONNECTION_READY);
                    this.startLoggingStats(); // Запуск логирования статистики
                    break;
                case "disconnected":
                case "failed":
                    this.stop();
                    break;
            }
        };

        this.conn.onicecandidate = (e) => {
            if (e.candidate) {
                log.warn(`[rtc] ICE candidate found: ${e.candidate.candidate}`);
                eventBus.notify(constants.RTC_EVENT.ICE_CANDIDATE_FOUND, btoa(JSON.stringify(e.candidate)));
            }
        };

        this.conn.onicegatheringstatechange = (e) => {
            log.info(`[rtc] Gathering state: ${this.conn.iceGatheringState}`);
            switch (this.conn.iceGatheringState) {
                case "complete":
                    eventBus.notify(constants.RTC_EVENT.ICE_GATHERING_COMPLETE);
                    break;
            }
        };

        this.setSDP(offer);
    }

    addice(candidate: RTCIceCandidate) {
        this.conn.addIceCandidate(candidate)
            .then(() => {
                log.info("[rtc] Candidate added");
            })
            .catch((e: any) => {
                log.error(`[rtc] Error adding candidate: ${e}`);
            });
    }

    stop() {
        stopUpdatingGamepadState();
        this.conn.close();
    }

    sendData(data: string) {
        console.log(this.dataChannel);
    }

    private async setSDP(offer: RTCSessionDescriptionInit) {
        await this.conn
            .setRemoteDescription(offer)
            .then(() => {
                log.info("[rtc] Remote description set");
            })
            .catch((e: any) => {
                log.error(`[rtc] Error setting remote description: ${e}`);
            });
        const answer = await this.conn.createAnswer();
        await this.conn
            .setLocalDescription(answer)
            .then((): void => {
                log.info("[rtc] Local description set");
                eventBus.notify(constants.RTC_EVENT.SDP_ANSWER_CREATED, btoa(JSON.stringify(answer)));
            })
            .catch((e: any) => {
                log.error(`[rtc] Error setting local description: ${e}`);
            });
    }

    // Метод для логирования статистики
    private async startLoggingStats() {
		try {
			while (this.conn.signalingState !== 'closed') {
				const stats = await this.conn.getStats();
				stats.forEach(report => {
					if (report.type === 'inbound-rtp') {
						log.info(`[rtc] ${report.type} - ${report.id}: ${JSON.stringify(report)}`);
					} else if (report.type === 'outbound-rtp') {
						log.info(`[rtc] ${report.type} - ${report.id}: ${JSON.stringify(report)}`);
					} else if (report.type === 'transport') {
						log.info(`[rtc] ${report.type} - ${report.id}: ${JSON.stringify(report)}`);
					} else if (report.type === 'candidate-pair') {
						log.info(`[rtc] ${report.type} - ${report.id}: ${JSON.stringify(report)}`);
					}
				});
				await new Promise(resolve => setTimeout(resolve, 5000)); // Логируем каждые 5 секунд
			}
		} catch (error) {
			log.error(`[rtc] Error logging stats: ${error}`);
		}
	}	
}

export default WebRTC;
