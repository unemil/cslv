import { useCallback, useState } from 'react'
import { useDropzone } from 'react-dropzone'
import axios from 'axios'
import './App.css'

const CaptchaSolve = () => {
  const [result, setResult] = useState('')
  const [imageURL, setImageURL] = useState('')

  const sendData = async (file) => {
    const formData = new FormData()
    formData.append('file', file)

    setImageURL(URL.createObjectURL(file))

    await axios
      .request({
        method: 'POST',
        url: 'http://80.90.186.78/api/v1/captcha/solve',
        data: formData,
        headers: {
          'Content-Type': 'multipart/form-data',
        }
      })
      .then(response => { setResult(response.data.solution) })
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
        <p>Перетащите файл сюда или нажмите, чтобы выбрать файл</p>
      </div>
      {imageURL && <img src={imageURL}></img>}
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
