import Link from "next/link";

import Image from "next/image";
//  import Layout from "../../components/Layout" 
 import { useRouter } from "next/router";
 import { useState, useEffect } from "react";
 import axios from "axios";
 
 export default function Donation() {
   const router = useRouter();
   const { id } = router.query;
   const [donation, setDonation] = useState(null);
 
   useEffect(() => {
     if (id) {
       const fetchDonation = async () => {
         try {
           const  res  = await axios.get(`/api/donations?id=${id}`);
           setDonation(res.data);
         } catch (error) {
           console.error("Error fetching donation:", error);
         }
       };
 
       fetchDonation();
     }
   }, [id]); //if id changes, then rerun function 
 
   if (!donation) {
     return <div>Loading...</div>;
   }
return(
  <>

  <Link href="/donations/donationList">Back to Donations</Link>
      <h1>{donation.title}</h1>
      <p>{donation.city}, {donation.location}</p>
      <p>{donation.description}</p>
        {/* <Image src={donation.src}/>  */}
              </>
)
}
