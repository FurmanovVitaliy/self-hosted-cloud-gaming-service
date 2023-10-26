const rooms = new Map();
 
export const getRooms = (roomID) => {
    return rooms.get(roomID);
}

export const setRooms = (room) => {
    rooms.set( room);
}

export default {
    getRooms,
    setRooms
}