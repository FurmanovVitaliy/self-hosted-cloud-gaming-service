import { msg, getClientDeviceInfo } from "@/common/utils";
import { MessageJSON } from "@/types/types";
//import { startUpdatingGamepadState } from "@/components/templates/input";
import WebRTC from "@/components/templates/webrtc";
import Page from "@/components/templates/page";
import constants from "@/common/constants";
import eventBus from "@/common/event-bus";
import log from "@/common/log";

class RoomPage extends Page {
	public title: string = "Room | Pixel Cloud";
	private html = (
	<null>
		<video class='main-display'></video>
		<div id="myAlert" style="display: flex;width: max-content; position: fixed; top: 50%; left: 50%; transform: translate(-50%, -50%); padding: 20px; ; z-index: 1000;">
    	<p>Please connect your gamepad	(other input devices are not supported yet )</p>
		</div>
	</null>
);

	private webrtc: WebRTC | null = null;
	private display: HTMLVideoElement;
	private websocket: WebSocket | null = null;
	private username: string = localStorage.getItem("username")!;
	private alert : HTMLElement | null = null;
	constructor() {
		super();
		this.loadTemplate(this.html);
		this.addEventListener("", "popstate", this.evetnGoOut);
		this.addEventListener("", "beforeunload", this.evetnGoOut);
		this.display = this.container.querySelector(".main-display")!;
		this.alert = this.container.querySelector("#myAlert");
		
	}

	evetnGoOut() {
		this.websocket!.close();
		this.webrtc!.stop();
		this.webrtc = null;
		this.websocket = null;
	}

	
	initControlDevice() {
		// Флаг для отслеживания того, был ли вызван upgradeconnection
		this.alert?.style.setProperty("display","flex");
		let isConnectionUpgraded = false;

		//alert("Please connect a gamepad to control the game");

		window.addEventListener("gamepadconnected", (e) => {
			console.log("Gamepad connected");

			// Если upgradeconnection уже был вызван, ничего не делаем
			if (isConnectionUpgraded) return;

			// Вызов upgradeconnection и установка флага
			this.alert?.style.setProperty("display","none");
			this.upgradeConnection();
			isConnectionUpgraded = true;
		});
	}
	

	upgradeConnection() {	
		
		let deviceInfoJson = getClientDeviceInfo();
		console.log(deviceInfoJson);
		const uuid = window.location.pathname.split("/")[2];
		this.websocket = new WebSocket(constants.websocket.server + constants.websocket.joinRoom(uuid, this.username));
		this.webrtc = new WebRTC(this.display);

		this.websocket.onopen = () => {
			this.websocket!.send(JSON.stringify(deviceInfoJson));
      log.debug("[ws]->Device info sent");
      
      eventBus.execute(constants.RTC_EVENT.CONNECTION_READY,() => {
        this.websocket!.send(msg(constants.WS_MSG_TAG.WEBRTC,constants.RTC_CONTENT_TYPE.CONNECTION_READY));
          log.info(`[ws]->Connection ready signal was sent`)},true);
      
      eventBus.execute(constants.RTC_EVENT.ICE_GATHERING_COMPLETE,() => {
        this.websocket!.send(msg(constants.WS_MSG_TAG.WEBRTC,constants.RTC_CONTENT_TYPE.CLENT_ICE_READY));
         log.debug(`[ws]->Ice ready signal send`)},true);


      eventBus.execute(constants.RTC_EVENT.SDP_ANSWER_CREATED,(answer) => {
        this.websocket!.send(msg(constants.WS_MSG_TAG.WEBRTC,constants.RTC_CONTENT_TYPE.ANSWER,answer));
        log.info(`[ws]->Answer was sent`)},true);

      eventBus.execute(constants.RTC_EVENT.ICE_CANDIDATE_FOUND,(candidate) => {
        this.websocket!.send(msg(constants.WS_MSG_TAG.WEBRTC,constants.RTC_CONTENT_TYPE.CANDIDATE, candidate));
        log.warn(`[ws]->Candidate was send`);});
		};
		this.websocket.onmessage = (e) => {
			const data = JSON.parse(e.data) as MessageJSON;
			switch (data.tag) {
				case constants.WS_MSG_TAG.CHAT:
					console.log(data.content);
					break;
				case constants.WS_MSG_TAG.ERROR:
					console.log(data.content);
					break;
				case constants.WS_MSG_TAG.WEBRTC:
					switch (data.content_type) {
						case constants.RTC_CONTENT_TYPE.OFFER:
							console.log("[ws] <- Offer received");
							this.webrtc?.init(JSON.parse(atob(data.content)));
							break;
						case constants.RTC_CONTENT_TYPE.CANDIDATE:
							console.log(`[ws] <- Candidate received : ${atob(data.content)}`);
							this.webrtc?.addice(JSON.parse(atob(data.content)));
							break;
						case constants.RTC_CONTENT_TYPE.SERVER_ICE_READY:
							console.log("[ws] <- Server ICE ready signal received");
							break;
					}
					break;
				default:
					console.log("Unknown message: " + data);
			}
		};

		this.websocket.onclose = () => {log.info("[ws] Socket is closed");};
		this.websocket.onerror = (error) => {log.error(`[ws] Socket error: ${error}`);};
	}

	

	async render() {
		this.initControlDevice();
		return this.container;
	}
}

export default RoomPage;
