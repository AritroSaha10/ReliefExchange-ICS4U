import Link from "next/link";
import Image from "next/image";
import Head from "next/head"
//.. is move out of folder . is move out of file


export default function donations({donations}) {

    return(
        <>
        <Head> <title>Donations</title></Head>
           <h1>Donations</h1>
        <h2>  <Link href="/">Back to home page</Link></h2>
<ul>
  {donations.map((donation)=>(
  <li key={donation.id}>
    <Link href={`/donations/${donation.id}`}>{donation.name}</Link>
    </li>
  ))}
</ul>
  {/* <Image src="/images/profile.jpg" height="144" width="120" alt="profile"/> */}
 
    </>
    ) 
  }
  export async function getServerSideProps()
  {
const res=await fetch("http://localhost:4000/donations/donationList") //get donations from server side 
const donations=await res.json();
return {
  props: {donations},
};
  }