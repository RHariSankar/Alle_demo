import logo from './logo.svg';
import './App.css';
import Input from './components/input'
import { ThemeProvider } from '@mui/material/styles';
import { CssBaseline } from '@mui/material';
import theme from './theme';
import View from './components/view';

function App() {
  return (
    <div className="App">
      <header className="App-header">
        <ThemeProvider theme={theme}>
          <CssBaseline />
          <View></View>

        </ThemeProvider>


      </header>
    </div>
  );
}

export default App;
