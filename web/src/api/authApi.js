import baseApi from './baseApi';

export const login = (email, password) => {
    const json= JSON.stringify({email, password})
    console.log(json)
    return fetch ('http://localhost:8000/login', {
        credentials: "include",
        method: 'POST',
        headers: {
            
        },
        body: json,
    })
}
export default {
    login
}