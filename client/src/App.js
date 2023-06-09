import React from "react"
import axios from "axios"
import { BrowserRouter as Router, Route, Routes } from "react-router-dom"

import HomePage from "./components/HomePage/HomePage"
import LoginPage from "./components/LoginPage/LoginPage"
import RegistrationPage from "./components/RegistrationPage/RegistrationPage"
import Dashboard from "./components/Dashboard/Dashboard"

function App() {

  // async function fetchItems() {
  //   try {
  //       const res = await axios.get("http://localhost:8000/items/", { 
  //           headers: {
  //               Authorization: `Bearer ${token}`,
  //           },
  //       })
  //       const data = await res.data
  //       setItems(data)
  //   } catch (error) {
  //       console.log(error)
  //       return
  //   }
  // }

  return (
    <Router>
      <div className="App">
       <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/register" element={<RegistrationPage />} />
        <Route path="/login" element={<LoginPage setToken={setToken}/>} />
        <Route path="/users" element={<PrivateRoute />} />
       </Routes>
      </div>
    </Router>
  );
}

export default App;
