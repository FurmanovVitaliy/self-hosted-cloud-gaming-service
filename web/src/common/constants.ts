const constants = {
  SERVER_IP: "192.168.1.13:8000",
  routes: {
    index: "/",
    games: "/games",
    rooms: "/rooms",
    room: "/rooms/:uuid",
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
  RTC_EVENT: {
    RTC_CONNECTION_CLOSED: "rtcConnectionClosed",
    RTC_CONNECTION_READY: "rtcConnectionReady",
    RTC_ICE_CANDIDATE_FOUND: "rtcIceCandidateFound",
    RTC_ICE_CANDIDATE_RECEIVED: "rtcIceCandidateReceived",
    RTC_SDP_ANSWER_CREATED: "rtcSdpAnswer",
    RTC_SDP_OFFER_RECEIVED: "rtcSdpOffer",
    RTC_INPUT_READY: "inputReady",
  },
  RTC_SIGNAL: {
    SERVER_ICE_GATHERING_COMPLETE : "sigc",
    CONNECTION_READY: 'connectionReady',
  },

  RTC_ENTITY_NAME: {
    OFFER: "offer",
    ANSWER: "answ",
    SERVER_CANDIDATE: "sc",
    CLIENT_CANDIDATE: "cc",
    SIGNAL: "sig",
    ERROR: "err",
  },
};


export default constants;
