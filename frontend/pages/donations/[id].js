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
           const  data  = await axios.get(`/api/donations?id=${id}`);
           setDonation(data);
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
  <Layout>

  <Link href="/donations">Back to Donations</Link>
        <Image src={donation.src}/> 
              </Layout>
)
}
