import Grid from '@mui/material/Grid';
import TextField from '@mui/material/TextField';
import UploadFileIcon from '@mui/icons-material/UploadFile';
import IconButton from '@mui/material/IconButton';
import SendIcon from '@mui/icons-material/Send';
import { FormControl } from '@mui/material';
import { useState } from 'react';
import { styled } from '@mui/material/styles';
import { postChat } from '../apis/chat';
import { uploadImage } from '../apis/image';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogContentText from '@mui/material/DialogContentText';
import DialogTitle from '@mui/material/DialogTitle';
import Button from '@mui/material/Button';
import CloseIcon from '@mui/icons-material/Close';
import Chip from '@mui/material/Chip';
import Stack from '@mui/material/Stack';

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
    const [imageUploadPopup, setImageUploadPopup] = useState(false)
    const [imageBlobUrl, setImageBlobUrl] = useState("")
    const [selectedFile, setSelectedFile] = useState(null);
    const [tags, setTags] = useState([]);
    const [inputValue, setInputValue] = useState('');


    const handleInputChange = (event) => {
        setInputValue(event.target.value);
    };

    const handleInputKeyDown = (event) => {
        if (event.key === 'Enter' && inputValue.trim() !== '') {
            setTags([...tags, inputValue.trim()]);
            setInputValue('');
        }
    };

    const handleDeleteTag = (tagToDelete) => {
        setTags(tags.filter((tag) => tag !== tagToDelete));
    };

    function handleTextChange(event) {
        setText(event.target.value)
    }

    function handleSubmit(event) {
        event.preventDefault();
        console.log('handle submit from input: ', text)
        let inputData = {
            type: 'chat',
            role: 'user',
            data: text,
            dateTime: date.toISOString()
        }
        props?.submit(event, inputData);
        setText('')
        postChat(inputData).then((response) => {
            console.log(response)
            props?.submit(event, response.data)
        }).catch((error) => {
            console.error("couldn't send chat", error)
        })
    }

    function fileUploadPopup(event) {
        console.log('file upload', event.target.files[0]);
        setSelectedFile(event.target.files[0])
        let blolUrl = URL.createObjectURL(event.target.files[0])
        setImageBlobUrl(blolUrl)
        setImageUploadPopup(true)
    }
    function handleClose() {
        setImageUploadPopup(false)
    }

    function fileUpload(event) {
        if (!selectedFile)
            return
        uploadImage(selectedFile, tags).then((response) => {
            console.log('image uploaded', response.data)
            props.submit(event, response.data)
        }).catch((error) => {
            console.error('Error uploading image', error)
        })
        setImageBlobUrl("")
        setImageUploadPopup(false)
        setSelectedFile(null)
        setTags([])

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
                                <VisuallyHiddenInput type="file" accept="image/*" onChange={fileUploadPopup} />
                            </IconButton>
                        </Grid>
                    </Grid>

                </Grid>
            </form>
            <Dialog
                open={imageUploadPopup}
                onClose={handleClose}
                aria-labelledby="alert-dialog-title"
                aria-describedby="alert-dialog-description"
            >
                <DialogTitle id="alert-dialog-title">
                    <Grid item xs={12} md={12} style={{ display: 'flex', flexDirection: 'row', alignItems: 'self-start', }}>
                        <Grid item xs={11} md={11} >
                            Image Upload
                        </Grid>
                        <Grid item xs={1} md={1} >
                            <IconButton onClick={handleClose}>
                                <CloseIcon />
                            </IconButton>
                        </Grid>
                    </Grid>
                </DialogTitle>
                <DialogContent>
                    <DialogContentText id="alert-dialog-description" style={{ marginBottom: '10px' }}>
                        <div>
                            <img alt="not found"
                                width={"500px"}
                                src={imageBlobUrl}>
                            </img>
                        </div>
                    </DialogContentText>
                </DialogContent>
                <DialogActions>
                    <Grid container>
                        <Grid item xs={12} md={12} >
                            <TextField
                                label="Add Tags"
                                variant="outlined"
                                value={inputValue}
                                onChange={handleInputChange}
                                onKeyDown={handleInputKeyDown}
                                fullWidth
                                required
                            />
                            <br />
                            <br />
                            <Stack direction="row" spacing={1}>
                                {tags.map((tag, index) => (
                                    <Chip
                                        key={index}
                                        label={tag}
                                        onDelete={() => handleDeleteTag(tag)}
                                        color="primary"
                                    />
                                ))}
                            </Stack>
                        </Grid>
                        <Grid container xs={12} md={12} style={{ justifyContent: 'end' }}>
                            <Button onClick={fileUpload} autoFocus disabled={tags.length === 0}>
                                Upload
                            </Button>
                        </Grid>
                    </Grid>
                </DialogActions>
            </Dialog>
        </FormControl>


    )
}

export default Input
