import { useState } from "react";
import { useNavigate } from "react-router-dom";



export default function CreatePage(){
    
    const roles = ['User','Student']
    const [email,setEmail]= useState('')
    const [password,setPassword] = useState('')
    const [username,setUsername]= useState('')
    const navigate = useNavigate()

    const createAccount= async (e: React.FormEvent)=>{
        e.preventDefault()

        try{
            const resp = await fetch('http://localhost:9090/users',
                {
                    method:"POST",
                    headers:{ 'Content-Type': 'application/json' },
                    body: JSON.stringify({'email': email, 'password': password, 'username': username, 'roles': roles})
                }    
            )
            if (resp.ok) {
                alert('account created successfully')
                navigate(`/${encodeURIComponent(email)}?password=${encodeURIComponent(password)}`)
            }
            else {
                alert(`${resp.statusText}`)
            }
        } catch(e){
            alert(`failed to create account ${e}`)
        }
        setEmail('')
        setUsername('')
        setPassword('')

    }
    return(
        <div>
            <h1>Create a New Account</h1>
            <form className='login-form' onSubmit={createAccount}>
                <div className='segment'>
                    <label title="email">email:</label>
                    <input type="text" onChange={(e) => setEmail(e.target.value)} />
                </div>
                <div className='segment'>
                    <label title="username"> username:</label>
                    <input type='text' onChange={(e) => setUsername(e.target.value)} />
                </div>
                <div className='segment'>
                    <label title="password"> password:</label>
                    <input type='password' onChange={(e) => setPassword(e.target.value)} />
                </div>
                <button type='submit'>Create</button>
            </form>
        </div>
    );
}