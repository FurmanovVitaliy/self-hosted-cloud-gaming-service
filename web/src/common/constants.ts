const constants = {
	log: {
		enabled: true,
		level: 3,
	},
	api: {
		server: "https://192.168.1.13:8000",
		login: "/signin",
		signup: "/signup",
		logout: "/logout",
		games: "/games",
		roomCreeate: "/room/create",
		rooms: "/rooms",
		roomIsAlive: (uuid: string, username: string) => `/room/state/${uuid}?username=${username}`,
	},
	websocket: {
		server: "wss://192.168.1.13:8000",
		joinRoom: (uuid: string, username: string) => `/room/join/${uuid}?username=${username}`,
	},
	routes: {
		home: "/",
		library: "/library",
		game: "/library/:uuid",
		rooms: "/rooms",
		room: "/rooms/:uuid",
		faq: "/faq",
		login: "/login",
		singup: "/signup",
	},
	search: {
		types: {
			games: "games",
			rooms: "rooms",
		},
	},
	storage: {
		keys: {
			token: "token",
		},
	},
	iceConfig: {
		iceServers: [
			{
				urls: "stun:stun.l.google.com:19302",
			},
		],
	},
	RTC_EVENT: {
		CONNECTION_CLOSED: "rtcConnectionClosed",
		CONNECTION_READY: "rtcConnectionReady",
		ICE_CANDIDATE_FOUND: "rtcIceCandidateFound",
		ICE_CANDIDATE_RECEIVED: "rtcIceCandidateReceived",
		ICE_GATHERING_COMPLETE: "rtcIceGatheringComplete",
		SDP_ANSWER_CREATED: "rtcSdpAnswer",
		SDP_OFFER_RECEIVED: "rtcSdpOffer",
		INPUT_READY: "inputReady",
	},

	WS_MSG_TAG: {
		WEBRTC: "wrtc",
		CHAT: "chat",
		ERROR: "error",
		DEVICE_INFO: "deviceInfo",
	},

	RTC_CONTENT_TYPE: {
		OFFER: "offer",
		ANSWER: "answer",
		CANDIDATE: "candidate",
		SERVER_ICE_READY: "server_ice_ready",
		CLENT_ICE_READY: "client_ice_ready",
		CONNECTION_READY: "connection_ready",
	},

	input: {
		CHECK_CONNECTION_INTERVAL: 2000,
		UPDATE_INTERVAL: 50,
	},
};

export default constants;
