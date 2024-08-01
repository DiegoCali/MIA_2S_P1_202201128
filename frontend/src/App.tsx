import Editor from "@monaco-editor/react";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import CssBaseline from "@mui/material/CssBaseline";
import {createTheme, ThemeProvider} from "@mui/material";
import {useRef} from "react";
import './App.css'

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

function App() {
    const editorRef = useRef(null);

    function handleEditorDidMount(editor: any){
        editorRef.current = editor;
    }

    function addText(){
        document.getElementById("Console").innerText = editorRef.current.getValue()
    }
  return (

      <>
          <ThemeProvider theme={darkTheme}>
              <CssBaseline />
              <h2>[MIA] Proyecto 1 Diego Cali</h2>
              <Editor
                  height="45vh"
                  width="180vh"
                  theme="vs-dark"
                  loading="Loading..."
                  defaultLanguage="javascript"
                  defaultValue="// some comment"
                  onMount={handleEditorDidMount}
              />
              <Stack spacing={2} direction="row">
                  <Button variant="outlined">File</Button>
                  <Button variant="outlined" onClick={addText}>Run</Button>
              </Stack>
              <textarea
                  id="Console"
                  style={{
                      width: "100%",
                      height: "100%",
                  }}
                  rows={5}
              ></textarea>
          </ThemeProvider>
      </>
  )
}

export default App
