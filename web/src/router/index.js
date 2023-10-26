import appConstants from "../common/constants";
import Route from 'route-parser'
import MainPage from '../pages/main.template'
import GamesPage from '../pages/games.template'
import RoomsPage from '../pages/rooms.template'
import Room from "../pages/room.template";


export const routes = {
    Main: new Route(appConstants.routes.index),
    Games: new Route(appConstants.routes.games),
    Rooms: new Route(appConstants.routes.rooms),
    Room: new Route(appConstants.routes.room),
}

const routesPages = [
    {route: routes.Main, page: MainPage},
    {route: routes.Games, page: GamesPage},
    {route: routes.Rooms, page: RoomsPage},
    {route: routes.Room, page: Room},
]

export const getPathRoute =(path) => {
    const target = routesPages.find(r => r.route.match(path))
    if (target) {
        const params = target.route.match(path)
        return {
            page: target.page,
            route: target.route,
            params
        }
    }
    return null
}

export const render = (path) => {
    let result = '<h1>404</h1>'

    const pathRoute = getPathRoute(path)

    if (pathRoute) {
        result = pathRoute.page(pathRoute.params)
    }

    document.querySelector('#app').innerHTML = result
}

export const navigateTo = (path) => {
    window.history.pushState({ path }, path, path)
    render(path)
}

export const getRouterParams = () => {
    const url = new URL(window.location.href).pathname
    return getPathRoute(url)
}

const initRouter = () => {
    window.addEventListener('popstate', e => {
        render( new URL(window.location.href).pathname)
    })
    document.querySelectorAll('[href^="/"]').forEach(el => {
        el.addEventListener('click', (env) => {
            env.preventDefault()
            const {pathname: path} = new URL(env.target.href)
            navigateTo(path)
        })
    })
    render(new URL(window.location.href).pathname)
}

export default initRouter