import axios from "axios"

export default async function handler(req, res) {
    const { id } = req.query
    try {
        const { data } = await axios.get(`http://localhost:8080/users/${id}`)
        res.status(200).json(data)
    }
    catch {
        res.status(404).json({ "error": error.message })
    }


}