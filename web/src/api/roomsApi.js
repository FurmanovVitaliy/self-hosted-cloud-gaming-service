import baseApi from './baseApi'

export const getRooms = (page) => {
    return baseApi.get(`/room`,{credentials: 'include'})
}

export default {
    getRooms,
}