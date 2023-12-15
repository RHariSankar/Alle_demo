import { Box, FormLabel, Grid, Typography } from '@mui/material';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';


function ChatDialog(props) {

    let gridCss = { display: 'flex', flexDirection: 'row', marginBottom: '1%' }
    let boxCss = { border: 'solid', backgroundColor: 'white', padding: '0.5rem', paddingRight: '2%', fontSize: '16px' }
    let typogrpahCss = { maxWidth: '30vw', wordWrap: 'break-word' }
    if (props?.data?.role === 'system') {
        gridCss.justifyContent = 'flex-start'
        typogrpahCss.textAlign = 'left'
    } else {
        gridCss.justifyContent = 'flex-end'
        typogrpahCss.textAlign = 'right'
    }
    console.log('dialog', props)
    return (
        <Grid item xs={12} md={12} style={gridCss}>
            <Box style={boxCss}>
                {props?.data?.type === 'chat' ?
                    <Typography variant="body1" style={typogrpahCss}>
                        {props?.data?.text}
                    </Typography> :
                    <img alt="not found"
                        width={"250px"}
                        src={props?.data?.text}
                    />
                }

            </Box>
        </Grid>
    )

}

export default ChatDialog