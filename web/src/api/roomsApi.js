import baseApi from './baseApi'


export const createRoom = (uuid,gameID,peer) => {
    const json = JSON.stringify({ uuid:uuid, game_id:gameID, peer:peer});
    return baseApi.post(`/room/create`, { credentials: 'include' }, json);
};
export const getRooms = () => {
    return baseApi.get(`/room`,{credentials: 'include'})
}

export default {
    createRoom,
    getRooms,
}