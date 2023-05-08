import Link from "next/link";
import Image from "next/image";
import Head from "next/head"
import axios from "axios";
import { useState, useEffect } from "react";
//.. is move out of folder . is move out of file


export default function displayDonations() {
  const [donations, setDonations] = useState([])
  useEffect(() => {
    const fetchDonations = async () => {
      try {
        const { data } = await axios.get("/api/donations");

        setDonations(data);
      } catch (error) {
        console.error("Error fetching donations:", error);
      }
    };

    fetchDonations();
  }, []);

  return (
    <>
      <Head> <title>Donations</title></Head>
      <h1>Donations</h1>
      <h2>  <Link href="/">Back to home page</Link></h2>
      <ul>
        {donations.map((donation) => (
          <li key={donation.id}>
            {/* <Image src={donation.images[0]}/> */}

            <Link href={`/donations/${donation.id}`}>{donation.title}</Link>
            <p>{donation.city}, {donation.location}</p>
            <p>{donation.description}</p>
          </li>
        ))}
      </ul>
      {/* <Image src="/images/profile.jpg" height="144" width="120" alt="profile"/> */}

    </>
  )
}
