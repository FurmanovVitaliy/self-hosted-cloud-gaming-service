 const appConstants = {
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
    storage:{
        keys:{
            token:'token',
        }
    }
}

export default appConstants