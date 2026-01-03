import { useState } from "react"
import { useNavigate } from "react-router-dom"
import '../styles/LoginPage.css'


export default function LoginPage() {


    const [email, setEmail] = useState("")

    const [password, setPassword] = useState("")

    const navigate = useNavigate()


    const login = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            const res = await fetch(
                `http://localhost:9090/users/${email}?password=${password}`,
                { method: "GET" }
            )
            if (!res.ok) {
                alert(`error ${res.statusText}`)
            }
            else {
                alert('login successful')
                navigate(`/${encodeURIComponent(email)}?password=${encodeURIComponent(password)}`)
            }
        }
        catch (e) {
            alert(`failed to send login request ${e}`)
        }

    }


    return (
        <div className='login-page'>
            <h1>login page</h1>
            <form className='login-form' onSubmit={login}>

                <div className='segment'>
                    <label title="email">email:</label>
                    <input type="text" onChange={(e) => setEmail(e.target.value)} />
                </div>
                <div className='segment'>
                    <label title="password"> password:</label>
                    <input type='password' onChange={(e) => setPassword(e.target.value)} />
                </div>
                <button type='submit'>login</button>
            </form>
        </div>
    )
}