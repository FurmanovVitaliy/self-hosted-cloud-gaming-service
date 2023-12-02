import appConstants from "../common/constants"

const apiServer = import.meta.env.DEV ? 'https://192.168.1.13:8000' : 'https://192.168.1.13:8000'

const baseFetch = (url, config = {}, params) => {
    return new Promise((resolve, reject) => {
        try{
            /*const token = window.sessionStorage.getItem(appConstants.storage.keys.token)*/
            /*const defaultHeaders = {
                'Content-Type': 'application/json',
            }
            /*if (token){
                defaultHeaders['Authorization'] = `Bearer ${token}`
            }*/
            const _config = {
                /*headers: {
                    ...defaultHeaders
                },*/
                ...config
            }
            if(params){
                _config['body'] = JSON.stringify(params)
            }
            window.fetch(`${apiServer}${url}`, {
                ..._config
            }).then(response=>response.json())
            .then(resolve, reject)
        } catch(e){
            reject(e)
        }
    })
}

const fetchGet = (url, config) => {
    return baseFetch(url, config)
}

const fetchPost = (url, config = {},json ={}) => {

    return baseFetch(url, {
        ...config,
        method: 'POST',
        body: json
    })
}

const fetchPut = (url, params ={}, config = {}) => {
    return baseFetch(url, {
        ...config,
        method: 'PUT'
    }, params)
}
const fetchPatch = (url, params ={}, config = {}) => {
    return baseFetch(url, {
        ...config,
        method: 'PATCH'
    }, params)
}

const fetchDelete = (url, params ={}, config = {}) => {
    return baseFetch(url, {
        ...config,
        method: 'DELETE'
    }, params)
}

export default {
    get: fetchGet,
    post: fetchPost,
    put: fetchPut,
    patch: fetchPatch,
    delete: fetchDelete,
}