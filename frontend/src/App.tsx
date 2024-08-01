import Editor from "@monaco-editor/react";
import Button from "@mui/material/Button";
import Stack from "@mui/material/Stack";
import Textarea from "@mui/joy/Textarea";
import './App.css'

function App() {

  return (
      <>
          <h2>[MIA] Proyecto 1 Diego Cali</h2>
          <Editor
              height="45vh"
              width="180vh"
              theme="vs-dark"
              loading="Loading..."
              defaultLanguage="javascript"
              defaultValue="// some comment"
          />
          <Stack spacing={2} direction="row">
              <Button variant="contained">File</Button>
              <Button variant="contained">Run</Button>
          </Stack>
          <Textarea></Textarea>
      </>
  )
}

export default App
