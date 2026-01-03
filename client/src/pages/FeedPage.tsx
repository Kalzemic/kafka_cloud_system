import { useEffect, useState } from "react"
import { useParams, useSearchParams } from "react-router-dom"
import '../styles/FeedPage.css'


export default function FeedPage() {

    const { email } = useParams<{ email: string }>()
    const [searchParams] = useSearchParams()
    const password = searchParams.get('password')
    const [content, setContent] = useState('')

    type Post = {
        email: string
        content: string
        timestamp: string
    }

    const [posts, setPosts] = useState<Post[]>([])
    useEffect(() => {
        if (!email || !password) return

        const url =
            `http://localhost:9090/posts/listen/${encodeURIComponent(email)}` +
            `?password=${encodeURIComponent(password)}`

        const es = new EventSource(url)

        es.addEventListener('post', (e) => {
            const post = JSON.parse(e.data)
            setPosts((prev) => [...prev, post])
        })

        es.onerror = (err) => {
            console.error("SSE error", err)
            es.close()
        }
        return () => {
            es.close()
        }

    }, [email, password])


    const createPost = async (e: React.FormEvent) => {
        e.preventDefault()
        try {
            const resp = await fetch(`http://localhost:9090/posts/produce/${email}?password=${password}`,
                {
                    method: "POST",
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ 'email': email, 'content': content })
                }
            )
            if (resp.ok) {
                console.log('post sent successfully')
                setContent('')
            }
            else {
                alert(`${resp.statusText}`)
            }
        } catch (e) {
            alert(`failed to send post ${e}`)
        }
    }

    return (
        <div className='feed-page'>
            <h2>Welcome {email}</h2>
            <form className='post-form' onSubmit={createPost}>
                <textarea rows={15} cols={60} onChange={(e) => setContent(e.target.value)}></textarea>
                <button type='submit'>post</button>
            </form>

            <div className='live-posts'>
                {posts.map((p, i) => (
                    <div key={i} className="post">
                        <b>{p.email}</b>
                        <p>{p.content}</p>
                    </div>
                ))}
            </div>
        </div>
    )
}