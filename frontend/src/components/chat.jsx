import { Box } from '@mui/material';
import Grid from '@mui/material/Grid';
import ChatDialog from './dialog';


function Chat(props) {
    let gridCss = { width: '100%', height: '80%', overflowY: 'auto', border: 'solid', margin: '10px' }
    return (
        <Grid item xs={12} md={12} sx={gridCss}>
            <Box sx={{ width: '100%', padding: '2%' }}>
                {props?.inputData.length === 0 ? <h2>No Chat History</h2> : <div>
                    {props?.inputData.map((e) => <ChatDialog data={e}></ChatDialog>)}
                </div>}
            </Box>
        </Grid>
    )

}

export default Chat