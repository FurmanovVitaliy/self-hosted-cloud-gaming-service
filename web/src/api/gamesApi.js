import baseApi from './baseApi'

export const getGames = (page) => {
    return baseApi.get(`/games`,{credentials: 'include'})
}

export default {
    getGames,
}