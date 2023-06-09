import React from "react"
import axios from "axios"

export default function RegistrationPage() {

    // state for data in the form
    const [formData, setFormData] = React.useState({
        username: "",
        email: "",
        password: "",
        passwordConfirmation: ""
    })

    // event handler triggered when form data changes
    function handleChange(event) {
        const { name, value } = event.target;
        setFormData(prevFormData => ({
            ...prevFormData,
            [name]: value
        }))
    }

    // event handler for form submission
    async function handleSubmit(event) {
        event.preventDefault();
        if (formData.username === "" || formData.email === "" || 
            formData.password === "" || formData.passwordConfirmation === "") {
                console.log("Missing registration info")
                return
        } else if (formData.password !== formData.passwordConfirmation) {
            console.log("Passwords do not match")
            return
        }

        // create JSON data for POST request
        const userData = {
            username: formData.username,
            email: formData.email,
            password: formData.password
        }

        // send request to registration route
        try {
            const res = await axios.post("http://localhost:8000/users/register", userData)
            console.log(res)
            console.log("Successfully registered user")
        } catch (error) {
            console.log(error)
            return
        }

        // reset state
        setFormData({
            username: "",
            email: "",
            password: "",
            passwordConfirmation: ""
        })
    };

    return (
        <div className="registration-container">
            <form className="registration-form" onSubmit={handleSubmit}>
                <input
                    type="text"
                    placeholder="Username"
                    className="registration-form--input"
                    name="username"
                    onChange={handleChange}
                    value={formData.username}
                />
                <input
                    type="email"
                    placeholder="Email Address"
                    className="registration-form--input"
                    name="email"
                    onChange={handleChange}
                    value={formData.email}
                />
                <input
                    type="password"
                    placeholder="Password"
                    className="registration-form--input"
                    name="password"
                    onChange={handleChange}
                    value={formData.password}
                />
                <input
                    type="password"
                    placeholder="Confirm password"
                    className="registration-form--input"
                    name="passwordConfirmation"
                    onChange={handleChange}
                    value={formData.passwordConfirmation}
                />
                <button className="registration-form--submit">
                    Register
                </button>
            </form>
        </div>
    )
}