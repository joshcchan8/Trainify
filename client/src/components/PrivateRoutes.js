import React from "react";

export default function PrivateRoutes() {

    async function handleLogin() {
        try {
            const res = await axios.post('http://localhost:8000/users/login', {
                email,
                password
            });

            const { token } = res.data;
            localStorage.setItem('token', token);
            navigate('/dashboard')
        } catch (error) {
            console.error('Login failed', error)
        }
    }

    return(
        <Outlet />
    )
}