
import axios from "axios"


export default async function handler(req, res) {
  const { id } = req.query
  try {
    if (id) {
      const { data } = await axios.get(`http://localhost:8080/donations/${id}`)
      res.status(200).json(data)
      console.log("getting from id page")
    }
    else {
      console.log("getting from all page")
      const response = await axios.get("http://localhost:8080/donations/donationList");
      console.log(response.data)
      res.status(200).json(response.data);

    }

  }
  catch (error) {
    res.status(500).json({ "error": error.message })
  }

}