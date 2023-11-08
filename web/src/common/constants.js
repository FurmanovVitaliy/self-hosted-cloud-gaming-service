const constants = {
    routes: {
        index: '/',
        games: '/games',
        rooms: '/rooms',
        room: '/rooms/:uuid'
    },
    search: {
        types: {
            games: 'games',
            rooms: 'rooms',
        }
    },
    storage: {
        keys: {
            token: 'token',
        }
    },
    events: {
            RTC_CONNECTION_CLOSED: 'rtcConnectionClosed',
            RTC_CONNECTION_READY: 'rtcConnectionReady',
            RTC_ICE_CANDIDATE_FOUND: 'rtcIceCandidateFound',
            RTC_ICE_CANDIDATE_RECEIVED: 'rtcIceCandidateReceived',
            RTC_NEW_CONNECTION: 'rtcNewConnection',
            RTC_SDP_ANSWER_CREATED: 'rtcSdpAnswer',
            RTC_SDP_OFFER_RECIVED: 'rtcSdpOffer',
    }
};

export default constants;
