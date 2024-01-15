import React, { useState } from 'react';
import axios from 'axios';
import { useNavigate } from 'react-router-dom';

export default function Login() {

    const navigate = useNavigate();

    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

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

    return (
        <div>
            <h2>Login</h2>
            <input type="text" placeholder="Email" value={email} onChange={(e) => setEmail(e.target.value)} />
            <input type="text" placeholder="Password" value={password} onChange={(e) => setPassword(e.target.value)} />
            <button onClick={handleLogin}>Login</button>
        </div>
    )
}