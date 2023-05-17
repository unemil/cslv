import { useCallback, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import axios from 'axios'
import './App.css'

const CaptchaSolve = () => {
  const [result, setResult] = useState('')

  const sendData = async (file) => {
    const formData = new FormData()
    formData.append('file', file)

    await axios
      .request({
        method: 'POST',
        url: 'http://80.90.186.78/api/v1/captcha/solve',
        data: formData,
        headers: {
          'Content-Type': 'multipart/form-data',
        }
      })
      .then(response => { setResult(response.data.text) })
      .catch(error => { setResult(error.response.data.error) })
  }

  const onDrop = useCallback((acceptedFiles) => {
    acceptedFiles.forEach((file) => {
      sendData(file)
    })
  }, [])
  const { getRootProps, getInputProps } = useDropzone({ onDrop })

  return (
    <>
      <h1>cslv</h1>
      <div id='dropzone' {...getRootProps()}>
        <input {...getInputProps()} />
        <p>Drag and drop file here, or click to select file</p>
      </div>
      <h2>{result}</h2>
    </>
  )
}

function App() {
  return (
    <>
      <CaptchaSolve />
    </>
  )
}

export default App
