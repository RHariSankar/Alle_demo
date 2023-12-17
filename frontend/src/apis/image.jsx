import axios from "axios";

const baseUrl = "http://localhost:8080/api/v1/image"

export function uploadImage(image) {
    const formData = new FormData()
    formData.append('image', image)
    return axios.post(baseUrl, formData)
}

export function getImage(id) {
    let url = baseUrl + "/" + id
    return axios.get(url, { responseType: 'blob' })
}