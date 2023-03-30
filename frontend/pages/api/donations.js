
import axios from "axios"
const api=axios.create({baseUrl:"http://localhost:8080"})

export default async function handler(req,res)
{
const {id}=req.query
try
{
if(id)
{
    const { data }=await api.get(`/donations/${id}`) 
    res.status(200).json(data)
    console.log("getting from id page")
}
else
{
  console.log("getting from all page")
  const { data } = await api.get("/donations/donationList");
  res.status(200).json(data);
    
}

}
catch (error){
res.status(500).json({"error":error.message})
}

}