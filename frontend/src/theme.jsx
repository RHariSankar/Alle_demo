import { createTheme } from "@mui/material/styles";

const theme = createTheme({
    palette: {
        mode: 'light',
        text: {
            primary: "#000000",
            secondary: 'rgba(0, 0, 0, 0.7)',
            disabled: 'rgba(0, 0, 0, 0.5)',
        },
        background: {
            default: "#cfdefc"
        }
    },
    typography: {
        allVariants: {
            color: "#000000",
        },
        fontFamily: [
            '-apple-system',
            'BlinkMacSystemFont',
            '"Segoe UI"',
            'Roboto',
            '"Helvetica Neue"',
            'Arial',
            'sans-serif',
            '"Apple Color Emoji"',
            '"Segoe UI Emoji"',
            '"Segoe UI Symbol"',
        ].join(','),
    },

});

export default theme