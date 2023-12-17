import Grid from '@mui/material/Grid';
import Chat from './chat';
import Input from './input';
import React, { useEffect, useState } from 'react';
import { getAllChats } from '../apis/chat';


let sampleData = [
    {
        "type": "chat",
        "role": "user",
        "data": "hi show me links of red tshirts",
        "dateTime": "2023-12-15T05:50:07.225Z"
    },
    {
        "type": "chat",
        "role": "system",
        "data": "here you go",
        "dateTime": "2023-12-15T05:50:07.225Z"
    },
    {
        "type": "chat",
        "role": "user",
        "data": "thanks",
        "dateTime": "2023-12-15T05:50:07.225Z"
    }
]

const generateRandomString = (length) => {
    let result = '';
    const characters =
        'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    const charactersLength = characters.length;
    for (let i = 0; i < length; i++) {
        result += characters.charAt(Math.floor(Math.random() * charactersLength));
    }
    return result;
};




function View(props) {

    const [inputValue, setInputValue] = useState([]);

    async function reply() {
        let res = generateRandomString(250)
        let date = new Date()
        let op = {
            type: 'chat',
            role: 'system',
            data: res,
            dateTime: date.toISOString()
        }
        return op
    }

    useEffect(() => {
        //TODO: call get api to get chat and render
        // setInputValue(sampleData)
        getAllChats().then((response) => {
            console.log("api response", response)
            setInputValue(response?.data)
        })
    }, [])


    const inputSubmit = (event, data) => {
        event.preventDefault();
        // TODO: api to save message to server
        console.log('view submit: ', data)
        // setInputValue((prevState) => [...prevState, ...data]);
        setInputValue((prevState) => prevState.concat(data));
        // if (data?.type == 'image') {

        // } else {
        //TODO: api to get reply
        //mock reply
        // reply().then(op => setInputValue((prevState => [...prevState, op])))
        // }

    };

    return (

        <Grid container rowSpacing={1} columnSpacing={1} style={{ height: '100vh' }}>
            {/* <div>{JSON.stringify(inputValue)}</div> */}
            <Chat inputData={inputValue}></Chat>
            <Input submit={inputSubmit}></Input>
        </Grid>

    )
}

export default View