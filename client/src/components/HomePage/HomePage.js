import React from "react"
import { Link } from "react-router-dom"

export default function HomePage() {
    return (
        <div className="homepage-container">
            <h1>Home Page</h1>
            <Link to="/register">Register</Link>
            <Link to="/login">Login</Link>
        </div>
    )
}