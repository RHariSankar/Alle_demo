import { useEffect, useState } from "react"
import { getImage } from "../apis/image"
import { Grid } from "@mui/material"


function Image(props) {
    console.log("image", props);
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
        <Grid container spacing={2}>
            <Grid item xs={12} md={12}>
                <div>
                    <img alt="not found"
                        width={"500px"}
                        src={imageUrl}>
                    </img>
                </div>
            </Grid>
            <Grid item xs={12} md={12} style={{ display: "flex", flexDirection: "row-reverse" }}>
                {
                    props?.data?.tags.join(" ")
                }
            </Grid>
        </Grid>

    )

}

export default Image