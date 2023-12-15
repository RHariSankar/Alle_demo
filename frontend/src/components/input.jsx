import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';
import UploadFileIcon from '@mui/icons-material/UploadFile';
import IconButton from '@mui/material/IconButton';
import SendIcon from '@mui/icons-material/Send';
import { FormControl } from '@mui/material';
import { useState } from 'react';
import { styled } from '@mui/material/styles';

let date = new Date()

const VisuallyHiddenInput = styled('input')({
    clip: 'rect(0 0 0 0)',
    clipPath: 'inset(50%)',
    height: 1,
    overflow: 'hidden',
    position: 'absolute',
    bottom: 0,
    left: 0,
    whiteSpace: 'nowrap',
    width: 1,
});


function Input(props) {

    const [text, setText] = useState('')

    function handleTextChange(event) {
        setText(event.target.value)
    }

    function handleSubmit(event) {
        event.preventDefault();
        console.log('handle submit from input: ', text)
        let inputData = {
            type: 'chat',
            role: 'user',
            text: text,
            dateTime: date.toISOString()
        }
        //TODO: async call api to get response from model 
        props?.submit(event, inputData);
        setText('')
        //TODO: call props?.submit again with response
    }

    function fileUpload(event) {
        console.log('file upload', event.target.files[0]);
        let url = URL.createObjectURL(event.target.files[0])
        // call post api to save image
        console.log('image', url)
        let imageData = {
            type: 'image',
            role: 'user',
            text: url,
            dateTime: date.toISOString()
        }
        props.submit(event, imageData)
    }

    return (
        <FormControl style={{ width: '100%', marginLeft: "3%", marginRight: "3%" }}>
            <form onSubmit={handleSubmit}>
                <Grid container spacing={2}>
                    <Grid item xs={11} md={11}>
                        <TextField
                            fullWidth
                            inputProps={{ min: 0, style: { textAlign: 'right', backgroundColor: 'white' } }}
                            id="outlined-multiline-flexible"
                            label="Type here"
                            onChange={handleTextChange}
                            value={text}
                        />

                    </Grid>
                    <Grid item xs={1} md={1} style={{ display: 'flex', flexDirection: 'row', alignItems: 'self-start', }}>
                        <Grid item xs={6} md={6} >
                            <IconButton style={{ color: 'black' }} type='submit'>
                                <SendIcon />
                            </IconButton>
                        </Grid>
                        <Grid item xs={6} md={6} >
                            <IconButton component="label" variant="contained" style={{ color: 'black' }}>
                                <UploadFileIcon />
                                <VisuallyHiddenInput type="file" accept="image/*" onChange={fileUpload} />
                            </IconButton>
                        </Grid>
                    </Grid>

                </Grid>
            </form>
        </FormControl>

    )
}

export default Input
