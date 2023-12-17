import { useEffect, useState } from "react"
import { getImage } from "../apis/image"



function Image(props) {

    let imageId = props?.data?.data

    const [imageUrl, setImageUrl] = useState('')

    useEffect(() => {
        getImage(imageId).then((response) => {
            const imageBlobUrl = URL.createObjectURL(new Blob([response.data]))
            console.log('image retrieved', imageBlobUrl)
            setImageUrl(imageBlobUrl)
        })
    }, [])

    return (
        <div>
            <img alt="not found"
                width={"500px"}
                src={imageUrl}>
            </img>
        </div>
    )

}

export default Image