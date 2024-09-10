import baseApi from './baseApi';

export const login = (email, password) => {
    const json= JSON.stringify({email, password})
    return fetch ('https://192.168.1.13:8000/login', {
        credentials: "include",
        method: 'POST',
        headers: {},
        body: json,
    })
}
export default {
    login
}