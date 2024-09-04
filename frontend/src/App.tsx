import Editor from "@monaco-editor/react";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import CssBaseline from "@mui/material/CssBaseline";
import {createTheme, ThemeProvider} from "@mui/material";
import {useRef, useState} from "react";
import './App.css'
import {POST} from "./Post.ts";

const darkTheme = createTheme({
    palette: {
        mode: 'dark',
    },
});

function App() {
    const editorRef = useRef(null);
    const [file, setFile] = useState<File>(null);

    function openFile(){
        if (file){
            const reader = new FileReader();
            reader.onload = function(e) {
                // eslint-disable-next-line @typescript-eslint/ban-ts-comment
                // @ts-expect-error
                editorRef.current.setValue(e.target.result);
            }
            reader.readAsText(file);
        } else {
            alert("No file selected")
        }
    }

    function handleEditorDidMount(editor: any){
        editorRef.current = editor;
    }

    function sendText(){
        // eslint-disable-next-line @typescript-eslint/ban-ts-comment
        // @ts-expect-error
        const dataToSend = {code: editorRef.current.getValue()}
        POST("http://localhost:8080/run-code", dataToSend).then((result) => {
            const output = document.getElementById("Console");
            // eslint-disable-next-line @typescript-eslint/ban-ts-comment
            // @ts-expect-error
            output.value = result.received;
        });
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
                  <input onChange={ (e) => {setFile(e.target.files[0])}} type="file"/>
                  <Button variant="outlined" onClick={openFile}>Open File</Button>
                  <Button variant="outlined" onClick={sendText}>Run</Button>
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
