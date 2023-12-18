import axios from "axios";

const baseUrl = "http://localhost:8080/api/v1/image"

export function uploadImage(image, tags) {
    console.log("tags", tags)
    const formData = new FormData()
    formData.append('image', image)
    tags.forEach(tag => {
        formData.append("tag", tag)
    });
    return axios.post(baseUrl, formData)
}

export function getImage(id) {
    let url = baseUrl + "/" + id
    return axios.get(url, { responseType: 'blob' })
}