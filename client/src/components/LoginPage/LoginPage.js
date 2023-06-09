import React from "react"
import axios from "axios"

export default function LoginPage({ setToken }) {

    const [formData, setFormData] = React.useState({
        email: "",
        password: ""
    })

    function handleChange(event) {
        const { name, value } = event.target;
        setFormData(prevFormData => ({
            ...prevFormData,
            [name]: value
        }))
    }

    async function handleSubmit(event) {
        event.preventDefault();
        if (formData.email === "" || formData.password === "") {
            console.log("Missing login info")
            return
        }
        
        // send request to login route
        try {
            const res = await axios.post("http://localhost:8000/users/login", formData)
            const data = await res.data
            const receivedToken = await data.token
            setToken(receivedToken)
            console.log("Successfully logged in user")
        } catch (error) {
            console.log(error)
            return
        }

        // reset state
        setFormData({
            email: "",
            password: ""
        })

        window.location.href = '/dashboard'
    }

    return (
        <div className="login-container">
            <form className="login-form" onSubmit={handleSubmit}>
                <input
                    type="email"
                    placeholder="Email Address"
                    className="login-form--input"
                    name="email"
                    onChange={handleChange}
                    value={formData.email}
                />
                <input
                    type="password"
                    placeholder="Password"
                    className="login-form--input"
                    name="password"
                    onChange={handleChange}
                    value={formData.password}
                />
                <button className="login-form--submit">
                    Login
                </button>
            </form>
        </div>
    )
}