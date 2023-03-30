
import axios from "axios"
const api=axios.create({baseUrl:"http://localhost:8080"})

export async function handler(req,res)
{
const {id}=req.query
try
{
if(id)
{
    const {data}=await api.get(`/donations/${id}`) 
    res.status(200).json(data)
}
else
{
    const {data}=await api.get("/donations/donationList") //get donations from server side 
    res.status(200).json(data)
    
}

}
catch{
res.status(500).json({error:error.message})
}

}