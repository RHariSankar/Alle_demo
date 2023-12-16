import axios from "axios";


const baseUrl = "http://localhost:8080/api/v1/chat"

export function getAllChats() {
    let url = baseUrl + '/all'
    return axios.get(url)
}

export function postChat(body) {
    let url = baseUrl
    return axios.post(url, body)
}