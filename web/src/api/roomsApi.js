import baseApi from './baseApi'


export const createRoom = (uuid,game, peer) => {
    const json = JSON.stringify({ game: game,uuid: uuid, peer: peer });
    return baseApi.post(`/room/create`, { credentials: 'include' }, json);
};
export const getRooms = () => {
    return baseApi.get(`/room`,{credentials: 'include'})
}

export default {
    createRoom,
    getRooms,
}